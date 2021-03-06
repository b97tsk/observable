package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestShare1(t *testing.T) {
	obs := rx.Ticker(Step(3)).Pipe(
		operators.Map(
			func(val interface{}, idx int) interface{} {
				return idx
			},
		),
		operators.Take(4),
		operators.Share(rx.Multicast),
	)

	NewTestSuite(t).Case(
		rx.Merge(
			obs,
			obs.Pipe(DelaySubscription(4)),
			obs.Pipe(DelaySubscription(8)),
			obs.Pipe(DelaySubscription(13)),
		),
		0, 1, 1, 2, 2, 2, 3, 3, 3, 0, 1, 2, 3, Completed,
	).TestAll()
}

func TestShare2(t *testing.T) {
	obs := rx.Ticker(Step(3)).Pipe(
		operators.Map(
			func(val interface{}, idx int) interface{} {
				return idx
			},
		),
		operators.Share(rx.Multicast),
		operators.Take(4),
	)

	NewTestSuite(t).Case(
		rx.Merge(
			obs,
			obs.Pipe(DelaySubscription(4)),
			obs.Pipe(DelaySubscription(8)),
			obs.Pipe(DelaySubscription(19)),
		),
		0, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 5, 0, 1, 2, 3, Completed,
	).TestAll()
}

func TestShare3(t *testing.T) {
	obs := rx.Ticker(Step(3)).Pipe(
		operators.Map(
			func(val interface{}, idx int) interface{} {
				return idx
			},
		),
		operators.Take(4),
		operators.Share(
			rx.MulticastReplayFactory(
				&rx.ReplayOptions{BufferSize: 1},
			),
		),
	)

	NewTestSuite(t).Case(
		rx.Merge(
			obs,
			obs.Pipe(DelaySubscription(4)),
			obs.Pipe(DelaySubscription(8)),
			obs.Pipe(DelaySubscription(13)),
		),
		0, 0, 1, 1, 1, 2, 2, 2, 3, 3, 3, 0, 1, 2, 3, Completed,
	).TestAll()
}

func TestShare4(t *testing.T) {
	obs := rx.Ticker(Step(3)).Pipe(
		operators.Map(
			func(val interface{}, idx int) interface{} {
				return idx
			},
		),
		operators.Share(
			rx.MulticastReplayFactory(
				&rx.ReplayOptions{BufferSize: 1},
			),
		),
		operators.Take(4),
	)

	NewTestSuite(t).Case(
		rx.Merge(
			obs,
			obs.Pipe(DelaySubscription(4)),
			obs.Pipe(DelaySubscription(8)),
			obs.Pipe(DelaySubscription(16)),
		),
		0, 0, 1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 0, 1, 2, 3, Completed,
	).TestAll()
}
