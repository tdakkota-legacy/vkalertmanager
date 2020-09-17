package hook

import (
	"context"
)

type Emitter interface {
	Emit(ctxt context.Context, m Message) error
}
