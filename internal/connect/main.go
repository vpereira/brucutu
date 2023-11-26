package connect

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/emersion/go-imap/client"
	"github.com/simia-tech/go-pop3"
	"golang.org/x/crypto/ssh"
)

// Arguments used to bring arguments from CLI
type Arguments struct {
	UseTLS   bool
	StartTLS bool
	Host     string
	User     string
	Password string
}

// HTTP Basic Auth Bruteforce
func HTTPBasicAuth(wg *sync.WaitGroup, throttler <-chan int, output chan string, ca Arguments) {
	defer wg.Done()

	var httpURL string
	if ca.UseTLS {
		httpURL = "https://" + ca.Host
	} else {
		httpURL = "http://" + ca.Host
	}

	req, err := http.NewRequest("GET", httpURL, nil)
	if err != nil {
		<-throttler
		return
	}

	auth := base64.StdEncoding.EncodeToString([]byte(ca.User + ":" + ca.Password))
	req.Header.Add("Authorization", "Basic "+auth)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		<-throttler
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		output <- fmt.Sprintf("%s:%s", ca.User, ca.Password)
	}

	<-throttler
}

// POP3 Bruteforce
func POP3(wg *sync.WaitGroup, throttler <-chan int, output chan string, ca Arguments) {
	defer wg.Done()
	var c *pop3.Client
	var err error

	if ca.UseTLS {
		c, err = pop3.Dial(ca.Host, pop3.UseTLS(&tls.Config{InsecureSkipVerify: true}))
	} else {
		c, err = pop3.Dial(ca.Host)
	}

	defer c.Quit()

	if err != nil {
		<-throttler
		return
	}
	err = c.Auth(ca.User, ca.Password)
	if err == nil {
		output <- fmt.Sprintf("%s:%s", ca.User, ca.Password)
	}

	<-throttler
}

// IMAP Bruteforce
func IMAP(wg *sync.WaitGroup, throttler <-chan int, output chan string, ca Arguments) {
	defer wg.Done()
	var c *client.Client
	var err error

	if ca.UseTLS {
		c, err = client.DialTLS(ca.Host, &tls.Config{InsecureSkipVerify: true})
	} else {
		c, err = client.Dial(ca.Host)
	}

	defer c.Logout()

	if err != nil {
		<-throttler
		return
	}
	err = c.Login(ca.User, ca.Password)

	if err == nil {
		output <- fmt.Sprintf("%s:%s", ca.User, ca.Password)
	}

	<-throttler
}

// SSH bruteforce
func SSH(wg *sync.WaitGroup, throttler <-chan int, output chan string, ca Arguments) {
	defer wg.Done()

	sshConfig := &ssh.ClientConfig{
		User: ca.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(ca.Password),
		},
		// it should be configurable
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshConfig.SetDefaults()

	c, err := ssh.Dial("tcp", ca.Host, sshConfig)
	if err != nil {
		<-throttler
		return
	}
	defer c.Close()
	output <- fmt.Sprintf("%s:%s", ca.User, ca.Password)
	<-throttler
}
