package util

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
)

// ReadFile trransform string in io
func ReadFile(f string) (data []string, err error) {
	b, err := os.Open(f)
	if err != nil {
		return
	}
	defer b.Close()
	scanner := bufio.NewScanner(b)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	return
}

// DialHost test if connection to host:port is possible
func DialHost(host string) (err error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return
	}
	conn.Close()
	return
}

var supportedProtocols = map[string]int{
	"ssh":   22,
	"pop3":  110,
	"imap":  143,
	"pop3s": 995,
	"imaps": 993,
}

//PrintSupportedProtocols can be improved
func PrintSupportedProtocols() {
	fmt.Println(supportedProtocols)
}

// ProtocolSupported there is convention to write boolean methods?
func ProtocolSupported(protocol string) bool {
	_, ok := supportedProtocols[protocol]
	return ok
}

// SetHostName return hostname in the format "host:port"
func SetHostName(cli *CliArgument, myURL *url.URL) *string {
	var host string
	if *cli.AlternativePort != 0 {
		host = fmt.Sprintf("%s:%d", myURL.Host, *cli.AlternativePort)
	} else {
		host = fmt.Sprintf("%s:%d", myURL.Host, supportedProtocols[myURL.Scheme])
	}
	return &host
}

// WriteLog goroutine used to log the password found
func WriteLog(outputChannel chan string, outputFile *os.File, quitFirstFound bool) {
	for {
		loginPassword, ok := <-outputChannel
		if ok {
			log.Info(loginPassword, " found")
			if outputFile != nil {
				outputFile.WriteString(fmt.Sprintf("%s\n", loginPassword))
			}
			if quitFirstFound {
				os.Exit(0)
			}
		} else {
			break
		}
	}
}
