package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_WindowCount(t *testing.T) {
	toSlice := func(val interface{}, idx int) Observable {
		if obs, ok := val.(Observable); ok {
			return obs.Pipe(
				operators.ToSlice(),
			)
		}
		return Throw(ErrNotObservable)
	}
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.WindowCount(2),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.WindowCount(3),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				WindowCountConfigure{3, 1}.MakeFunc(),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				WindowCountConfigure{3, 2}.MakeFunc(),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				WindowCountConfigure{3, 4}.MakeFunc(),
				operators.MergeMap(toSlice),
				toString,
			),
		},
		"[A B]", "[C D]", "[E F]", "[G]", xComplete,
		"[A B C]", "[D E F]", "[G]", xComplete,
		"[A B C]", "[B C D]", "[C D E]", "[D E F]", "[E F G]", "[F G]", "[G]", "[]", xComplete,
		"[A B C]", "[C D E]", "[E F G]", "[G]", xComplete,
		"[A B C]", "[E F G]", "[]", xComplete,
	)
}