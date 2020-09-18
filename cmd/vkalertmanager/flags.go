package main

import (
	"github.com/SevereCloud/vksdk/v2"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func (app *App) getEnvNames(names ...string) []string {
	return names
}

func (app *App) addFileConfig(flagName string, command *cli.Command) {
	prev := command.Before

	command.Before = func(context *cli.Context) error {
		if prev != nil {
			err := prev(context)
			if err != nil {
				return err
			}
		}

		path := context.Path(flagName)
		fileContext, err := altsrc.NewYamlSourceFromFile(path)
		if err != nil {
			app.logger.Warn().Err(err).Msgf("failed to load config from %s", path)
			return nil
		}

		return altsrc.ApplyInputSourceValues(context, fileContext, command.Flags)
	}
}

func (app *App) flags() []cli.Flag {
	vkFlags := []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "vk.server",
			Required: false,
			Value:    "https://api.vk.com/method/",
			Usage:    "VK API Method URL",
			EnvVars:  app.getEnvNames("VK_SERVER"),
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "vk.user_agent",
			Required: false,
			Value:    "vksdk/" + vksdk.Version + " (+https://github.com/SevereCloud/vksdk)",
			Usage:    "VK API Client useragent",
			EnvVars:  app.getEnvNames("VK_USER_AGENT"),
		}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name:     "vk.tokens",
			Aliases:  []string{"tokens"},
			Required: true,
			Usage:    "VK API tokens",
		}),
	}
	hookFlags := []cli.Flag{
		altsrc.NewIntSliceFlag(&cli.IntSliceFlag{
			Name:     "hook.receivers",
			Aliases:  []string{"receivers"},
			Required: true,
			Usage:    "Sets receiver users.",
			EnvVars:  app.getEnvNames("RECEIVERS"),
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "hook.template.body",
			Aliases:  []string{"template"},
			Required: false,
			Usage:    "Message template body.",
		}),
		altsrc.NewPathFlag(&cli.PathFlag{
			Name:     "hook.template.file",
			Aliases:  []string{"template-file"},
			Required: false,
			Usage:    "Message template file.",
		}),
	}
	serverFlags := []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "server.bind",
			Aliases:  []string{"bind"},
			Required: false,
			Value:    ":8080",
			Usage:    "Sets bind address",
			EnvVars:  app.getEnvNames("ADDR"),
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "server.tls.key_file",
			Required: false,
			Usage:    "TLS key file.",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "server.tls.cert_file",
			Required: false,
			Usage:    "TLS cert file.",
		}),
	}

	flags := []cli.Flag{
		&cli.PathFlag{
			Name:    "config.file",
			Value:   "vkalertmanager.yml",
			Usage:   "path to config file",
			EnvVars: app.getEnvNames("CONFIG_FILE", "CONFIG"),
		},

		// logging
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "logging.format",
			Value:   "human",
			Usage:   "logging format(json/human)",
			EnvVars: app.getEnvNames("LOG_FORMAT"),
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "logging.level",
			Value:   "debug",
			Usage:   "logging level",
			EnvVars: app.getEnvNames("LOG_LEVEL"),
		}),
	}
	flags = append(flags, vkFlags...)
	flags = append(flags, serverFlags...)
	flags = append(flags, hookFlags...)
	return flags
}
