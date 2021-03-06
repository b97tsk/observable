package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestMerge(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Merge(
			rx.Just("A", "B").Pipe(AddLatencyToValues(3, 5)),
			rx.Just("C", "D").Pipe(AddLatencyToValues(2, 4)),
			rx.Just("E", "F").Pipe(AddLatencyToValues(1, 3)),
		),
		"E", "C", "A", "F", "D", "B", Completed,
	).Case(
		rx.Merge(),
		Completed,
	).TestAll()
}
