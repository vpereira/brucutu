package util

import "testing"

func TestParseURL(t *testing.T) {
	_, err := ParseURL(":foo")
	if err == nil {
		t.Errorf("URL shouldn't be parsable")
	}
}
