package controllers

import (
	"github.com/sriramr98/todo_task_service/services"
	"github.com/sriramr98/todo_task_service/utils"
	"net/http"
)

func ListTasks(taskService services.TaskService) utils.ApiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		userId := r.Context().Value("userId").(int64)

		tasks, err := taskService.ListTasks(userId)
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
