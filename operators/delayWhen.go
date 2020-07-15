package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type delayWhenObservable struct {
	Source           rx.Observable
	DurationSelector func(interface{}, int) rx.Observable
}

func (obs delayWhenObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)
	sink = sink.WithCancel(cancel)

	type X struct {
		Index           int
		Active          int
		SourceCompleted bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				sourceIndex := x.Index
				sourceValue := t.Value
				x.Index++
				x.Active++

				cx <- x

				scheduleCtx, scheduleCancel := context.WithCancel(ctx)

				var observer rx.Observer
				observer = func(t rx.Notification) {
					observer = rx.Noop
					scheduleCancel()
					if x, ok := <-cx; ok {
						x.Active--
						switch {
						case t.HasValue:
							sink.Next(sourceValue)
							if x.Active == 0 && x.SourceCompleted {
								close(cx)
								sink.Complete()
								return
							}
							cx <- x
						case t.HasError:
							close(cx)
							sink(t)
						default:
							if x.Active == 0 && x.SourceCompleted {
								close(cx)
								sink(t)
								return
							}
							cx <- x
						}
					}
				}

				obs := obs.DurationSelector(sourceValue, sourceIndex)
				obs.Subscribe(scheduleCtx, observer.Sink)

			case t.HasError:
				close(cx)
				sink(t)

			default:
				x.SourceCompleted = true
				if x.Active == 0 {
					close(cx)
					sink(t)
					return
				}
				cx <- x
			}
		}
	})
}

// DelayWhen creates an Observable that delays the emission of items from
// the source Observable by a given time span determined by the emissions of
// another Observable.
//
// It's like Delay, but the time span of the delay duration is determined by
// a second Observable.
func DelayWhen(durationSelector func(interface{}, int) rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return delayWhenObservable{source, durationSelector}.Subscribe
	}
}
