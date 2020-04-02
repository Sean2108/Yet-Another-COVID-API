package utils

import "testing"

func TestEditDistance(t *testing.T) {
	tables := []struct {
		str1     string
		str2     string
		expected int
	}{
		{"test", "tst", 1},
		{"test", "tast", 1},
		{"testt", "test", 1},
		{"test", "test", 0},
		{"testt", "tst", 2},
		{"abcd", "bcde", 2},
		{"", "", 0},
	}

	for _, table := range tables {
		result := EditDistance([]rune(table.str1), []rune(table.str2))
		if result != table.expected {
			t.Errorf("Result of EditDistance was incorrect, got: %d, want: %d.", result, table.expected)
		}
	}
}
