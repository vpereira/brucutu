package connect

import (
	"crypto/tls"
	"fmt"
	"sync"

	"github.com/go-ldap/ldap/v3"
)

// LDAP bruteforces via simple bind. The bind DN is taken verbatim from ca.User,
// supporting UPN (user@domain), downlevel (DOMAIN\user), or full DN (cn=user,dc=...).
func LDAP(wg *sync.WaitGroup, throttler <-chan int, output chan string, ca Arguments) {
	defer wg.Done()

	var conn *ldap.Conn
	var err error

	if ca.UseTLS {
		conn, err = ldap.DialURL(fmt.Sprintf("ldaps://%s", ca.Host), ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	} else {
		conn, err = ldap.DialURL(fmt.Sprintf("ldap://%s", ca.Host))
	}

	if err != nil {
		<-throttler
		return
	}
	defer conn.Close()

	if ca.StartTLS {
		if err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
			<-throttler
			return
		}
	}

	if err = conn.Bind(ca.User, ca.Password); err == nil {
		output <- fmt.Sprintf("%s:%s", ca.User, ca.Password)
	}

	<-throttler
}
