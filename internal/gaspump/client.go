package gaspump

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"tg-botv1/internal/logger"
)

type Client struct {
	log *logrus.Logger
}

type CommentData []struct {
	Text string `json:"text"`
}

func New() *Client {
	log := logger.New()
	return &Client{log: log}
}

func (c *Client) GetComments(ca string) ([]string, error) {
	url := fmt.Sprintf("https://api.gas111.com/api/v1/community-notes/list?token_address=%s", ca)
	method := "GET"

	c.log.Infof("Starting HTTP request to %s", url)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		c.log.Errorf("Error when creating HTTP request: %v", err)

		return nil, err
	}
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		c.log.Errorf("Error when making HTTP request: %v", err)

		return nil, err
	}

	defer res.Body.Close()
	c.log.Infof("Received HTTP response with status code: %d", res.StatusCode)

	if res.StatusCode != http.StatusOK {
		c.log.Errorf("Unexpected status code: %d", res.StatusCode)

		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.log.Errorf("Error when reading response body: %v", err)

		return nil, err
	}

	var commentData CommentData
	err = json.Unmarshal(body, &commentData)
	if err != nil {
		c.log.Errorf("Error when unmarshaling response body: %v", err)

		return nil, err
	}

	c.log.Infof("Successfully parsed comments, total comments: %d", len(commentData))

	var comments []string
	for _, comment := range commentData {
		comments = append(comments, comment.Text)
	}

	c.log.Infof("Returning %d comments", len(comments))

	return comments, nil
}
