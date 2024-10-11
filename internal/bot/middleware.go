package bot

import (
	"fmt"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

func Logger(log *logrus.Logger) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			data := fmt.Sprintf("%v: %v", c.Sender().Username, c.Message().Text)
			log.Info(data)

			return next(c)
		}
	}
}

func OnErrorWithLogger(log *logrus.Logger) func(err error, c tele.Context) {
	return func(err error, c tele.Context) {
		log.Errorf("[ERROR]: %v", err)
		_ = c.Send("Произошла ошибка")
	}
}
