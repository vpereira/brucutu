package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/bytbox/go-pop3"
	"github.com/emersion/go-imap/client"
	"golang.org/x/crypto/ssh"
)

func connectPOP3(wg *sync.WaitGroup, throttler <-chan int, output chan string, host string, user string, password string) {
	defer wg.Done()
	c, err := pop3.Dial(host)
	if err != nil {
		<-throttler
		return
	}
	err = c.Auth(user, password)
	if err == nil {
		output <- fmt.Sprintf("%s:%s", user, password)
	}
	defer c.Quit()

	<-throttler
}

func connectIMAP(wg *sync.WaitGroup, throttler <-chan int, output chan string, host string, user string, password string) {
	defer wg.Done()
	c, err := client.Dial(host)
	if err != nil {
		<-throttler
		return
	}
	err = c.Login(user, password)
	if err == nil {
		output <- fmt.Sprintf("%s:%s", user, password)
	}
	defer c.Logout()

	<-throttler
}
func connectSSH(wg *sync.WaitGroup, throttler <-chan int, output chan string, host string, user string, password string) {
	defer wg.Done()

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		// it should be configurable
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshConfig.SetDefaults()

	c, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		<-throttler
		return
	}
	output <- fmt.Sprintf("%s:%s", user, password)
	defer c.Close()
	<-throttler
}
