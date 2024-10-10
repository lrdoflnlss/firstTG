package bot

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"strings"
	"tg-botv1/internal/ai"
	"tg-botv1/internal/gaspump"
	"time"
)

const TG_KEY = "7720831150:AAGLDfKOHarhRpKP73xakKkddbDKYzj6S4A"

type Bot struct {
	bot *telebot.Bot
	gp  *gaspump.Client
	ai  *ai.Client
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

	aiClient, err := ai.New()
	if err != nil {
		return nil, err
	}

	bot := &Bot{bot: b, ai: aiClient}

	b.Handle("/comment", bot.SendComments)
	b.Handle("/review", bot.HandleReviewCA)

	return bot, nil
}

func (b *Bot) Start() {
	b.bot.Start()
}

func (b *Bot) ParseArgsComments(c telebot.Context) (string, error) {
	args := c.Args()
	if len(args) < 1 {
		return "", c.Send("Пожалуйста, укажите адрес смарт-контракта")
	}

	contractAddress := args[0]
	if !strings.HasPrefix(contractAddress, "EQ") {
		return "", c.Send("Укажите правильный адрес смарт-контракта")
	}

	return contractAddress, nil

}

func (b *Bot) SendComments(c telebot.Context) error {
	contractAddress, err := b.ParseArgsComments(c)
	if err != nil {
		return err
	}

	comments, err := b.gp.GetComments(contractAddress)
	if err != nil {
		return fmt.Errorf("Произошла ошибка при получении комментариев: %v", err)
	}

	var sb strings.Builder
	sb.WriteString("Все комментарии: \n\n")
	for _, comment := range comments {
		sb.WriteString(comment)
	}

	return c.Send(sb.String())
}

func (b *Bot) HandleReviewCA(c telebot.Context) error {
	contractAddress, err := b.ParseArgsComments(c)
	if err != nil {
		return err
	}

	comments, err := b.gp.GetComments(contractAddress)
	if err != nil {
		return fmt.Errorf("Произошла ошибка при получении комментариев: %v", err)
	}

	aiResp, err := b.ai.ReviewComments(comments)
	if err != nil {
		return err
	}

	return c.Send(aiResp)
}
