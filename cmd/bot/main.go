package main

import (
	"log"
	"tg-botv1/internal/bot"
)

func main() {
	b, err := bot.New()
	if err != nil {
		log.Fatalf("Ошибка инициализации: %v\n", err)
	}

	b.Start()
}
