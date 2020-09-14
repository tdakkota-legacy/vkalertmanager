package main

import (
	"errors"
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

var ErrAtLeastOneTokenExpected = errors.New("at least one token expected")
var ErrAtLeastOneReceiverExpected = errors.New("at least one alert receiver expected")

func parseTemplate(c *cli.Context) (emitter.Template, error) {
	if c.IsSet("hook.template.file") {
		return emitter.ParseFiles(c.Path("hook.template.file"))
	}

	return emitter.Parse(c.String("hook.template.body"))
}

func (app *App) createHook(c *cli.Context) (hook.Hook, error)  {
	t, err := parseTemplate(c)
	if err != nil {
		return hook.Hook{}, err
	}

	tokens := c.StringSlice("vk.tokens")
	if len(tokens) < 1 {
		return hook.Hook{}, ErrAtLeastOneTokenExpected
	}

	receivers := c.IntSlice("hook.receivers")
	if len(tokens) < 1 {
		return hook.Hook{}, ErrAtLeastOneReceiverExpected
	}

	vk := sdkutil.BuildSDK(tokens[0], tokens[1:]...).WithMiddleware(zlog.LoggingMiddleware(
		log.With().Str("type", "vksdk").Logger().Level(zerolog.DebugLevel),
	)).Complete()

	emit := emitter.NewVK(vk, receivers, t)
	h := hook.NewHook(emit, app.logger.With().Str("type", "hook").Logger())

	return h, nil
}

func (app *App) setup(c *cli.Context) error {
	app.logger = log.Logger

	h, err := app.createHook(c)
	if err != nil {
		return err
	}

	app.server = server.NewHookServer(h, app.logger, app.parseListenerConfig(c))
	return nil
}
