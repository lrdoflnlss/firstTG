package bot

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log"
)

func Logger() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			data := fmt.Sprintf("%v: %v", c.Sender().Username, c.Message().Text)
			log.Println(data)

			return next(c)
		}
	}
}

func OnError(err error, c tele.Context) {
	log.Println("[ERROR]:", err)

	_ = c.Send("Произошла ошибка")
}
