package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"task-manager/internal/config"
	"task-manager/internal/database"
	"task-manager/internal/handlers"
	"task-manager/internal/middleware"
	"task-manager/internal/repository"
	"task-manager/internal/routes"
	"task-manager/internal/services"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading env")
	}

	cfg := config.Primary{
		MongoUri: os.Getenv("MONGO_URI"),
		Database: os.Getenv("DATABASE_NAME"),
		Port:     os.Getenv("PORT"),
	}

	client, err := database.Connect(cfg.MongoUri)
	if err != nil {
		log.Fatalln("Error Connecting to database:", err)
	}
	defer client.Disconnect(context.Background())

	taskRepo := repository.NewTaskRespository(client, cfg.Database)
	taskService := services.NewTaskService(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskService)

	userRepo := repository.NewUserRespository(client, cfg.Database)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	mux := http.NewServeMux()
	limiter := middleware.NewRateLimiter(1, 2.0)
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	routes.TaskRouter(mux, taskHandler)
	routes.UserRouter(mux, userHandler)

	secureMux := middleware.ApplyMiddleware(mux, limiter.LimitMiddleware, middleware.JWTMiddleware)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: secureMux,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		fmt.Printf("Server running on port :%s...\n", cfg.Port)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalln("Error starting server:", err)
		}
	}()

	// shutdown shit
	<-ctx.Done()
	fmt.Println("Shutting down server...")

	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = server.Shutdown(shutdownCtx)
	if err != nil {
		log.Println("Server forced to shutdown:", err)
	}

	err = client.Disconnect(shutdownCtx)
	if err != nil {
		log.Println("Error closing database:", err)
	}

	fmt.Println("Shutdown complete.")
}
