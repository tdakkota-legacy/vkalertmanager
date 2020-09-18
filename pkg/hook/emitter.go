package hook

import (
	"context"
)

type Emitter interface {
	Emit(ctxt context.Context, m Message) error
}

type EmitterFunc func(ctxt context.Context, m Message) error

func (e EmitterFunc) Emit(ctxt context.Context, m Message) error {
	return e(ctxt, m)
}
