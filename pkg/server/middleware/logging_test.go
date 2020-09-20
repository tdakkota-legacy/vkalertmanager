package middleware

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

type handler struct{}

func (h handler) ServeHTTP(http.ResponseWriter, *http.Request) {
}

func TestZeroLog(t *testing.T) {
	counter := int64(0)

	hook := zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, message string) {
		atomic.AddInt64(&counter, 1)
	})
	l := zerolog.New(nil).Hook(hook)
	h := ZeroLog(l)(handler{})

	s := httptest.NewServer(h)
	defer s.Close()

	a := require.New(t)
	_, err := s.Client().Get(s.URL)
	a.NoError(err)
	a.Equal(int64(1), counter)
}
