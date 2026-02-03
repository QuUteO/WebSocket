package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/QuUteO/video-communication/internal/app"
)

func main() {
	// Создание приложения
	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	// Канал для обработки сигналов
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Запуск сервера в горутине
	go func() {
		if err := application.Run(); err != nil {
			panic(err)
		}
	}()

	// Ожидание сигнала для graceful shutdown
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		panic(err)
	}
}
