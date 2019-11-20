package util

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

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

// SupportedProtocols all our supported protocols
var SupportedProtocols = map[string]int{
	"ssh":  22,
	"pop3": 110,
	"imap": 143,
}

//PrintSupportedProtocols can be improved
func PrintSupportedProtocols() {
	fmt.Println(SupportedProtocols)
}

// ProtocolSupported there is convention to write boolean methods?
func ProtocolSupported(protocol string) bool {
	_, ok := SupportedProtocols[protocol]
	return ok
}
