package controllers

import (
	"encoding/json"
	"github.com/sriramr98/todo_task_service/services"
	"github.com/sriramr98/todo_task_service/utils"
	"log"
	"net/http"
	"strings"
)

type createTaskPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (c createTaskPayload) Validate() map[string]string {
	errors := make(map[string]string)

	if strings.TrimSpace(c.Title) == "" {
		errors["title"] = "Title is required"
	}

	if strings.TrimSpace(c.Body) == "" {
		errors["body"] = "Body is required"
	}

	return errors
}

func CreateTask(taskService services.TaskService) utils.ApiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {

		userId := r.Context().Value("userId").(int64)

		decoder := json.NewDecoder(r.Body)
		var taskBody createTaskPayload

		if err := decoder.Decode(&taskBody); err != nil {
			return utils.ApiError{
				StatusCode: http.StatusBadRequest,
				Code:       utils.ErrInvalidRequestPayload,
				Message:    "Invalid request body",
			}
		}

		if errors := taskBody.Validate(); len(errors) > 0 {
			return utils.ValidationError{
				StatusCode: http.StatusBadRequest,
				Code:       utils.ErrInvalidRequestPayload,
				Errors:     errors,
			}
		}

		id, err := taskService.CreateTask(taskBody.Title, taskBody.Body, userId)

		if err != nil {
			log.Println(err)
			return utils.ApiError{
				StatusCode: http.StatusInternalServerError,
				Code:       utils.ErrInternalServer,
				Message:    "Error creating task",
			}
		}

		utils.WriteSuccessMessage(w, http.StatusCreated, map[string]int64{
			"id": id,
		})

		return nil
	}
}
