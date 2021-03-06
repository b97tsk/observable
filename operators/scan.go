package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Scan applies an accumulator function over the source, and emits each
// intermediate result, with an optional seed value.
//
// It's like Reduce, but emits the current accumulation whenever the source
// emits a value.
func Scan(accumulator func(interface{}, interface{}, int) interface{}) rx.Operator {
	return ScanConfigure{Accumulator: accumulator}.Make()
}

// A ScanConfigure is a configure for Scan.
type ScanConfigure struct {
	Accumulator func(interface{}, interface{}, int) interface{}
	Seed        interface{}
	HasSeed     bool
}

// Make creates an Operator from this configure.
func (configure ScanConfigure) Make() rx.Operator {
	if configure.Accumulator == nil {
		panic("Scan: Accumulator is nil")
	}

	return func(source rx.Observable) rx.Observable {
		return scanObservable{source, configure}.Subscribe
	}
}

type scanObservable struct {
	Source rx.Observable
	ScanConfigure
}

func (obs scanObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	sourceIndex := -1

	acc := struct {
		Value    interface{}
		HasValue bool
	}{obs.Seed, obs.HasSeed}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			sourceIndex++

			if acc.HasValue {
				acc.Value = obs.Accumulator(acc.Value, t.Value, sourceIndex)
			} else {
				acc.Value = t.Value
				acc.HasValue = true
			}

			sink.Next(acc.Value)

		default:
			sink(t)
		}
	})
}
