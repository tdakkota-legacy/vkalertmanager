package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tdakkota/vkalertmanager/pkg/hook"
)

func prepareMessage(a *require.Assertions) (hook.Message, []byte) {
	m := hook.Message{}
	data, err := json.Marshal(m)
	a.NoError(err)
	return m, data
}

func getArgs(host, apiUrl string) []string {
	return []string{
		"app",
		"run",
		"--bind", host,
		"--tokens", "test_token",
		"--receivers", "1",
		"--vk.server", apiUrl,
	}
}

func TestApp_run(t *testing.T) {
	a := require.New(t)

	ctxt, cancel := context.WithCancel(context.Background())
	_, data := prepareMessage(a)

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response": 1}`))
		cancel()
	}))
	defer s.Close()

	host := "localhost:31512"
	url := "http://" + host + "/"
	args := getArgs(host, s.URL+"/")

	app := NewApp()
	go func() {
		_ = app.RunContext(ctxt, args)
	}()

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	a.NoError(err)
	_, err = http.DefaultClient.Do(req)
	a.NoError(err)

	<-ctxt.Done()
	err = ctxt.Err()
	if !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
