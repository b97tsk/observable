package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestCongestingMergeAll(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just(
				rx.Just("A", "B").Pipe(AddLatencyToValues(3, 5)),
				rx.Just("C", "D").Pipe(AddLatencyToValues(2, 4)),
				rx.Just("E", "F").Pipe(AddLatencyToValues(1, 3)),
			).Pipe(operators.CongestingMergeAll()),
			rx.Just(
				rx.Just("A", "B").Pipe(AddLatencyToValues(3, 5)),
				rx.Just("C", "D").Pipe(AddLatencyToValues(2, 4)),
				rx.Just("E", "F").Pipe(AddLatencyToValues(1, 3)),
			).Pipe(operators.CongestingMergeConfigure{Concurrency: 1}.Make()),
		},
		[][]interface{}{
			{"E", "C", "A", "F", "D", "B", Completed},
			{"A", "B", "C", "D", "E", "F", Completed},
		},
	)
}