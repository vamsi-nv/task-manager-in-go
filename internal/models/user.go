package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                          primitive.ObjectID `bson:"_id"`
	Username                    string             `bson:"username"`
	Email                       string             `bson:"email"`
	Password                    string             `bson:"password"`
	Verified                    bool               `bson:"verified"`
	CreatedAt                   time.Time          `bson:"created_at"`
	UpdatedAt                   time.Time          `bson:"updated_at"`
	VerificationToken           string             `bson:"verification_token,omitempty"`
	VerificationTokenExpiresAt  *time.Time         `bson:"verification_token_expires_at,omitempty"`
	PasswordResetToken          string             `bson:"password_reset_token,omitempty"`
	PasswordResetTokenExpiresAt *time.Time         `bson:"password_reset_token_expires_at,omitempty"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"_id"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type Credentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type UpdatePasswordRequest struct {
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}
