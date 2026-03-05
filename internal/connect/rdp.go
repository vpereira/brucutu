package connect

import (
	"fmt"
	"sync"

	"github.com/icodeface/grdp"
	"github.com/icodeface/grdp/glog"
)

// RDP bruteforces via NTLMv2 authentication.
// User format: bare username or "DOMAIN\user" — both handled by grdp internally.
func RDP(wg *sync.WaitGroup, throttler <-chan int, output chan string, ca Arguments) {
	defer wg.Done()

	c := grdp.NewClient(ca.Host, glog.NONE)

	if err := c.Login(ca.User, ca.Password); err == nil {
		output <- fmt.Sprintf("%s:%s", ca.User, ca.Password)
	}

	<-throttler
}
