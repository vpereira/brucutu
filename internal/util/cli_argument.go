package util

import "flag"

// CliArgument handles the command arguments
type CliArgument struct {
	URL                *string
	Login              *string
	Password           *string
	LoginList          *string
	PasswordList       *string
	OutputFile         *string
	Concurrency        *int
	SupportedProtocols *bool
	TryLoginReverse    *bool
	AlternativePort    *int
	QuitFirstFound     *bool
	UseTLS             *bool
	StartTLS           *bool
}

// ReadParameters reads the flags
func (c *CliArgument) ReadParameters() {
	c.URL = flag.String("u", "", "set url")
	c.Login = flag.String("l", "", "set single login")
	c.Password = flag.String("p", "", "set single password")
	c.LoginList = flag.String("L", "", "set list of logins")
	c.PasswordList = flag.String("P", "", "set list of passwords")
	c.OutputFile = flag.String("o", "", "set output file")
	c.Concurrency = flag.Int("c", 8, "number of concurrent goroutines")
	c.SupportedProtocols = flag.Bool("m", false, "print the supported protocols")
	c.TryLoginReverse = flag.Bool("r", false, "try login reverse as password")
	c.AlternativePort = flag.Int("a", 0, "set alternative port for service")
	c.QuitFirstFound = flag.Bool("f", false, "Quit as soon first password was found")
	c.UseTLS = flag.Bool("tls", false, "Use SSL/TLS")
	c.StartTLS = flag.Bool("starttls", false, "Use starttls")
}
