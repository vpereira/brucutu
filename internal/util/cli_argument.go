package util

import (
	"errors"
	"flag"
	"net/url"
	"os"
)

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

// NewCliArgument create new CliArgument
func NewCliArgument() *CliArgument {
	c := CliArgument{}
	c.setFlags()

	return &c
}

func (c *CliArgument) setFlags() {
	c.URL = flag.String("u", "", "set url")
	c.Login = flag.String("l", "", "set single login")
	c.LoginList = flag.String("L", "", "set list of logins")
	c.Password = flag.String("p", "", "set single password")
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

// ParseURL return the parsed url or error
func (c *CliArgument) ParseURL() (parsedURL *url.URL, err error) {
	myURL, err := url.Parse(*c.URL)
	if err != nil {
		return nil, err
	}
	return myURL, nil
}

// ReadParameters reads the flags and implement some validation logic
func (c *CliArgument) ReadParameters() (err error) {
	flag.Parse()

	// just print the list of supported protocols and exit as successful
	if *c.SupportedProtocols {
		PrintSupportedProtocols()
		os.Exit(0)
	}

	if *c.URL == "" {
		return errors.New("url taget cannot be empty")
	}
	if *c.Login != "" && *c.LoginList != "" {
		return errors.New("-L and -l are mutually exclusive")
	}

	if *c.Password != "" && *c.PasswordList != "" {
		return errors.New("-P and -p are mutually exclusive")
	}

	if *c.UseTLS && *c.StartTLS {
		return errors.New("-starttls and -tls are mutally exclusive")
	}
	return nil
}
