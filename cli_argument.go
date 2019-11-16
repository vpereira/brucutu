package main

import "flag"

type cliArgument struct {
	url                *string
	login              *string
	password           *string
	loginList          *string
	passwordList       *string
	outputFile         *string
	concurrency        *int
	supportedProtocols *bool
	tryLoginReverse    *bool
}

func (c *cliArgument) readParameters() {
	c.url = flag.String("u", "", "set url")
	c.login = flag.String("l", "", "set single login")
	c.password = flag.String("p", "", "set single password")
	c.loginList = flag.String("L", "", "set list of logins")
	c.passwordList = flag.String("P", "", "set list of passwords")
	c.outputFile = flag.String("o", "", "set output file")
	c.concurrency = flag.Int("c", 8, "number of concurrent goroutines")
	c.supportedProtocols = flag.Bool("m", false, "print the supported protocols")
	c.tryLoginReverse = flag.Bool("r", false, "try login reverse as password")
}
