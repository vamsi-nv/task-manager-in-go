package handlers

import (
	"net/http"

	"task-manager/internal/models"
	"task-manager/internal/services"
	"task-manager/internal/utils"
	"task-manager/internal/validation"
)

type UserHandler struct {
	Service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		Service: service,
	}
}

func (h *UserHandler) SignupUser(w http.ResponseWriter, r *http.Request) error {
	var user models.CreateUserRequest

	err := DecodeStrict(r.Body, &user)
	if err != nil {
		return utils.BadRequest("Invalid JSON", nil)
	}

	err = validation.Validate.Struct(user)
	if err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.BadRequest("Validation Failed", errs)
	}

	created, err := h.Service.CreateUser(r.Context(), &user)
	if err != nil && created != nil {
		utils.ResponseJSON(w, http.StatusAccepted, err.Error(), created)
		return nil
	}

	if err != nil {
		return err
	}

	utils.ResponseJSON(w, http.StatusCreated, "Signup successful. Please verify your email", created)
	return nil
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) error {

	var creds models.Credentials

	err := DecodeStrict(r.Body, &creds)
	if err != nil {
		return utils.BadRequest("Invalid JSON", nil)
	}

	err = validation.Validate.Struct(creds)
	if err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.BadRequest("Validation Failed", errs)
	}

	data, err := h.Service.LoginUser(r.Context(), &creds)
	if err != nil {
		return err
	}

	utils.ResponseJSON(w, http.StatusOK, "Loggedin", data)
	return nil
}

func (h *UserHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) error {
	token := r.URL.Query().Get("token")
	if token == "" {
		utils.BadRequest("Missing email verification token", nil)
	}

	err := h.Service.VerifyEmail(r.Context(), token)
	if err != nil {
		return err
	}

	utils.ResponseJSON(w, http.StatusOK, "Email Verified", nil)
	return nil
}

func (h *UserHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) error {
	var user models.CreateUserRequest

	err := DecodeStrict(r.Body, &user)
	if err != nil {
		return utils.BadRequest("Invalid JSON", nil)
	}

	err = h.Service.ForgotPassword(r.Context(), user.Email)
	if err != nil {
		return err
	}

	utils.ResponseJSON(w, http.StatusOK, "Reset password email has been sent to your email", nil)
	return nil
}

func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) error {
	var req models.UpdatePasswordRequest

	token := r.URL.Query().Get("token")
	if token == "" {
		return utils.BadRequest("Missing reset password token", nil)
	}

	err := DecodeStrict(r.Body, &req)
	if err != nil {
		return utils.BadRequest("Invalid JSON", nil)
	}

	err = validation.Validate.Struct(req)
	if err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.BadRequest("Validation Failed", errs)
	}

	err = h.Service.ResetPassword(r.Context(), token, &req)
	if err != nil {
		return err
	}

	utils.ResponseJSON(w, http.StatusOK, "Password updated", nil)
	return nil
}

func (h *UserHandler) ResendVerificationEmail(w http.ResponseWriter, r *http.Request) error {

	var user models.CreateUserRequest
	err := DecodeStrict(r.Body, &user)
	if err != nil {
		return utils.BadRequest("Invalid JSON", nil)
	}

	if user.Email == "" {
		return utils.BadRequest("Email is required", nil)
	}

	err = h.Service.ResendVerificationEmail(r.Context(), user.Email)
	if err != nil {
		return err
	}

	utils.ResponseJSON(w, http.StatusOK, "Verification email sent to your email", nil)
	return nil
}
