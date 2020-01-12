package util

import (
	"testing"
)

func TestSetHostnameNoPortSet(t *testing.T) {
	sshURL := "ssh://127.0.0.1"
	altPort := 0
	cli := &CliArgument{}

	cli.URL = &sshURL
	cli.AlternativePort = &altPort

	myURL, _ := cli.ParseURL()

	myHostName := SetHostName(cli, myURL)

	if *myHostName == "" {
		t.Errorf("SetHostName failed")
	}
}
