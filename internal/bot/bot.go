package bot

import (
	"context"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"os"
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

func (b *Bot) Start(ctx context.Context) {
	go func() {
		<-ctx.Done()
		b.log.Infoln("Bot stopped")
		b.bot.Stop()
	}()

	b.log.Info("Bot started: ", b.bot.Me.Username)
	b.bot.Start()
}
