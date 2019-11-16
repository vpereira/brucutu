package main

import (
	"net"
	"sync"
	"time"

	"github.com/bytbox/go-pop3"
	"github.com/emersion/go-imap/client"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

func dialHost(host string) (err error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return
	}
	conn.Close()
	return
}

func connectPOP3(wg *sync.WaitGroup, throttler <-chan int, host string, user string, password string) {
	defer wg.Done()
	c, err := pop3.Dial(host)
	if err != nil {
		<-throttler
		return
	}
	err = c.Auth(user, password)
	if err == nil {
		log.Info(user, ":", password, " was found")
	}
	defer c.Quit()

	<-throttler
}

func connectIMAP(wg *sync.WaitGroup, throttler <-chan int, host string, user string, password string) {
	defer wg.Done()
	c, err := client.Dial(host)
	if err != nil {
		<-throttler
		return
	}
	err = c.Login(user, password)
	if err == nil {
		log.Info(user, ":", password, " was found")
	}
	defer c.Logout()

	<-throttler
}
func connectSSH(wg *sync.WaitGroup, throttler <-chan int, host string, user string, password string) {
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
	log.Info(user, ":", password, " was found")
	defer c.Close()
	<-throttler
}
