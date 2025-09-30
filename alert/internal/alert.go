package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Notifier interface {
	Notify(ctx context.Context, msg string) error
}

type Discord struct {
	URL    string
	client *http.Client
}

func NewDiscord(url string, timeout time.Duration) *Discord {
	return &Discord{
		URL:    url,
		client: &http.Client{Timeout: timeout},
	}
}

func (discord *Discord) Notify(ctx context.Context, msg string) error {
	if discord.URL == "" {
		return nil
	}
	b, _ := json.Marshal(map[string]string{"content": msg})
	request, _ := http.NewRequestWithContext(ctx, http.MethodPost, discord.URL, bytes.NewReader(b))
	request.Header.Set("Content-Type", "application/json")
	response, err := discord.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}
