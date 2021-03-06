package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestDelay(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			operators.Delay(Step(1)),
		),
		"A", "B", "C", "D", "E", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(0, 1),
			operators.Delay(Step(2)),
		),
		"A", "B", "C", "D", "E", Completed,
	).Case(
		rx.Concat(
			rx.Just("A", "B", "C", "D", "E"),
			rx.Throw(ErrTest),
		).Pipe(
			AddLatencyToNotifications(0, 2),
			operators.Delay(Step(1)),
		),
		"A", "B", "C", "D", "E", ErrTest,
	).Case(
		rx.Empty().Pipe(
			operators.Delay(Step(1)),
		),
		Completed,
	).TestAll()
}
