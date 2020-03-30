package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestConcat(t *testing.T) {
	subscribe(
		t,
		Concat(
			Just("A", "B").Pipe(addLatencyToValue(3, 5)),
			Just("C", "D").Pipe(addLatencyToValue(2, 4)),
			Just("E", "F").Pipe(addLatencyToValue(1, 3)),
		),
		"A", "B", "C", "D", "E", "F", Complete,
	)
}
