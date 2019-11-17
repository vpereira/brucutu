package main

import (
	"flag"
	"testing"
)

func TestGenerateLists(t *testing.T) {

	var url string
	var userFile string
	var passwdFile string
	cli := &cliArgument{}
	flag.StringVar(&url, "u", "ssh://127.0.0.1", "")
	flag.StringVar(&userFile, "L", "sample/users.txt", "")
	flag.StringVar(&passwdFile, "P", "sample/passwd.txt", "")

	cli.url = &url
	cli.loginList = &userFile
	cli.passwordList = &passwdFile

	ul, error := generateUserList(cli)

	if error != nil {
		t.Errorf("user list generation failed")
	}
	if len(ul) != 3 {
		t.Errorf("user list size is wrong")
	}
	pl, error := generatePasswordList(cli)

	if error != nil {
		t.Errorf("user list generation failed")
	}
	if len(pl) != 5 {
		t.Errorf("user list size is wrong")
	}
}
