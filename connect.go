package main

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/emersion/go-imap/client"
	"github.com/simia-tech/go-pop3"
	"golang.org/x/crypto/ssh"
)

// ConnectArguments used to bring arguments from CLI
type ConnectArguments struct {
	UseTLS   bool
	StartTLS bool
	Host     string
	User     string
	Password string
}

func connectPOP3(wg *sync.WaitGroup, throttler <-chan int, output chan string, ca ConnectArguments) {
	defer wg.Done()
	var c *pop3.Client
	var err error

	if ca.UseTLS {
		c, err = pop3.Dial(ca.Host, pop3.UseTLS(&tls.Config{InsecureSkipVerify: true}))
	} else {
		c, err = pop3.Dial(ca.Host)
	}
	if err != nil {
		<-throttler
		return
	}
	err = c.Auth(ca.User, ca.Password)
	if err == nil {
		output <- fmt.Sprintf("%s:%s", ca.User, ca.Password)
	}
	defer c.Quit()

	<-throttler
}

func connectIMAP(wg *sync.WaitGroup, throttler <-chan int, output chan string, ca ConnectArguments) {
	defer wg.Done()
	var c *client.Client
	var err error

	if ca.UseTLS {
		c, err = client.DialTLS(ca.Host, &tls.Config{InsecureSkipVerify: true})
	} else {
		c, err = client.Dial(ca.Host)
	}
	if err != nil {
		<-throttler
		return
	}
	err = c.Login(ca.User, ca.Password)
	if err == nil {
		output <- fmt.Sprintf("%s:%s", ca.User, ca.Password)
	}
	defer c.Logout()

	<-throttler
}
func connectSSH(wg *sync.WaitGroup, throttler <-chan int, output chan string, ca ConnectArguments) {
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
	output <- fmt.Sprintf("%s:%s", ca.User, ca.Password)
	defer c.Close()
	<-throttler
}
