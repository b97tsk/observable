package rx

import (
	"context"
)

type firstOperator struct {
	source Operator
}

func (op firstOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	var mutableObserver Observer

	mutableObserver = func(t Notification) {
		switch {
		case t.HasValue:
			mutableObserver = NopObserver
			ob.Next(t.Value)
			ob.Complete()
			cancel()
		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()
		default:
			ob.Error(ErrEmpty)
			cancel()
		}
	}

	op.source.Call(ctx, func(t Notification) { t.Observe(mutableObserver) })

	return ctx, cancel
}

// First creates an Observable that emits only the first value (or the first
// value that meets some condition) emitted by the source Observable.
func (o Observable) First() Observable {
	op := firstOperator{o.Op}
	return Observable{op}
}
