package controllers

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sriramr98/todo_auth_service/services"
	"github.com/sriramr98/todo_auth_service/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (lb loginBody) Validate() map[string]string {
	errorDetails := make(map[string]string)

	if strings.TrimSpace(lb.Email) == "" {
		errorDetails["email"] = "Email is required"
	}

	if strings.TrimSpace(lb.Password) == "" {
		errorDetails["password"] = "Password is required"
	}

	return errorDetails
}

func LoginHandler(authService services.AuthService) utils.ApiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var loginBody loginBody
		if err := json.NewDecoder(r.Body).Decode(&loginBody); err != nil {
			return utils.ApiError{
				StatusCode: http.StatusBadRequest,
				Code:       utils.ErrInvalidRequestPayload,
				Message:    "Invalid request body",
			}
		}

		if errors := loginBody.Validate(); len(errors) > 0 {
			return utils.ValidationError{
				StatusCode: http.StatusBadRequest,
				Code:       utils.ErrInvalidRequestPayload,
				Errors:     errors,
			}
		}

		userInfo, err := authService.GetUserInfoByEmail(loginBody.Email)
		if err != nil {
			return utils.ApiError{
				StatusCode: http.StatusBadRequest,
				Code:       utils.ErrUserNotFound,
				Message:    "User not found",
			}
		}

		isPasswordValid := utils.ValidateHashPassword(loginBody.Password, userInfo.Password)
		if !isPasswordValid {
			return utils.ApiError{
				StatusCode: http.StatusBadRequest,
				Code:       utils.ErrInvalidPassword,
				Message:    "Invalid password",
			}
		}

		jwtToken, err := createJWTToken(userInfo)
		if err != nil {
			return utils.ApiError{
				StatusCode: http.StatusInternalServerError,
				Code:       utils.ErrInternalServer,
				Message:    "Error generating JWT token",
			}
		}

		utils.WriteSuccessMessage(w, http.StatusOK, map[string]string{"message": "Login succeeded", "token": jwtToken})
		return nil
	}
}

func createJWTToken(userInfo services.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	hmacSampleSecret := []byte(secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(userInfo.ID),
		"iss": "todo_app",
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	})

	return token.SignedString(hmacSampleSecret)
}
