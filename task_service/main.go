package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sriramr98/todo_task_service/controllers"
	"github.com/sriramr98/todo_task_service/services"
	"github.com/sriramr98/todo_task_service/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type GetEnvFunc func(string) string

func main() {

	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
		os.Exit(1)
	}

	ctx := context.Background()
	if err := run(ctx, os.Getenv); err != nil {
		log.Printf("Error running server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, getEnv GetEnvFunc) error {
	pg := initPG(getEnv)
	defer func(pg *sql.DB) {
		err := pg.Close()
		if err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}(pg)

	taskService := services.NewTaskService(pg)

	srv := NewServer(taskService)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", getEnv("PORT")),
		Handler: srv,
	}

	go runServer(httpServer)
	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	go initGracefulShutdown(ctx, httpServer, &wg)
	wg.Wait()
	return nil
}

func initPG(env GetEnvFunc) *sql.DB {
	host := env("PG_HOST")
	user := env("PG_USER")
	password := env("PG_PASSWORD")
	dbname := env("PG_DBNAME")

	return InitDB(host, user, password, dbname)
}

func NewServer(service services.TaskService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /task", utils.ErrorHandler(controllers.CreateTask(service)))
	mux.HandleFunc("GET /task", utils.ErrorHandler(controllers.ListTasks(service)))
	mux.HandleFunc("GET /task/{id}", utils.ErrorHandler(controllers.GetTask(service)))
	mux.HandleFunc("PATCH /task/{id}", utils.ErrorHandler(controllers.UpdateTask(service)))
	mux.HandleFunc("DELETE /task/{id}", utils.ErrorHandler(controllers.DeleteTask(service)))
	return mux
}

func runServer(server *http.Server) {
	log.Printf("Listening on %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("error listening and serving: %s\n", err)
	}
}

func initGracefulShutdown(ctx context.Context, server *http.Server, wg *sync.WaitGroup) {
	defer wg.Done()
	<-ctx.Done()
	shutdownCtx := context.Background()
	shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
	defer cancel()

	log.Println("Shutting down server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Failed to shutdown server: %v\n", err)
	}
	log.Println("Server stopped")
}
