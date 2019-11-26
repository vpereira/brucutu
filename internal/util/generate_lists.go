package util

/*
it tries to read a file assuming that -L <file> or -P <file> were given
if it returns an error, we try to read the -l <login> or -p <password> options
*/

// GenerateUserList method to read file and give back a list
func GenerateUserList(cli *CliArgument) (data []string, err error) {
	return generateList("login", cli)
}

// GeneratePasswordList method to read file and give back a list
func GeneratePasswordList(cli *CliArgument) (data []string, err error) {
	return generateList("password", cli)
}

func generateList(listType string, cli *CliArgument) (data []string, err error) {

	singleValue, listValue := getValues(listType, *cli)

	values, err := ReadFile(*listValue)

	if err != nil {
		if *singleValue != "" {
			return []string{*singleValue}, nil
		}
	}
	return values, err
}

func getValues(listType string, cli CliArgument) (*string, *string) {
	if listType == "login" {
		return cli.Login, cli.LoginList
	}
	return cli.Password, cli.PasswordList
}
