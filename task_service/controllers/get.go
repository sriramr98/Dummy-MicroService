package controllers

import (
	"github.com/sriramr98/todo_task_service/services"
	"github.com/sriramr98/todo_task_service/utils"
	"net/http"
	"strconv"
)

func GetTask(taskService services.TaskService) utils.ApiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)

		userId := r.Context().Value("userId").(int64)
		if err != nil {
			return utils.ApiError{
				StatusCode: http.StatusBadRequest,
				Code:       utils.ErrInvalidRequestPayload,
				Message:    "Invalid task id",
			}
		}

		task, err := taskService.GetTask(int64(id), userId)
		if err != nil {
			return utils.ApiError{
				StatusCode: http.StatusInternalServerError,
				Code:       utils.ErrInternalServer,
				Message:    "Error getting task",
			}
		}

		utils.WriteSuccessMessage(w, http.StatusOK, task)
		return nil
	}
}
