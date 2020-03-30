package operators_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Audit(t *testing.T) {
	subscribe(
		t,
		Just("A", "B", "C", "D", "E", "F", "G").Pipe(
			addLatencyToValue(1, 2),
			operators.Audit(func(interface{}) Observable {
				return Interval(step(3))
			}),
		),
		"B", "D", "F", Complete,
	)
}