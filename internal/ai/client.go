package ai

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"strings"
)

const API_KEY = "4&L^%JbNxw^D9:C"

var instructions = "Твоя задача смешно подстебывать криптокидов насчет их дебильных комментариев, можешь использовать свободный стиль и изощряться как хочешь"

type Client struct {
	ai *openai.Client
}

func New() (*Client, error) {
	client := openai.NewClient(API_KEY)

	return &Client{ai: client}, nil
}

func (c *Client) ReviewComments(comments []string) (string, error) {
	resp, err := c.ai.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: instructions,
				},
				{Role: openai.ChatMessageRoleUser,
					Content: strings.Join(comments, " ")}, //или лучше через builder
			},
		},
	)

	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("No choices found\n")
	}

	return resp.Choices[0].Message.Content, nil
}
