package controllers

import (
	"github.com/sriramr98/todo_task_service/services"
	"github.com/sriramr98/todo_task_service/utils"
	"net/http"
)

func GetTask(taskService services.TaskService) utils.ApiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Get task"))
		return nil
	}
}
