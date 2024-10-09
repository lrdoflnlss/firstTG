package bot

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"strings"
	"tg-botv1/internal/gaspump"
	"time"
)

const TG_KEY = "7720831150:AAGLDfKOHarhRpKP73xakKkddbDKYzj6S4A"

type Bot struct {
	bot *telebot.Bot
	gp  *gaspump.Client
}

func New() (*Bot, error) {
	pref := telebot.Settings{
		Token:  TG_KEY,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}

	bot := &Bot{bot: b}

	b.Handle("/comment", bot.SendComments)

	return bot, nil
}

func (b *Bot) Start() {
	b.bot.Start()
}

func (b *Bot) CommentCA(c telebot.Context, ca string) ([]string, error) {
	cm, err := b.gp.GetComments(ca)
	if err != nil {
		return nil, fmt.Errorf("Не удалось получить комментарии: %v", err)
	}

	return cm, nil
}

func (b *Bot) ParseArgs(c telebot.Context) (string, error) {
	args := c.Args()
	if len(args) < 1 {
		return "", c.Send("Пожалуйста, укажите адрес смарт-контракта")
	}

	contractAdress := args[0]
	if !strings.Contains(contractAdress, "EQ") {
		return "", c.Send("Укажите правильный адрес смарт-контракта")
	}

	return contractAdress, nil

}

func (b *Bot) SendComments(c telebot.Context) error {
	contractAdress, err := b.ParseArgs(c)
	if err != nil {
		return err
	}

	comments, err := b.CommentCA(c, contractAdress)
	if err != nil {
		return c.Send(err)
	}

	for _, comment := range comments {
		if err := c.Send(comment); err != nil {
			return fmt.Errorf("Ошибка при отправке сообщения: %v", err)
		}
	}

	return nil
}
