package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/theluminousartemis/inkspire/internal/ratelimiter"
)

func TestRateLimiterMiddleware(t *testing.T) {
	cfg := config{
		ratelimiter: ratelimiter.Config{
			RequestsPerTimeFrame: 20,
			Timeframe:            time.Second,
			Enabled:              true,
		},
		addr: ":8080",
	}

	app := newTestApplication(t, cfg)
	ts := httptest.NewServer(app.mount())
	defer ts.Close()

	client := &http.Client{}
	mockIP := "192.168.1.1"
	marginOfError := 2

	for i := 0; i < cfg.ratelimiter.RequestsPerTimeFrame+marginOfError; i++ {
		req, err := http.NewRequest("GET", ts.URL+"/v1/health", nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		req.Header.Set("X-Forwarded-For", mockIP)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not send request: %v", err)
		}
		defer resp.Body.Close()
		if i < cfg.ratelimiter.RequestsPerTimeFrame {
			checkResponseCode(t, http.StatusOK, resp.StatusCode)
		} else {
			checkResponseCode(t, http.StatusTooManyRequests, resp.StatusCode)
		}
	}
}
