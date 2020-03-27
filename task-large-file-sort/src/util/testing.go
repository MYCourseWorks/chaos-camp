package util

import "testing"

func CheckExpectedError(got, want error, t *testing.T) {
	if want != nil {
		// expect error
		if want != got {
			t.Fatalf("Not the expected error -> want=%v got=%v", want, got)
		} else {
			// all good
			return
		}
	} else if got != nil {
		t.Fatalf("Unexpected error -> got=%v", got)
		return
	}
}
