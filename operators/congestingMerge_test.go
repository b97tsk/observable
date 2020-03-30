package operators_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_CongestingMergeAll(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			Just(
				Just("A", "B").Pipe(addLatencyToValue(3, 5)),
				Just("C", "D").Pipe(addLatencyToValue(2, 4)),
				Just("E", "F").Pipe(addLatencyToValue(1, 3)),
			).Pipe(operators.CongestingMergeAll()),
			Just(
				Just("A", "B").Pipe(addLatencyToValue(3, 5)),
				Just("C", "D").Pipe(addLatencyToValue(2, 4)),
				Just("E", "F").Pipe(addLatencyToValue(1, 3)),
			).Pipe(CongestingMergeConfigure{ProjectToObservable, 1}.Use()),
		},
		[][]interface{}{
			{"E", "C", "A", "F", "D", "B", Complete},
			{"A", "B", "C", "D", "E", "F", Complete},
		},
	)
}