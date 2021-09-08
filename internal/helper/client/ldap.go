package client

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

func DialAndBind(c *Config) (*ldap.Conn, error) {
	conn, err := dial(c)
	if err != nil {
		return nil, err
	}

	if c.UseGSSAPI {
		ccache := strings.TrimLeft(c.CCache, "FILE:")
		spn := fmt.Sprintf("ldap/%s", c.LDAPHost)
		err = conn.GSSAPICCBind("/etc/krb5.conf", ccache, spn)
	} else {
		err = conn.Bind(c.BindUser, c.BindPassword)
	}
	if err != nil {
		conn.Close()
		return nil, err
	}

	// return the LDAP connection
	return conn, nil
}

func dial(c *Config) (*ldap.Conn, error) {
	uri := fmt.Sprintf("%s:%d", c.LDAPHost, c.LDAPPort)

	if c.TLS {
		return ldap.DialTLS("tcp", uri, &tls.Config{
			ServerName:         c.LDAPHost,
			InsecureSkipVerify: c.TLSInsecure,
		})
	}

	conn, err := ldap.Dial("tcp", uri)
	if err != nil {
		return nil, err
	}

	if c.StartTLS {
		err = conn.StartTLS(&tls.Config{
			InsecureSkipVerify: true,
		})
		if err != nil {
			return nil, err
		}
	}
	return conn, err
}
