package connect

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
	log "github.com/sirupsen/logrus"
)

// LDAPEnumerate uses the LDAP Ping (NetLogon cLDAP) technique from ldapnomnom
// to validate which usernames exist in Active Directory before brute-forcing.
// Returns the filtered list; on connection failure returns the original list unchanged.
func LDAPEnumerate(host string, users []string) []string {
	conn, err := ldap.DialURL(fmt.Sprintf("ldap://%s", host))
	if err != nil {
		log.Warnf("LDAP enumeration: cannot connect to %s: %v — returning full list", host, err)
		return users
	}
	defer conn.Close()

	var valid []string
	for _, username := range users {
		req := ldap.NewSearchRequest(
			"",
			ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false,
			fmt.Sprintf("(&(NtVer=\x06\x00\x00\x00)(AAC=\x10\x00\x00\x00)(User=%s))", username),
			[]string{"NetLogon"},
			nil,
		)
		resp, err := conn.Search(req)
		if err != nil {
			continue
		}
		if len(resp.Entries) > 0 &&
			len(resp.Entries[0].Attributes) > 0 &&
			len(resp.Entries[0].Attributes[0].ByteValues) > 0 {
			res := resp.Entries[0].Attributes[0].ByteValues[0]
			if len(res) > 2 && res[0] == 0x17 && res[1] == 0x00 {
				valid = append(valid, username)
			}
		}
	}
	return valid
}
