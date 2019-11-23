package util

import (
	"flag"
	"testing"
)

func TestGenerateLists(t *testing.T) {

	var url string
	var userFile string
	var passwdFile string
	cli := &CliArgument{}
	flag.StringVar(&url, "u", "ssh://127.0.0.1", "")
	flag.StringVar(&userFile, "L", "../../samples/users.txt", "")
	flag.StringVar(&passwdFile, "P", "../../samples/passwd.txt", "")

	cli.URL = &url
	cli.LoginList = &userFile
	cli.PasswordList = &passwdFile

	ul, error := GenerateUserList(cli)

	if error != nil {
		t.Errorf("user list generation failed")
	}
	if len(ul) != 3 {
		t.Errorf("user list size is wrong")
	}
	pl, error := GeneratePasswordList(cli)

	if error != nil {
		t.Errorf("user list generation failed")
	}
	if len(pl) != 5 {
		t.Errorf("user list size is wrong")
	}
}
