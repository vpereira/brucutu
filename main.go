package main

import (
	"flag"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/vpereira/brucutu/internal/connect"
	"github.com/vpereira/brucutu/internal/util"
)

func main() {
	cli := util.NewCliArgument()

	err := cli.ReadParameters()

	if err != nil {
		log.Fatal(err.Error())
		flag.PrintDefaults()
		os.Exit(1)
	}

	myURL, err := cli.ParseURL()

	if err != nil {
		log.Fatal(*cli.URL, " can't be parsed")
		os.Exit(1)
	}

	if !util.ProtocolSupported(myURL.Scheme) {
		log.Fatal("Protocol ", myURL.Scheme, " not supported")
		os.Exit(1)
	}

	users, err := util.GenerateUserList(cli)

	if err != nil {
		log.Fatal("Can't read user list, exiting.")
		os.Exit(1)
	}

	passwords, err := util.GeneratePasswordList(cli)

	if err != nil {
		log.Fatal("Can't read password list, exiting:", err.Error())
		os.Exit(1)
	}

	throttler := make(chan int, *cli.Concurrency)
	outputChannel := make(chan string)

	host := util.SetHostName(cli, myURL)

	// test connection
	if err := util.DialHost(*host); err != nil {
		log.Fatal("util.DialHost", err.Error())
		os.Exit(1)
	}

	var outputFile *os.File

	if *cli.OutputFile != "" {
		outputFile, err := os.Create(*cli.OutputFile)

		if err != nil {
			log.Fatal("Output file", err.Error())
			os.Exit(1)
		}

		defer outputFile.Close()
	}

	// invert the logins, i.e foobar, become rabfoo and use it as password
	if *cli.TryLoginReverse {
		var reverseValues = make([]string, len(users))

		for _, user := range users {
			reverseValues = append(reverseValues, util.Reverse(user))
		}
		passwords = append(passwords, reverseValues...)
	}

	go util.WriteLog(outputChannel, outputFile, *cli.QuitFirstFound)
	var wg sync.WaitGroup

	for _, user := range users {
		for _, password := range passwords {
			throttler <- 0
			wg.Add(1)
			ca := connect.Arguments{StartTLS: *cli.StartTLS, UseTLS: *cli.UseTLS, Host: *host, User: user, Password: password}
			switch myURL.Scheme {
			case "pop3", "pop3s":
				go connect.POP3(&wg, throttler, outputChannel, ca)
			case "ssh":
				go connect.SSH(&wg, throttler, outputChannel, ca)
			case "imap", "imaps":
				go connect.IMAP(&wg, throttler, outputChannel, ca)
			default:
				log.Fatal("not implemented")
			}
		}
	}
	wg.Wait()
	close(outputChannel)
}
