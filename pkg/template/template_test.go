package template

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tdakkota/vkalertmanager/pkg/hook"
)

func TestDefault(t *testing.T) {
	alerts := hook.Alerts{
		{
			Status:   "firing",
			StartsAt: time.Now(),
			Labels: hook.KV{
				"alertname": "test-alert-label",
			},
			EndsAt: time.Now().Add(time.Second),
		},
	}

	m, err := Default().Produce(hook.Message{
		Data: hook.Data{
			Alerts: alerts,
		},
	})

	require.NoError(t, err)
	require.Contains(t, m, "test-alert-label")
	require.Contains(t, m, "Duration: 1 second")
}
