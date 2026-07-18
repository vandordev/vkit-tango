package contract

import "context"

type Command[I any, O any] interface {
	Execute(context.Context, I) (O, error)
}
