package operators_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Scan(t *testing.T) {
	max := func(acc, val interface{}, idx int) interface{} {
		if acc.(int) > val.(int) {
			return acc
		}
		return val
	}

	sum := func(acc, val interface{}, idx int) interface{} {
		return acc.(int) + val.(int)
	}

	subscribeN(
		t,
		[]Observable{
			Range(1, 7).Pipe(operators.Scan(max)),
			Just(42).Pipe(operators.Scan(max)),
			Empty().Pipe(operators.Scan(max)),
			Range(1, 7).Pipe(operators.Scan(sum)),
			Just(42).Pipe(operators.Scan(sum)),
			Empty().Pipe(operators.Scan(sum)),
			Throw(errTest).Pipe(operators.Scan(sum)),
		},
		[][]interface{}{
			{1, 2, 3, 4, 5, 6, Complete},
			{42, Complete},
			{Complete},
			{1, 3, 6, 10, 15, 21, Complete},
			{42, Complete},
			{Complete},
			{errTest},
		},
	)
}