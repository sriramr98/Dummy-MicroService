package controllers

import (
	"github.com/sriramr98/todo_task_service/services"
	"github.com/sriramr98/todo_task_service/utils"
	"net/http"
)

func ListTasks(taskService services.TaskService) utils.ApiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		tasks, err := taskService.ListTasks()
		if err != nil {
			return utils.ApiError{
				StatusCode: http.StatusInternalServerError,
				Code:       utils.ErrInternalServer,
				Message:    "Error listing tasks",
			}
		}

		utils.WriteSuccessMessage(w, http.StatusOK, tasks)
		return nil
	}
}
