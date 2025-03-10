package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sriramr98/todo_auth_service/controllers"
	"github.com/sriramr98/todo_auth_service/services"
	"github.com/sriramr98/todo_auth_service/utils"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
		os.Exit(1)
	}

	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Printf("Error running server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	pg := initPG()
	defer func(pg *sql.DB) {
		err := pg.Close()
		if err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}(pg)

	authService := services.NewAuthService(pg)

	srv := NewServer(authService)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler: srv,
	}

	grpcServer, tcpLn := GetGRPCServer(authService)

	go runServer(httpServer)
	go func(grpcServer *grpc.Server, tcpLn net.Listener) {
		err := grpcServer.Serve(tcpLn)
		if err != nil {
			log.Fatalf("Error starting gRPC Server: %v", err)
		}
	}(grpcServer, tcpLn)

	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	go initGracefulShutdown(ctx, httpServer, grpcServer, &wg)
	wg.Wait()
	return nil
}

func initPG() *sql.DB {
	host := os.Getenv("PG_HOST")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	dbname := os.Getenv("PG_DBNAME")

	return InitDB(host, user, password, dbname)
}

func NewServer(authService services.AuthService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", utils.ErrorHandler(controllers.LoginHandler(authService)))
	mux.HandleFunc("POST /register", utils.ErrorHandler(controllers.RegisterHandler(authService)))
	return mux
}

func runServer(server *http.Server) {
	log.Printf("Listening on %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("error listening and serving: %s\n", err)
	}
}

func initGracefulShutdown(ctx context.Context, server *http.Server, grpcServer *grpc.Server, wg *sync.WaitGroup) {
	defer wg.Done()
	<-ctx.Done()
	shutdownCtx := context.Background()
	shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
	defer cancel()

	log.Println("Shutting down HTTP server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Failed to shutdown HTTP server: %v\n", err)
	}

	grpcShutdownCtx := context.Background()
	grpcShutdownCtx, cancel = context.WithTimeout(grpcShutdownCtx, 10*time.Second)
	defer cancel()

	log.Println("Shutting down GRPC server...")
	// This also shuts down the TCP connection, so no need to do it manually
	grpcServer.GracefulStop()
	log.Println("GRPC server stopped")
}
