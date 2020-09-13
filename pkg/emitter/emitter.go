package emitter

import (
	"context"

	"github.com/tdakkota/vkalertmanager/pkg/hook"
)

type Emitter interface {
	Emit(ctxt context.Context, m hook.Message) error
}
