package util

import "testing"

func TestListenSignal(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{name: "current"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ListenSignal()
		})
	}
}
