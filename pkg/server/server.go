package server

import (
	"context"
	"net"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/tdakkota/vkalertmanager/pkg/hook"
	"golang.org/x/sync/errgroup"
)

func Create(emitter hook.Emitter, logger zerolog.Logger, config ...ListenerConfig) HookServer {
	h := hook.NewHook(emitter, logger.With().Str("type", "hook").Logger())
	return NewHookServer(h, logger.With().Str("type", "server").Logger(), config...)
}

type HookServer struct {
	hook   hook.Hook
	config []ListenerConfig

	logger zerolog.Logger
	server *http.Server
}

func NewHookServer(hook hook.Hook, logger zerolog.Logger, config ...ListenerConfig) HookServer {
	server := &http.Server{
		Handler: hook,
	}

	return HookServer{
		hook:   hook,
		config: config,
		logger: logger,
		server: server,
	}
}

func (h HookServer) listen(c ListenerConfig) error {
	h.logger.Info().Msgf("binding to %s", c.Bind)

	l, err := net.Listen("tcp", c.Bind)
	if err != nil {
		return err
	}

	if c.TLS != nil {
		tls := c.TLS
		return h.server.ServeTLS(l, tls.CertFile, tls.KeyFile)
	}

	return h.server.Serve(l)
}

func (h HookServer) Run(ctx context.Context) error {
	h.server.BaseContext = func(net.Listener) context.Context {
		return ctx
	}

	g, ctx := errgroup.WithContext(ctx)
	for i := range h.config {
		c := h.config[i]
		g.Go(func() error {
			err := h.listen(c)
			if err != nil {
				h.logger.Error().Err(err).Msgf("listener %s stopped", c.Bind)
			}
			return err
		})
	}

	return g.Wait()
}
