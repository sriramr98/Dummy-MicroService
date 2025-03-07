package controllers

import (
	"fmt"
	"github.com/sriramr98/todo_task_service/services"
	"github.com/sriramr98/todo_task_service/utils"
	"net/http"
)

func ListTasks(taskService services.TaskService) utils.ApiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		fmt.Fprintf(w, "List controllers")
		return nil
	}
}
