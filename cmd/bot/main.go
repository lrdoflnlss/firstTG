package main

import (
	"context"
	"flag"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
	"tg-botv1/internal/bot"
	"tg-botv1/internal/logger"
)

var configPath = flag.String("c", ".env", "Path to config")

func main() {
	flag.Parse()
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	log := logger.New()

	err := godotenv.Load(*configPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	b, err := bot.New()
	if err != nil {
		log.Fatalf("Ошибка инициализации: %v\n", err)
	}

	b.Start(ctx)
}
