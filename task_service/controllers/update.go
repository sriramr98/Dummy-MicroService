package controllers

import (
	"encoding/json"
	"github.com/sriramr98/todo_task_service/services"
	"github.com/sriramr98/todo_task_service/utils"
	"log"
	"net/http"
	"strconv"
)

type updateTaskPayload struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	Completed bool   `json:"completed"`
}

func UpdateTask(taskService services.TaskService) utils.ApiHandler {
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

		decoder := json.NewDecoder(r.Body)
		var taskBody updateTaskPayload

		if err := decoder.Decode(&taskBody); err != nil {
			return utils.ApiError{
				StatusCode: http.StatusBadRequest,
				Code:       utils.ErrInvalidRequestPayload,
				Message:    "Invalid request body",
			}
		}

		err = taskService.UpdateTask(int64(id), taskBody.Title, taskBody.Body, taskBody.Completed, userId)

		if err != nil {
			log.Println(err)
			return utils.ApiError{
				StatusCode: http.StatusInternalServerError,
				Code:       utils.ErrInternalServer,
				Message:    "Error updating task",
			}
		}

		utils.WriteSuccessMessage(w, http.StatusOK, map[string]string{"message": "Task updated successfully"})
		return nil
	}
}
