package main

import "github.com/vpereira/brucutu/internal/util"
/*
it tries to read a file assuming that -L <file> or -P <file> were given
if it returns an error, we try to read the -l <login> or -p <password> options
*/
func generateUserList(cli *cliArgument) (data []string, err error) {
	users, err := util.ReadFile(*cli.loginList)

	if err != nil {
		if *cli.login != "" {
			return []string{*cli.login}, nil
		}
	}
	return users, err
}

func generatePasswordList(cli *cliArgument) (data []string, err error) {
	passwords, err := util.ReadFile(*cli.passwordList)

	if err != nil {
		if *cli.password != "" {
			return []string{*cli.password}, nil
		}
	}
	return passwords, err
}
