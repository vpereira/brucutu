package main

import (
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/vpereira/brucutu/internal/connect"
	"github.com/vpereira/brucutu/internal/util"
)

type target struct {
	scheme string
	host   string
}

func dispatch(scheme string, wg *sync.WaitGroup, throttler chan int, out chan string, ca connect.Arguments) {
	switch scheme {
	case "ftp":
		connect.FTP(wg, throttler, out, ca)
	case "http", "https":
		connect.HTTPBasicAuth(wg, throttler, out, ca)
	case "pop3", "pop3s":
		connect.POP3(wg, throttler, out, ca)
	case "ssh":
		connect.SSH(wg, throttler, out, ca)
	case "imap", "imaps":
		connect.IMAP(wg, throttler, out, ca)
	case "ldap", "ldaps":
		connect.LDAP(wg, throttler, out, ca)
	case "rdp":
		connect.RDP(wg, throttler, out, ca)
	default:
		log.Fatalf("not implemented: %s", scheme)
	}
}

func main() {
	cli := util.NewCliArgument()

	if err := cli.ReadParameters(); err != nil {
		log.Fatal(err.Error())
	}

	myURL, err := cli.ParseURL()
	if err != nil {
		log.Fatal(*cli.URL, " can't be parsed")
	}

	if !util.ProtocolSupported(myURL.Scheme) {
		log.Fatal("Protocol ", myURL.Scheme, " not supported")
	}

	users, err := util.GenerateUserList(cli)
	if err != nil {
		log.Fatal("Can't read user list, exiting.")
	}

	passwords, err := util.GeneratePasswordList(cli)
	if err != nil {
		log.Fatal("Can't read password list, exiting:", err.Error())
	}

	host := util.SetHostName(cli, myURL)

	// test connection to primary target
	if err := util.DialHost(*host); err != nil {
		log.Fatal("util.DialHost", err.Error())
	}

	// LDAP username enumeration pre-step (Active Directory only, no -D needed)
	if *cli.LDAPEnum {
		if myURL.Scheme == "ldap" || myURL.Scheme == "ldaps" {
			log.Info("Running LDAP username enumeration...")
			users = connect.LDAPEnumerate(*host, users)
			log.Infof("%d valid users after enumeration", len(users))
		} else {
			log.Warn("-e flag ignored: only valid for ldap/ldaps")
		}
	}

	var outputFile *os.File
	if *cli.OutputFile != "" {
		outputFile, err = os.Create(*cli.OutputFile)
		if err != nil {
			log.Fatal("Output file", err.Error())
		}
		defer outputFile.Close()
	}

	if *cli.TryLoginReverse {
		reverseValues := make([]string, 0, len(users))
		for _, user := range users {
			reverseValues = append(reverseValues, util.Reverse(user))
		}
		passwords = append(passwords, reverseValues...)
	}

	// Build target list. -H overrides the primary protocol for round-robin hopping.
	// Each protocol uses its default port; -a is only applied to the primary target.
	var targets []target
	if *cli.HopProtocols != "" {
		for _, proto := range strings.Split(*cli.HopProtocols, ",") {
			proto = strings.TrimSpace(proto)
			targets = append(targets, target{
				scheme: proto,
				host:   util.HostForProtocol(myURL.Host, proto),
			})
		}
	} else {
		targets = []target{{scheme: myURL.Scheme, host: *host}}
	}

	throttler := make(chan int, *cli.Concurrency)
	outputChannel := make(chan string)

	go util.WriteLog(outputChannel, outputFile, *cli.QuitFirstFound)

	var wg sync.WaitGroup
	hopIdx := 0
	for _, user := range users {
		for _, password := range passwords {
			t := targets[hopIdx%len(targets)]
			hopIdx++
			throttler <- 0
			wg.Add(1)
			ca := connect.Arguments{
				StartTLS: *cli.StartTLS,
				UseTLS:   *cli.UseTLS || t.scheme == "ldaps",
				Host:     t.host,
				User:     user,
				Password: password,
			}
			go dispatch(t.scheme, &wg, throttler, outputChannel, ca)
		}
	}

	wg.Wait()
	close(outputChannel)
}
