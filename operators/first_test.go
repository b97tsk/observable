package operators_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_First(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			Empty().Pipe(operators.First()),
			Throw(errTest).Pipe(operators.First()),
			Just("A").Pipe(operators.First()),
			Just("A", "B").Pipe(operators.First()),
			Concat(Just("A"), Throw(errTest)).Pipe(operators.First()),
			Concat(Just("A", "B"), Throw(errTest)).Pipe(operators.First()),
		},
		[][]interface{}{
			{ErrEmpty},
			{errTest},
			{"A", Complete},
			{"A", Complete},
			{"A", Complete},
			{"A", Complete},
		},
	)
}