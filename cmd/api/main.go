package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "test-task/internal/infrastructure/postgres"
    postgresRepo "test-task/internal/repository/postgres"
    httpTransport "test-task/internal/transport/http"
    "test-task/internal/transport/http/handlers"
    "test-task/internal/usecase/task"
)

func main() {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        dbURL = "postgres://postgres:postgres@localhost:5432/taskdb?sslmode=disable"
    }

    pool, err := postgres.NewPool(context.Background(), dbURL)
    if err != nil {
        log.Fatalf("failed to create pool: %v", err)
    }
    defer pool.Close()

    taskRepo := postgresRepo.NewTaskRepository(pool)
    scheduler := task.NewScheduler(taskRepo)
    taskService := task.NewService(taskRepo, scheduler)
    taskHandler := handlers.NewTaskHandler(taskService)

    router := httpTransport.NewRouter(taskHandler)

    go func() {
        ticker := time.NewTicker(1 * time.Hour)
        defer ticker.Stop()

        for range ticker.C {
            if err := scheduler.ScheduleRecurringTasks(context.Background()); err != nil {
                log.Printf("failed to schedule recurring tasks: %v", err)
            }
        }
    }()

    srv := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }

    go func() {
        log.Println("starting server on :8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("failed to start server: %v", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("server forced to shutdown: %v", err)
    }

    log.Println("server exited")
}