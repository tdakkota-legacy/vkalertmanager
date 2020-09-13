package hook

import (
	"encoding/json"
	"github.com/tdakkota/vkalertmanager/pkg/emitter"
	"net/http"

	"github.com/rs/zerolog"
)

type Hook struct {
	emitter emitter.Emitter
	logger  zerolog.Logger
}

func NewHook(emitter emitter.Emitter, logger zerolog.Logger) Hook {
	return Hook{emitter: emitter, logger: logger}
}

func (h Hook) decode(r *http.Request) (m Message, err error) {
	err = json.NewDecoder(r.Body).Decode(&m)
	return
}

func (h Hook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	m, err := h.decode(r)
	if err != nil {
		h.logger.Info().Err(err).Msg("failed to unmarshal message")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.emitter.Emit(r.Context(), m)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to emit alert event")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
