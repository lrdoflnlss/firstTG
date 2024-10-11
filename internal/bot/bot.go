package bot

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"os"
	"strings"
	"tg-botv1/internal/ai"
	"tg-botv1/internal/gaspump"
	"tg-botv1/internal/logger"
	"time"
)

type Bot struct {
	bot *telebot.Bot
	gp  *gaspump.Client
	ai  *ai.Client
	log *logrus.Logger
}

func New() (*Bot, error) {
	log := logger.New()

	pref := telebot.Settings{
		Token:   os.Getenv("TG_KEY"),
		Poller:  &telebot.LongPoller{Timeout: 10 * time.Second},
		OnError: OnErrorWithLogger(log),
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Errorf("failed to initialize bot: %v", err)

		return nil, err
	}
	log.Info("Bot initialized successfully")

	aiClient, err := ai.New()
	if err != nil {
		return nil, err
	}

	gpClient := gaspump.New()

	bot := &Bot{bot: b, ai: aiClient, log: log, gp: gpClient}

	b.Use(middleware.Recover())
	b.Use(Logger(log))
	b.Use(middleware.AutoRespond())

	b.Handle("/comment", bot.SendComments)
	b.Handle("/review", bot.HandleReviewCA)

	return bot, nil
}

func (b *Bot) Start() {
	b.log.Info("Bot started: ", b.bot.Me.Username)
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
		b.log.Info("error when receiving smart contract")

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
