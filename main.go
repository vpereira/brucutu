package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
)

// SupportedProtocols all our supported protocols
var SupportedProtocols = map[string]int{
	"ssh":  22,
	"pop3": 110,
	"imap": 143,
}

func printSupporteProtocols() {
	fmt.Println(SupportedProtocols)
}

// convention to write boolean methods?
func protocolSupported(protocol string) bool {
	_, ok := SupportedProtocols[protocol]
	return ok
}

func main() {
	cli := &cliArgument{}
	cli.readParameters()
	flag.Parse()
	if *cli.supportedProtocols == true {
		printSupporteProtocols()
		os.Exit(0)
	}

	// url is mandatory
	if *cli.url == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	myURL, err := parseURL(*cli.url)

	if err != nil {
		log.Fatal(*cli.url, "can't be parsed")
		os.Exit(1)
	}

	if !protocolSupported(myURL.Scheme) {
		log.Fatal("Protocol ", myURL.Scheme, " not supported")
		os.Exit(1)
	}

	// you are just allowed to choose one option for login and one option for password
	// -L and -l or -P and -p aren't allowed at the same time
	if (*cli.loginList != "" && *cli.login != "") || (*cli.password != "" && *cli.passwordList != "") {
		flag.PrintDefaults()
		os.Exit(1)
	}

	users, err := generateUserList(cli)

	if err != nil {
		log.Fatal("Can't read user list, exiting.")
		os.Exit(1)
	}

	passwords, err := generatePasswordList(cli)

	if err != nil {
		log.Fatal("Can't read password list, exiting.")
		os.Exit(1)
	}

	throttler := make(chan int, *cli.concurrency)
	outputChannel := make(chan string, len(users)*len(passwords))

	var host string
	if *cli.alternativePort != 0 {
		host = fmt.Sprintf("%s:%d", myURL.Host, *cli.alternativePort)
	} else {
		host = fmt.Sprintf("%s:%d", myURL.Host, SupportedProtocols[myURL.Scheme])
	}
	// test connection
	if err := dialHost(host); err != nil {
		log.Fatal("Couldn't connect to host", host, " exiting.")
		os.Exit(1)
	}

	go writeLog(outputChannel, *cli.quitFirstFound)
	var wg sync.WaitGroup

	for _, user := range users {
		for _, password := range passwords {
			throttler <- 0
			wg.Add(1)
			switch myURL.Scheme {
			case "pop3":
				go connectPOP3(&wg, throttler, outputChannel, host, user, password)
			case "ssh":
				go connectSSH(&wg, throttler, outputChannel, host, user, password)
			case "imap":
				go connectIMAP(&wg, throttler, outputChannel, host, user, password)
			default:
				log.Fatal("not implemented")
			}
		}
	}
	wg.Wait()
}

func writeLog(outputChannel chan string, quitFirstFound bool) {
	for {
		loginPassword, ok := <-outputChannel
		if ok {
			log.Info(loginPassword, " found")
			if quitFirstFound {
				os.Exit(0)
			}
		} else {
			break
		}
	}
}
