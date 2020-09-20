package hook

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func prepareTestMessage() (Message, []byte, error) {
	testMessage := Message{
		Data:            Data{},
		Version:         "4",
		GroupKey:        "",
		TruncatedAlerts: 0,
	}

	data, err := json.Marshal(testMessage)
	if err != nil {
		return Message{}, nil, err
	}

	return testMessage, data, nil
}

func testHook(e EmitterFunc, cb func(*httptest.Server)) {
	hook := NewHook(e, log.Logger)

	s := httptest.NewServer(hook)
	defer s.Close()
	cb(s)
}

var testError = errors.New("test-error")

func TestHook(t *testing.T) {
	testMessage, body, err := prepareTestMessage()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("ok", func(t *testing.T) {
		a := require.New(t)

		e := EmitterFunc(func(ctxt context.Context, m Message) error {
			a.Equal(testMessage, m)
			return nil
		})

		testHook(e, func(s *httptest.Server) {
			r, err := s.Client().Post(
				s.URL,
				"application/json",
				bytes.NewReader(body),
			)
			a.NoError(err)
			a.Equal(http.StatusOK, r.StatusCode)
		})
	})

	t.Run("bad-emit", func(t *testing.T) {
		a := require.New(t)

		e := EmitterFunc(func(ctxt context.Context, m Message) error {
			a.Equal(testMessage, m)
			return testError
		})

		testHook(e, func(s *httptest.Server) {
			r, err := s.Client().Post(
				s.URL,
				"application/json",
				bytes.NewReader(body),
			)
			a.NoError(err)
			a.Equal(http.StatusInternalServerError, r.StatusCode)
		})
	})

	t.Run("bad-method", func(t *testing.T) {
		a := require.New(t)

		testHook(nil, func(s *httptest.Server) {
			r, err := s.Client().Get(s.URL)
			a.NoError(err)
			a.Equal(http.StatusMethodNotAllowed, r.StatusCode)
		})
	})

	t.Run("bad-request", func(t *testing.T) {
		a := require.New(t)

		testHook(nil, func(s *httptest.Server) {
			r, err := s.Client().Post(
				s.URL,
				"application/json",
				bytes.NewBufferString(""),
			)
			a.NoError(err)
			a.Equal(http.StatusBadRequest, r.StatusCode)
		})
	})
}
