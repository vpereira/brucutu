package util

/*
it tries to read a file assuming that -L <file> or -P <file> were given
if it returns an error, we try to read the -l <login> or -p <password> options
*/

// GenerateUserList method to read file and give back a list
func GenerateUserList(cli *CliArgument) (data []string, err error) {
	users, err := ReadFile(*cli.LoginList)

	if err != nil {
		if *cli.Login != "" {
			return []string{*cli.Login}, nil
		}
	}
	return users, err
}

// GeneratePasswordList method to read file and give back a list
func GeneratePasswordList(cli *CliArgument) (data []string, err error) {
	passwords, err := ReadFile(*cli.PasswordList)

	if err != nil {
		if *cli.Password != "" {
			return []string{*cli.Password}, nil
		}
	}
	return passwords, err
}
