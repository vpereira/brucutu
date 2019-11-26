package util

import "testing"

func TestParseURL(t *testing.T) {
	url := ":foo"
	cli := &CliArgument{URL: &url}
	_, err := cli.ParseURL()
	if err == nil {
		t.Errorf("URL shouldn't be parsable")
	}
}
