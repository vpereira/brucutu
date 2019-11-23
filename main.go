package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/vpereira/brucutu/internal/util"
)

func main() {
	cli := &cliArgument{}
	cli.readParameters()
	flag.Parse()
	if *cli.supportedProtocols == true {
		util.PrintSupportedProtocols()
		os.Exit(0)
	}

	// url is mandatory
	if *cli.url == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	myURL, err := util.ParseURL(*cli.url)

	if err != nil {
		log.Fatal(*cli.url, " can't be parsed")
		os.Exit(1)
	}

	if !util.ProtocolSupported(myURL.Scheme) {
		log.Fatal("Protocol ", myURL.Scheme, " not supported")
		os.Exit(1)
	}

	// you are just allowed to choose one option for login and one option for password
	// -L and -l or -P and -p aren't allowed at the same time
	if (*cli.loginList != "" && *cli.login != "") || (*cli.password != "" && *cli.passwordList != "") {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// either use TLS or starttls. Both AFAIC arent to be used together
	if *cli.startTLS && *cli.useTLS {
		log.Fatal("starttls and use ssl are mutual exclusive")
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
	outputChannel := make(chan string)

	var host string
	if *cli.alternativePort != 0 {
		host = fmt.Sprintf("%s:%d", myURL.Host, *cli.alternativePort)
	} else {
		host = fmt.Sprintf("%s:%d", myURL.Host, util.SupportedProtocols[myURL.Scheme])
	}
	// test connection
	if err := util.DialHost(host); err != nil {
		log.Fatal("Couldn't connect to host", host, " exiting.")
		os.Exit(1)
	}

	go util.WriteLog(outputChannel, *cli.quitFirstFound)
	var wg sync.WaitGroup

	for _, user := range users {
		for _, password := range passwords {
			throttler <- 0
			wg.Add(1)
			ca := ConnectArguments{UseTLS: *cli.useTLS, Host: host, User: user, Password: password}
			switch myURL.Scheme {
			case "pop3", "pop3s":
				go connectPOP3(&wg, throttler, outputChannel, ca)
			case "ssh":
				go connectSSH(&wg, throttler, outputChannel, ca)
			case "imap", "imaps":
				go connectIMAP(&wg, throttler, outputChannel, ca)
			default:
				log.Fatal("not implemented")
			}
		}
	}
	wg.Wait()
	close(outputChannel)
}
