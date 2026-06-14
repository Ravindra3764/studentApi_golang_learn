package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ravindra3764/student-api/student-api/internal/config"
	"github.com/ravindra3764/student-api/student-api/internal/http/handlers/student"
)

func main() {

	//load config

	cfg := config.MustLoad()
	// db setup
	// setup router

	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())
	// setup server

	server := http.Server{
		Addr:    cfg.HttpServer.Addr,
		Handler: router,
	}

	slog.Info("Server started", slog.String("Address", cfg.HttpServer.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}

	}()

	<-done

	slog.Info("Shutting down your server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {

		slog.Error("failed to shutdown server", slog.String("Error", err.Error()))

	}

	slog.Info("Server shut successfull")

}
