package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestSkipWhile(t *testing.T) {
	lessThanFive := func(val interface{}, idx int) bool {
		return val.(int) < 5
	}

	NewTestSuite(t).Case(
		rx.Just(1, 2, 3, 4, 5, 4, 3, 2, 1).Pipe(
			operators.SkipWhile(lessThanFive),
		),
		5, 4, 3, 2, 1, Completed,
	).Case(
		rx.Concat(
			rx.Range(1, 9),
			rx.Throw(ErrTest),
		).Pipe(
			operators.SkipWhile(lessThanFive),
		),
		5, 6, 7, 8, ErrTest,
	).Case(
		rx.Concat(
			rx.Range(1, 5),
			rx.Throw(ErrTest),
		).Pipe(
			operators.SkipWhile(lessThanFive),
		),
		ErrTest,
	).TestAll()
}
