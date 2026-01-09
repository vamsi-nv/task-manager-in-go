package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"task-manager/internal/models"
	"task-manager/internal/repository"
	"task-manager/internal/utils"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		Repo: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.CreateUserRequest) (*models.UserResponse, error) {
	existingUser, _ := s.Repo.GetUserByEmail(ctx, user.Email)
	if existingUser != nil {
		return nil, utils.BadRequest("Email already exists", nil)
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, utils.Internal("Error processing password", nil)
	}

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, utils.Internal("Internal security error", nil)
	}
	token := hex.EncodeToString(b)

	expAt := time.Now().UTC().Add(24 * time.Hour)

	newUser := &models.User{
		Username:                   user.Username,
		Email:                      user.Email,
		Password:                   hashedPassword,
		Verified:                   false,
		VerificationToken:          token,
		VerificationTokenExpiresAt: &expAt,
	}

	err = s.Repo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	err = utils.SendVerificationEmail(newUser.Email, token)

	resUser := &models.UserResponse{
		ID:        newUser.ID,
		Username:  newUser.Username,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
	}

	if err != nil {
		return resUser, utils.Internal("Account created, but failed to send verification email. Please try 'Resend Verification'.", nil)
	}

	return resUser, nil
}

func (s *UserService) LoginUser(ctx context.Context, creds *models.Credentials) (any, error) {
	user, _ := s.Repo.GetUserByEmail(ctx, creds.Email)
	if user == nil {
		return nil, utils.Unauthorized("Invalid email or password", nil)
	}

	if !user.Verified {
		return nil, utils.Forbidden("Please verify your email before logging in", nil)
	}

	err := utils.VerifyPassword(creds.Password, user.Password)
	if err != nil {
		return nil, utils.Unauthorized("Invalid email or password", nil)
	}

	token, err := utils.CreateTokenWithClaims(*user)
	if err != nil {
		return nil, utils.Internal("Error Logging In", nil)
	}

	return struct {
		Token string              `json:"token"`
		User  models.UserResponse `json:"user"`
	}{
		Token: token,
		User: models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

func (s *UserService) ForgotPassword(ctx context.Context, email string) error {
	if email == "" {
		return utils.BadRequest("Email is required", nil)
	}
	existingUser, err := s.Repo.GetUserByEmail(ctx, email)
	if existingUser == nil {
		return utils.NotFound("User not found", nil)
	}

	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return utils.Internal("Error sending email", nil)
	}

	token := hex.EncodeToString(b)
	err = utils.SendForgotPasswordEmail(email, token)
	if err != nil {
		return utils.Internal("Error sending email", nil)
	}
	expAt := time.Now().UTC().Add(10 * time.Minute)
	existingUser.PasswordResetToken = token
	existingUser.PasswordResetTokenExpiresAt = &expAt

	err = s.Repo.UpdatePasswordToken(ctx, existingUser)
	if err != nil {
		return utils.Internal("Error saving reset token", nil)
	}

	return nil
}

func (s *UserService) ResetPassword(ctx context.Context, token string, req *models.UpdatePasswordRequest) error {

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	req.Password = hashedPassword

	err = s.Repo.UpdatePassword(ctx, token, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) VerifyEmail(ctx context.Context, token string) error {

	err := s.Repo.VerifyEmail(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) ResendVerificationEmail(ctx context.Context, email string) error {

	existingUser, err := s.Repo.GetUserByEmail(ctx, email)
	if existingUser == nil {
		return utils.NotFound("User not found. Please signup before verification", nil)
	}

	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		utils.Internal("Error sending verification email", nil)
	}

	token := hex.EncodeToString(b)
	err = utils.SendVerificationEmail(email, token)
	if err != nil {
		return utils.Internal("Error sending verification email", nil)
	}

	return nil
}
