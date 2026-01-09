package middleware

import (
	"net/http"

	"task-manager/internal/utils"
)

func WithError(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err == nil {
			return
		}

		if appErr, ok := err.(*utils.AppError); ok {
			utils.ErrorJSON(w, appErr.Code, appErr.Error(), appErr.Errors)

			return
		}
		utils.ErrorJSON(w, http.StatusInternalServerError, "Internal Server error", nil)
	}
}
