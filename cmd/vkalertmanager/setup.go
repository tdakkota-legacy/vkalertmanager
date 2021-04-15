package main

import (
	"errors"

	"github.com/tdakkota/vkalertmanager/pkg/template"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	sdkutil "github.com/tdakkota/vksdkutil/v2"
	zlog "github.com/tdakkota/vksdkutil/v2/middleware/log/zerolog"
	"github.com/urfave/cli/v2"

	"github.com/tdakkota/vkalertmanager/pkg/emitter"
	"github.com/tdakkota/vkalertmanager/pkg/hook"
	"github.com/tdakkota/vkalertmanager/pkg/server"
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

func parseTemplate(c *cli.Context) (*template.Template, error) {
	if c.IsSet("hook.template.file") {
		return template.ParseFiles(c.Path("hook.template.file"))
	}

	if c.IsSet("hook.template.body") {
		return template.Parse(c.String("hook.template.body"))
	}

	return template.Default(), nil
}

func (app *App) createEmitter(c *cli.Context) (hook.Emitter, error) {
	t, err := parseTemplate(c)
	if err != nil {
		return nil, err
	}

	tokens := c.StringSlice("vk.tokens")
	if len(tokens) < 1 {
		return nil, ErrAtLeastOneTokenExpected
	}

	receivers := c.IntSlice("hook.receivers")
	if len(tokens) < 1 {
		return nil, ErrAtLeastOneReceiverExpected
	}

	vk := sdkutil.BuildSDK(tokens[0], tokens[1:]...).
		WithMiddleware(zlog.LoggingMiddleware(
			app.logger.With().
				Str("type", "vksdk").
				Logger().
				Level(zerolog.DebugLevel),
		)).
		WithUserAgent(c.String("vk.user_agent")).
		WithMethodURL(c.String("vk.server")).
		Complete()

	return emitter.NewVK(vk, receivers, emitter.WithTemplate(t)), nil
}

func (app *App) setup(c *cli.Context) error {
	app.logger = log.Logger

	emit, err := app.createEmitter(c)
	if err != nil {
		return err
	}

	app.server = server.Create(emit, app.logger, app.parseListenerConfig(c))
	return nil
}
