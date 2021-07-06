package mongo

import (
	"errors"
	"fmt"
	"strings"
)

type CredentialConfig struct {
	Username string
	Password string
}

type ConnectionOptions struct {
	SSL            bool
	SSLCert        string
	ReplicaSet     string
	ReadPreference string
	RetryWrites    bool
}

func (c *ConnectionOptions) generateParams() string {
	var params []string

	params = append(params, fmt.Sprintf("ssl=%t", c.SSL))

	if c.SSLCert != "" {
		params = append(params, fmt.Sprintf("ssl_ca_certs=%s", c.SSLCert))
	}

	if c.ReplicaSet != "" {
		params = append(params, fmt.Sprintf("replicaSet=%s", c.ReplicaSet))
	}

	if c.ReadPreference != "" {
		params = append(params, fmt.Sprintf("readPreference=%s", c.ReadPreference))
	}

	params = append(params, fmt.Sprintf("retryWrites=%t", c.RetryWrites))

	return strings.Join(params, "&")
}

type ClientConfig struct {
	Host         string
	Port         uint
	Database     string
	Clustered    bool
	DBNameInPath bool
	Credentials  *CredentialConfig
	Options      *ConnectionOptions
}

func (c *ClientConfig) generateURI() (string, error) {
	if c.Host == "" {
		return "", errors.New("host is empty")
	}

	if c.Port == 0 {
		return "", errors.New("port is not set")
	}

	if c.Database == "" {
		return "", errors.New("database is not set")
	}

	connectURI := ""

	if c.Clustered {
		connectURI = "mongodb+srv"
	} else {
		connectURI = "mongodb"
	}

	if c.Credentials != nil {
		if c.Credentials.Username == "" || c.Credentials.Password == "" {
			return "", errors.New("credentials are not set")
		}

		connectURI = fmt.Sprintf("%s://%s:%s@%s", connectURI, c.Credentials.Username, c.Credentials.Password, c.Host)
	} else {
		connectURI = fmt.Sprintf("%s://%s", connectURI, c.Host)
	}

	if !c.Clustered {
		connectURI = fmt.Sprintf("%s:%d", connectURI, c.Port)
	}

	if c.DBNameInPath {
		connectURI = fmt.Sprintf("%s/%s", connectURI, c.Database)
	} else {
		connectURI = fmt.Sprintf("%s/", connectURI)
	}

	if c.Options != nil {
		connectURI = fmt.Sprintf("%s?%s", connectURI, c.Options.generateParams())
	}

	return connectURI, nil
}
