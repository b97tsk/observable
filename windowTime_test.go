package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_WindowTime(t *testing.T) {
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
				operators.WindowTime(step(2)),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.WindowTime(step(4)),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.WindowTime(step(6)),
				operators.MergeMap(toSlice),
				toString,
			),
		},
		"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", "[G]", xComplete,
		"[A B]", "[C D]", "[E F]", "[G]", xComplete,
		"[A B C]", "[D E F]", "[G]", xComplete,
	)
	t.Log("----------")
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				WindowTimeConfigure{step(8), 0, 0}.MakeFunc(),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				WindowTimeConfigure{step(8), 0, 3}.MakeFunc(),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				WindowTimeConfigure{step(8), 0, 2}.MakeFunc(),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				WindowTimeConfigure{step(8), 0, 1}.MakeFunc(),
				operators.MergeMap(toSlice),
				toString,
			),
		},
		"[A B C D]", "[E F G]", xComplete,
		"[A B C]", "[D E F]", "[G]", xComplete,
		"[A B]", "[C D]", "[E F]", "[G]", xComplete,
		"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", "[G]", "[]", xComplete,
	)
	t.Log("----------")
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				WindowTimeConfigure{step(2), step(2), 0}.MakeFunc(),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				WindowTimeConfigure{step(2), step(4), 0}.MakeFunc(),
				operators.MergeMap(toSlice),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				WindowTimeConfigure{step(4), step(2), 0}.MakeFunc(),
				operators.MergeMap(toSlice),
				toString,
			),
		},
		"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", "[G]", xComplete,
		"[A]", "[C]", "[E]", "[G]", xComplete,
		"[A B]", "[B C]", "[C D]", "[D E]", "[E F]", "[F G]", "[G]", xComplete,
	)
}