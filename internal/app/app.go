// Package app - пакет-аккумулятор собирающий приложение.
package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"makves-task/internal/handler"
	"makves-task/internal/repository"
)

// Run - запускает приложение.
func Run() {
	// Инициализируем БД.
	database := repository.NewGormLiteSQL()
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler.NewRouter(database),
	}
	// Запускаем сервер.
	go func() {
		log.Println("starting server at:", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	// Канал Grace-ful Shutdown.
	sigs := make(chan os.Signal)
	signal.Notify(sigs,
		syscall.SIGINT,
		os.Interrupt)
	// После получения сигнала закрываем приложение.
	<-sigs
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatal("server shutdown error")
	}
	log.Println("shutting down")
	os.Exit(0)
}
