package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tdakkota/vkalertmanager/pkg/emitter"
	"github.com/tdakkota/vkalertmanager/pkg/hook"
	"github.com/tdakkota/vkalertmanager/pkg/server"
	sdkutil "github.com/tdakkota/vksdkutil/v2"
	zlog "github.com/tdakkota/vksdkutil/v2/middleware/log/zerolog"
	"github.com/urfave/cli/v2"
)

func (app *App) parseListenerConfig(c *cli.Context) server.ListenerConfig {
	var tls *server.TLSConfig
	if c.IsSet("server.tls.key_file") && c.IsSet("server.tls.cert_file") {
		tls = &server.TLSConfig{
			KeyFile:  c.String("server.tls.key_file"),
			CertFile: c.String("server.tls.cert_file"),
		}
	}

	config := server.ListenerConfig{
		Bind: c.String("server.bind"),
		TLS:  tls,
	}

	return config
}

func (app *App) setup(c *cli.Context) error {
	app.logger = log.Logger

	t, err := emitter.Parse(c.String(""))
	if err != nil {
		return err
	}

	var receivers []int
	vk := sdkutil.BuildSDK("token").WithMiddleware(zlog.LoggingMiddleware(
		log.With().Str("type", "vksdk").Logger().Level(zerolog.DebugLevel),
	)).Complete()

	emit := emitter.NewVK(vk, receivers, t)
	h := hook.NewHook(emit, app.logger.With().Str("type", "hook").Logger())

	app.server = server.NewHookServer(h, app.logger, app.parseListenerConfig(c))

	return nil
}
