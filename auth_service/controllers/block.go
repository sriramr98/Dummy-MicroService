package controllers

import (
	"github.com/sriramr98/todo_auth_service/services"
	"github.com/sriramr98/todo_auth_service/utils"
	"net/http"
	"os"
	"strconv"
)

func BlockUserHandler(authService services.AuthService) utils.ApiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		adminPwd := r.Header.Get("AdminPassword")
		if adminPwd != os.Getenv("ADMIN_PASSWORD") {
			return utils.ApiError{
				StatusCode: http.StatusUnauthorized,
				Code:       utils.ErrUnauthorized,
				Message:    "Unauthorized",
			}
		}

		userIdStr := r.PathValue("userId")
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			return utils.ApiError{
				StatusCode: http.StatusBadRequest,
				Code:       utils.ErrInvalidRequestPayload,
				Message:    "Invalid user ID",
			}
		}

		if blocked := authService.BlockUser(int64(userId)); !blocked {
			return utils.ApiError{
				StatusCode: http.StatusInternalServerError,
				Code:       utils.ErrInternalServer,
				Message:    "Error blocking user",
			}
		}

		return nil
	}
}
