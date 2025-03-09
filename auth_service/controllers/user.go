package controllers

import (
	"github.com/sriramr98/todo_auth_service/utils"
	"net/http"
)

func GetUserInfoHandler() utils.ApiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		utils.WriteSuccessMessage(w, http.StatusOK, map[string]string{"message": "User info retrieved"})
		return nil
	}
}
