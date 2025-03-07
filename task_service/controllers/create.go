package controllers

import (
	"encoding/json"
	"github.com/sriramr98/todo_task_service/services"
	"github.com/sriramr98/todo_task_service/utils"
	"net/http"
)

type createTaskPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func CreateTask(taskService services.TaskService) utils.ApiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		decoder := json.NewDecoder(r.Body)
		var taskBody createTaskPayload

		if err := decoder.Decode(&taskBody); err != nil {
			return utils.ApiError{
				StatusCode: http.StatusBadRequest,
				Code:       utils.ErrInvalidRequestPayload,
				Message:    "Invalid request body payload",
				//Errors: map[string]string{
				//	"title": "Invalid request payload",
				//},
			}
		}

		return nil
	}
}
