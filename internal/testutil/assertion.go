package testutil

import "testing"

func AssertEqual[T comparable](t *testing.T, received T, expected T) {
	if received != expected {
		t.Errorf("Expected %v but got %v", expected, received)
	}
}
