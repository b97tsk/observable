package rx

import (
	"context"
)

type bufferWhenOperator struct {
	ClosingSelector func() Observable
}

func (op bufferWhenOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	type X struct {
		Buffers []interface{}
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	var (
		openBuffer     func()
		avoidRecursive avoidRecursiveCalls
	)

	openBuffer = func() {
		if isDone(ctx) {
			return
		}

		ctx, cancel := context.WithCancel(ctx)

		var observer Observer
		observer = func(t Notification) {
			observer = NopObserver
			cancel()
			if x, ok := <-cx; ok {
				if t.HasError {
					close(cx)
					sink(t)
					return
				}
				sink.Next(x.Buffers)
				x.Buffers = nil
				cx <- x
				avoidRecursive.Do(openBuffer)
			}
		}

		closingNotifier := op.ClosingSelector()
		closingNotifier.Subscribe(ctx, observer.Notify)
	}

	avoidRecursive.Do(openBuffer)

	if isDone(ctx) {
		return Done()
	}

	source.Subscribe(ctx, func(t Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				x.Buffers = append(x.Buffers, t.Value)
				cx <- x
			default:
				close(cx)
				sink.Next(x.Buffers)
				sink(t)
			}
		}
	})

	return ctx, cancel
}

// BufferWhen buffers the source Observable values, using a factory function
// of closing Observables to determine when to close, emit, and reset the
// buffer.
//
// BufferWhen collects values from the past as a slice. When it starts
// collecting values, it calls a function that returns an Observable that
// tells when to close the buffer and restart collecting.
//
// Dead loop could happen if closing Observables emit a value or complete as
// soon as they are subscribed to.
func (Operators) BufferWhen(closingSelector func() Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := bufferWhenOperator{closingSelector}
		return source.Lift(op.Call)
	}
}
