package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestWithLatestFrom(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("A", "B").Pipe(
			AddLatencyToValues(1, 2),
			operators.WithLatestFrom(
				rx.Range(1, 4).Pipe(
					AddLatencyToNotifications(0, 2),
				),
			),
			ToString(),
		),
		"[A 1]", "[B 2]", Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(1, 2),
			operators.WithLatestFrom(
				rx.Range(1, 4).Pipe(
					AddLatencyToNotifications(0, 2),
				),
			),
			ToString(),
		),
		"[A 1]", "[B 2]", "[C 3]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D").Pipe(
			AddLatencyToValues(1, 2),
			operators.WithLatestFrom(
				rx.Range(1, 4).Pipe(
					AddLatencyToNotifications(0, 2),
				),
			),
			ToString(),
		),
		"[A 1]", "[B 2]", "[C 3]", "[D 3]", Completed,
	).TestAll()

	NewTestSuite(t).Case(
		rx.Just("A", "B").Pipe(
			AddLatencyToValues(1, 2),
			operators.WithLatestFrom(
				rx.Concat(
					rx.Range(1, 4),
					rx.Throw(ErrTest),
				).Pipe(
					AddLatencyToNotifications(0, 2),
				),
			),
			ToString(),
		),
		"[A 1]", "[B 2]", Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			AddLatencyToValues(1, 2),
			operators.WithLatestFrom(
				rx.Concat(
					rx.Range(1, 4),
					rx.Throw(ErrTest),
				).Pipe(
					AddLatencyToNotifications(0, 2),
				),
			),
			ToString(),
		),
		"[A 1]", "[B 2]", "[C 3]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D").Pipe(
			AddLatencyToValues(1, 2),
			operators.WithLatestFrom(
				rx.Concat(
					rx.Range(1, 4),
					rx.Throw(ErrTest),
				).Pipe(
					AddLatencyToNotifications(0, 2),
				),
			),
			ToString(),
		),
		"[A 1]", "[B 2]", "[C 3]", ErrTest,
	).TestAll()
}
