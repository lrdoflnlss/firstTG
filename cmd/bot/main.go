package main

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"tg-botv1/internal/bot"
)

var configPath = flag.String("c", ".env", "Path to config")

func main() {
	flag.Parse()

	err := godotenv.Load(*configPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	b, err := bot.New()
	if err != nil {
		log.Fatalf("Ошибка инициализации: %v\n", err)
	}

	b.Start()
}
