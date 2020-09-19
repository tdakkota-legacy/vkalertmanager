package emitter

import (
	"bytes"
	"context"
	"strconv"
	"testing"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/stretchr/testify/require"
	"github.com/tdakkota/vkalertmanager/pkg/hook"
	"github.com/tdakkota/vkalertmanager/pkg/template"
	"github.com/tdakkota/vksdkutil/v2/testutil"
)

type IntSlice []int

func (s IntSlice) String() string {
	// max length from -9223372036854775808
	b := make([]byte, len(s)*20)
	cut := make([][]byte, len(s))

	for i := 0; i < len(s); i++ {
		cutIndex := i * 20
		cut[i] = b[cutIndex : cutIndex+20]
	}

	for i, v := range s {
		cut[i] = strconv.AppendInt(cut[i][:0], int64(v), 10)
	}

	return string(bytes.Join(cut, []byte(",")))
}

func TestVK_Emit(t *testing.T) {
	receivers := []int{1, 2, 3, 4, 5, 6, 7}

	t.Run("group", func(t *testing.T) {
		vk, cse := testutil.CreateSDK(t)
		defer cse.ExpectationsWereMet()

		emit := NewVK(vk, receivers)
		cse.ExpectCall("messages.send").WithParamsF(func() api.Params {
			return api.Params{
				"user_ids": IntSlice(receivers).String(),
			}
		})

		err := emit.Emit(context.Background(), hook.Message{})
		require.NoError(t, err)
	})

	t.Run("user", func(t *testing.T) {
		vk, cse := testutil.CreateSDK(t)
		defer cse.ExpectationsWereMet()

		emit := NewVK(vk, receivers,
			WithTemplate(template.Default()),
			WithIsUser(true),
		)

		for i := range receivers {
			cse.ExpectCall("messages.send").WithParamsF(func() api.Params {
				return api.Params{
					"peer_id": receivers[i],
				}
			})
		}

		err := emit.Emit(context.Background(), hook.Message{})
		require.NoError(t, err)
	})
}
