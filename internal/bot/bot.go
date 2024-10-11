package bot

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"log"
	"os"
	"strings"
	"tg-botv1/internal/ai"
	"tg-botv1/internal/gaspump"
	"time"
)

type Bot struct {
	bot *telebot.Bot
	gp  *gaspump.Client
	ai  *ai.Client
}

func New() (*Bot, error) {
	pref := telebot.Settings{
		Token:   os.Getenv("TG_KEY"),
		Poller:  &telebot.LongPoller{Timeout: 10 * time.Second},
		OnError: OnError,
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

	b.Use(middleware.Recover())
	b.Use(Logger())
	b.Use(middleware.AutoRespond())

	b.Handle("/comment", bot.SendComments)
	b.Handle("/review", bot.HandleReviewCA)

	return bot, nil
}

func (b *Bot) Start() {
	log.Println("Bot started:", b.bot.Me.Username)
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
		return fmt.Errorf("произошла ошибка при получении комментариев: %v", err)
	}

	var sb strings.Builder
	sb.WriteString("Все комментарии: ")
	for _, comment := range comments {
		sb.WriteString("\n\n")
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
		return fmt.Errorf("произошла ошибка при получении комментариев: %v", err)
	}

	aiResp, err := b.ai.ReviewComments(comments)
	if err != nil {
		return err
	}

	return c.Send(aiResp)
}
