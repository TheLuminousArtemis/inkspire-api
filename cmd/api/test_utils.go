package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theluminousartemis/socialnews/internal/auth"
	"github.com/theluminousartemis/socialnews/internal/ratelimiter"
	"github.com/theluminousartemis/socialnews/internal/store"
	"github.com/theluminousartemis/socialnews/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()
	logger := zap.NewNop().Sugar()
	mockStore := store.NewMockStore()
	mockCache := cache.NewMockStore()
	auth := &auth.TestAuthenticator{}
	ratelimiter := ratelimiter.NewFixedWindowRateLimiter(
		cfg.ratelimiter.RequestsPerTimeFrame,
		cfg.ratelimiter.Timeframe,
	)
	return &application{
		l:             logger,
		storage:       mockStore,
		cache:         mockCache,
		config:        cfg,
		authenticator: auth,
		rateLimiter:   ratelimiter,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected response code %d got %d", expected, actual)
	}
}
