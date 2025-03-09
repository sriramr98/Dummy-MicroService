package controllers

import (
	"encoding/json"
	"errors"
	"github.com/sriramr98/todo_auth_service/services"
	"github.com/sriramr98/todo_auth_service/utils"
	"net/http"
)

func RegisterHandler(authService services.AuthService) utils.ApiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var user services.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			return utils.ApiError{
				StatusCode: http.StatusBadRequest,
				Code:       utils.ErrInvalidRequestPayload,
				Message:    "Invalid request body",
			}
		}

		if err := user.Validate(); len(err) > 0 {
			return utils.ValidationError{
				StatusCode: http.StatusBadRequest,
				Code:       utils.ErrInvalidRequestPayload,
				Errors:     err,
			}
		}

		password, err := utils.HashPassword(user.Password)
		if err != nil {
			return utils.ApiError{
				StatusCode: http.StatusBadRequest,
				Code:       utils.ErrInvalidPassword,
				Message:    err.Error(),
			}
		}

		userInfo, err := authService.CreateNewUser(user.Email, password)

		if err != nil {
			statusCode := http.StatusInternalServerError
			code := utils.ErrInternalServer
			if errors.Is(err, services.ErrUserExists) {
				statusCode = http.StatusBadRequest
				code = utils.ErrUserExists
			}
			return utils.ApiError{
				StatusCode: statusCode,
				Code:       code,
				Message:    err.Error(),
			}
		}

		userInfo.Password = ""

		utils.WriteSuccessMessage(w, http.StatusOK, userInfo)
		return nil
	}
}
