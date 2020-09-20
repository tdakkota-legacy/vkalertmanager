package middleware

import (
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

func ZeroLog(l zerolog.Logger) func(http.Handler) http.Handler {
	c := alice.New()

	return func(next http.Handler) http.Handler {
		// Install the logger handler
		c = c.Append(hlog.NewHandler(l))
		// Set access handler
		c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg("")
		}))
		// Install some provided extra handler to set some request's context fields.
		// Thanks to those handler, all our logs will come with some pre-populated fields.
		c = c.Append(hlog.RemoteAddrHandler("ip"))
		c = c.Append(hlog.UserAgentHandler("user_agent"))
		c = c.Append(hlog.RefererHandler("referer"))

		// Here is your final handler
		return c.Then(next)
	}
}
