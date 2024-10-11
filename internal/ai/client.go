package ai

import (
	"context"
	"errors"
	"github.com/sashabaranov/go-openai"
	"os"
	"strings"
)

var (
	ErrNoResponse = errors.New("no choices found")
)

var instructions = "Твоя задача смешно подстебывать криптокидов насчет их дебильных комментариев, можешь использовать свободный стиль и изощряться как хочешь"

type Client struct {
	ai *openai.Client
}

func New() (*Client, error) {
	client := openai.NewClient(os.Getenv("AI_KEY"))

	return &Client{ai: client}, nil
}

func (c *Client) ReviewComments(comments []string) (string, error) {
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
					Content: strings.Join(comments, " ")},
			},
		},
	)

	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", ErrNoResponse
	}

	return resp.Choices[0].Message.Content, nil
}
