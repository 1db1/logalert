package main

import "testing"

func TestNewBuffer(t *testing.T) {

	tests := []struct {
		name     string
		bufSize  string
		expected int
	}{
		{
			"1. Less then 1 Kb",
			"500b",
			0,
		},
		{
			"2. More then limit",
			"1000Mb",
			0,
		},
		{
			"3. Incorrect input 1",
			"10kbb",
			0,
		},
		{
			"4. Incorrect input 3",
			"a kb",
			0,
		},
		{
			"5. Incorrect input 4",
			"",
			0,
		},
		{
			"6. Incorrect input 5",
			"kb",
			0,
		},
		{
			"7. New buffer 1Kb",
			"1Kb",
			1024,
		},
		{
			"8. New buffer 1 KB",
			"1 KB",
			1024,
		},
		{
			"9. New buffer 8Mb",
			"8MB",
			8 * 1024 * 1024,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := newBuffer(tt.bufSize)
			if len(buf) != tt.expected {
				t.Errorf("Expected buffer with length %d, received %d", tt.expected, len(buf))
			}
		})
	}
}
