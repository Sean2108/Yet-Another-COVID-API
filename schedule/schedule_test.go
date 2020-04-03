package schedule

import (
	"testing"
	"time"
)

var functionCalled bool

func TestCallFunctionDaily_CallsFunction(t *testing.T) {
	functionCalled = false
	testFn := func() {
		functionCalled = true
	}
	CallFunctionDaily(testFn, 1)
	if !functionCalled {
		t.Error("Function was not called.")
	}
}

func TestGetTimeTillUpdate(t *testing.T) {
	tables := []struct {
		now      time.Time
		expected string
	}{
		{time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), "0us"},
		{time.Date(2020, 1, 1, 23, 0, 0, 0, time.UTC), "1h"},
		{time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC), "23h"},
		{time.Date(2020, 1, 1, 23, 59, 0, 0, time.UTC), "1m"},
	}

	for _, table := range tables {
		result := getTimeTillUpdate(0, table.now)
		expectedDuration, _ := time.ParseDuration(table.expected)
		if result != expectedDuration {
			t.Errorf("duration returned by getTimeTillUpdate was incorrect, got: %d, want %d.", result, expectedDuration)
		}
	}
}
