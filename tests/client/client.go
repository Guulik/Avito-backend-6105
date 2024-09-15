package client

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
	"zadanie-6105/config"
)

type Suite struct {
	*testing.T
	Cfg    *config.Config
	Client *http.Client
}

var (
	BaseURL = "http://localhost:8080/api"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadPath(configPath())

	ctx, cancelCtx := context.WithTimeout(context.Background(), 30*time.Second)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	return ctx, &Suite{
		T:      t,
		Cfg:    cfg,
		Client: &http.Client{},
	}
}

func FormRequest(
	method string,
	url string,
	body io.Reader,
) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil
	}
	req.Header.Add("Content-Type", "application/json")

	return req
}

func configPath() string {
	const key = "CONFIG_PATH"

	if v := os.Getenv(key); v != "" {
		return v
	}

	return "../../config/stage.yaml"
}
