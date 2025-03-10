package controllers

import (
	"github.com/sriramr98/todo_task_service/services"
	"github.com/sriramr98/todo_task_service/utils"
	"net/http"
	"strconv"
)

func DeleteTask(taskService services.TaskService) utils.ApiHandler {
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

		if err := taskService.DeleteTask(id, userId); err != nil {
			return utils.ApiError{
				StatusCode: http.StatusInternalServerError,
				Code:       utils.ErrInternalServer,
				Message:    "Error deleting task",
			}
		}

		utils.WriteSuccessMessage(w, http.StatusOK, map[string]string{"message": "Task deleted successfully"})
		return nil
	}
}
