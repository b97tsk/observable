package rx

import (
	"context"
	"sync"
)

type switchMapOperator struct {
	project func(interface{}, int) Observable
}

func (op switchMapOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()
	childCtx, childCancel := canceledCtx, doNothing

	var (
		mu             sync.Mutex
		outerIndex     = -1
		activeIndex    = -1
		completeSignal = make(chan struct{}, 1)
	)

	source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			mu.Lock()
			defer mu.Unlock()

			outerValue := t.Value
			outerIndex++
			outerIndex := outerIndex

			activeIndex = outerIndex
			childCancel()

			obsv := op.project(outerValue, outerIndex)

			childCtx, childCancel = context.WithCancel(ctx)

			go obsv.Subscribe(childCtx, func(t Notification) {
				switch {
				case t.HasValue:
					sink(t)

				case t.HasError:
					sink(t)
					cancel()

				default:
					mu.Lock()

					if activeIndex != outerIndex {
						mu.Unlock()
						break
					}

					activeIndex = -1
					mu.Unlock()

					select {
					case completeSignal <- struct{}{}:
					default:
					}
				}
			})

		case t.HasError:
			sink(t)
			cancel()

		default:
			mu.Lock()
			if activeIndex != -1 {
				go func() {
					for activeIndex != -1 {
						mu.Unlock()
						select {
						case <-done:
							return
						case <-completeSignal:
						}
						mu.Lock()
					}
					mu.Unlock()
					sink.Complete()
					cancel()
				}()
				return
			}
			mu.Unlock()
			sink(t)
			cancel()
		}
	})

	return ctx, cancel
}

// Switch converts a higher-order Observable into a first-order Observable by
// subscribing to only the most recently emitted of those inner Observables.
//
// Switch flattens an Observable-of-Observables by dropping the previous inner
// Observable once a new one appears.
func (o Observable) Switch() Observable {
	op := switchMapOperator{projectToObservable}
	return o.Lift(op.Call).Mutex()
}

// SwitchMap creates an Observable that projects each source value to an
// Observable which is merged in the output Observable, emitting values only
// from the most recently projected Observable.
//
// SwitchMap maps each value to an Observable, then flattens all of these inner
// Observables using Switch.
func (o Observable) SwitchMap(project func(interface{}, int) Observable) Observable {
	op := switchMapOperator{project}
	return o.Lift(op.Call).Mutex()
}

// SwitchMapTo creates an Observable that projects each source value to the
// same Observable which is flattened multiple times with Switch in the output
// Observable.
//
// It's like SwitchMap, but maps each value always to the same inner Observable.
func (o Observable) SwitchMapTo(inner Observable) Observable {
	return o.SwitchMap(func(interface{}, int) Observable { return inner })
}
