package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shanisharrma/gopher-social/cmd/api/config"
	"github.com/shanisharrma/gopher-social/internal/auth"
	"github.com/shanisharrma/gopher-social/internal/ratelimiter"
	"github.com/shanisharrma/gopher-social/internal/store"
	"github.com/shanisharrma/gopher-social/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, cfg config.Config) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()

	mockStore := store.NewMockStore()
	mockCacheStorage := cache.NewMockStore()

	testAuth := &auth.TestAuthenticator{}

	ratelimiter := ratelimiter.NewFixedWindowRateLimiter(
		cfg.Ratelimiter.RequestPerTimeFrame,
		cfg.Ratelimiter.TimeFrame,
	)

	return &application{
		config:        cfg,
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockCacheStorage,
		authenticator: testAuth,
		ratelimiter:   ratelimiter,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected response code %d. got %d", expected, actual)
	}
}
