package ai

import (
	"context"
	"errors"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"tg-botv1/internal/logger"
)

var (
	ErrNoResponse = errors.New("no choices found")
)

var instructions = "Твоя задача смешно подстебывать криптокидов насчет их дебильных комментариев, можешь использовать свободный стиль и изощряться как хочешь"

type Client struct {
	ai  *openai.Client
	log *logrus.Logger
}

func New() (*Client, error) {
	client := openai.NewClient(os.Getenv("AI_KEY"))

	log := logger.New()

	return &Client{ai: client, log: log}, nil
}

func (c *Client) ReviewComments(comments []string) (string, error) {
	c.log.Infof("Starting review of %d comments", len(comments))

	commentsText := fmt.Sprintf("Комментари: %s", strings.Join(comments, ","))
	c.log.Debugf("Generated comments text for AI: %s", commentsText)

	resp, err := c.ai.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: instructions,
				},
				{Role: openai.ChatMessageRoleUser,
					Content: commentsText},
			},
		},
	)

	if err != nil {
		c.log.Errorf("Error when creating chat completion: %v", err)

		return "", err
	}

	if len(resp.Choices) == 0 {
		c.log.Warn("No response choices received from AI")

		return "", ErrNoResponse
	}

	c.log.Info("Received response from AI successfully")

	return resp.Choices[0].Message.Content, nil
}
