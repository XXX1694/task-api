package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task-api/internal/handlers"
	"task-api/internal/middleware"
	"task-api/internal/models"
	"time"
)

func main() {
	store := models.NewTaskStore()
	taskHandler := handlers.NewTaskHandler(store)

	apiMux := http.NewServeMux()

	apiMux.HandleFunc("/v1/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.GetTasks(w, r)
		case http.MethodPost:
			taskHandler.CreateTask(w, r)
		case http.MethodPatch:
			taskHandler.UpdateTask(w, r)
		case http.MethodDelete:
			taskHandler.DeleteTask(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	apiMux.HandleFunc("/v1/external/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		taskHandler.GetExternalTasks(w, r)
	})

	protectedHandler := middleware.RateLimitMiddleware(middleware.RequestIDMiddleware(middleware.LoggingMiddleware(middleware.AuthMiddleware(apiMux))))

	mainMux := http.NewServeMux()

	mainMux.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handlers.SwaggerHandler(w, r)
	})

	mainMux.HandleFunc("/swagger-ui", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handlers.SwaggerUIHandler(w, r)
	})

	mainMux.Handle("/", protectedHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mainMux,
	}

	go func() {
		log.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
