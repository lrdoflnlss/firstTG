package bot

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"strings"
)

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
