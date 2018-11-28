package rx

import (
	"context"
)

type fromSliceOperator struct {
	Slice []interface{}
}

func (op fromSliceOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	done := ctx.Done()
	for _, val := range op.Slice {
		select {
		case <-done:
			return canceledCtx, nothingToDo
		default:
		}
		sink.Next(val)
	}
	sink.Complete()
	return canceledCtx, nothingToDo
}

func just(one interface{}) Observable {
	return Observable{}.Lift(
		func(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
			sink.Next(one)
			sink.Complete()
			return canceledCtx, nothingToDo
		},
	)
}

// FromSlice creates an Observable that emits values from a slice, one after
// the other, and then completes.
func FromSlice(slice []interface{}) Observable {
	switch {
	case len(slice) > 1:
		op := fromSliceOperator{slice}
		return Observable{}.Lift(op.Call)
	case len(slice) == 1:
		return just(slice[0])
	default:
		return Empty()
	}
}

// Just creates an Observable that emits some values you specify as arguments,
// one after the other, and then completes.
func Just(values ...interface{}) Observable {
	return FromSlice(values)
}
