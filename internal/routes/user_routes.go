package routes

import (
	"net/http"

	"task-manager/internal/handlers"
	"task-manager/internal/middleware"
)

func UserRouter(mux *http.ServeMux, h *handlers.UserHandler) {
	mux.HandleFunc("POST /api/auth/sign-up", middleware.WithError(h.SignupUser))
	mux.HandleFunc("POST /api/auth/login", middleware.WithError(h.LoginUser))
	mux.HandleFunc("POST /api/auth/forgot-password", middleware.WithError(h.ForgotPassword))
	mux.HandleFunc("POST /api/auth/reset-password", middleware.WithError(h.ResetPassword))
	mux.HandleFunc("GET /api/auth/verify-email", middleware.WithError(h.VerifyEmail))
	mux.HandleFunc("POST /api/auth/resend-verification", middleware.WithError(h.ResendVerificationEmail))

}
