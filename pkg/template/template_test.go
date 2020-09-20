package template

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tdakkota/vkalertmanager/pkg/hook"
)

func TestDefault(t *testing.T) {
	a := require.New(t)

	alerts := hook.Alerts{
		{
			Status: "firing",
			Labels: hook.KV{
				"alertname": "test-alert-label",
			},
			Annotations: hook.KV{
				"message":     "message",
				"summary":     "summary",
				"description": "description",
			},
			StartsAt: time.Now(),
			EndsAt:   time.Now().Add(time.Second),
		},
	}

	m, err := Default().Produce(hook.Message{
		Data: hook.Data{
			Alerts: alerts,
		},
	})

	a.NoError(err)
	a.Contains(m, "test-alert-label")
	a.Contains(m, "Duration: 1 second")
}
