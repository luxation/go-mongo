package mongo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParamsGeneration(t *testing.T) {
	tests := []struct {
		options ConnectionOptions
		output  string
	}{
		{
			options: ConnectionOptions{},
			output:  "ssl=false&retryWrites=false",
		},
		{
			options: ConnectionOptions{
				SSL: true,
			},
			output: "ssl=true&retryWrites=false",
		},
		{
			options: ConnectionOptions{
				SSL:     true,
				SSLCert: "cert-1",
			},
			output: "ssl=true&ssl_ca_certs=cert-1&retryWrites=false",
		},
		{
			options: ConnectionOptions{
				SSL:        true,
				SSLCert:    "cert-1",
				ReplicaSet: "replica-1",
			},
			output: "ssl=true&ssl_ca_certs=cert-1&replicaSet=replica-1&retryWrites=false",
		},
		{
			options: ConnectionOptions{
				SSL:            true,
				SSLCert:        "cert-1",
				ReplicaSet:     "replica-1",
				ReadPreference: "preference",
			},
			output: "ssl=true&ssl_ca_certs=cert-1&replicaSet=replica-1&readPreference=preference&retryWrites=false",
		},
		{
			options: ConnectionOptions{
				SSL:            true,
				SSLCert:        "cert-1",
				ReplicaSet:     "replica-1",
				ReadPreference: "preference",
				RetryWrites:    true,
			},
			output: "ssl=true&ssl_ca_certs=cert-1&replicaSet=replica-1&readPreference=preference&retryWrites=true",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.output, test.options.generateParams())
	}

}

func TestUriGeneration(t *testing.T) {
	tests := []struct {
		config     ClientConfig
		output     string
		hasError   bool
		stacktrace string
	}{
		{
			config:     ClientConfig{},
			hasError:   true,
			output:     "",
			stacktrace: "host is empty",
		},
		{
			config: ClientConfig{
				Host: "host",
			},
			hasError:   true,
			output:     "",
			stacktrace: "port is not set",
		},
		{
			config: ClientConfig{
				Host: "host",
				Port: 1,
			},
			hasError:   true,
			output:     "",
			stacktrace: "database is not set",
		},
		{
			config: ClientConfig{
				Host:     "host",
				Port:     1,
				Database: "db",
			},
			hasError:   false,
			output:     "mongodb://host:1",
			stacktrace: "",
		},
		{
			config: ClientConfig{
				Host:     "host",
				Port:     1,
				Database: "db",
				Credentials: &CredentialConfig{
					Username: "",
					Password: "",
				},
			},
			hasError:   true,
			output:     "",
			stacktrace: "credentials are not set",
		},
		{
			config: ClientConfig{
				Host:     "host",
				Port:     1,
				Database: "db",
				Credentials: &CredentialConfig{
					Username: "user",
					Password: "",
				},
			},
			hasError:   true,
			output:     "",
			stacktrace: "credentials are not set",
		},
		{
			config: ClientConfig{
				Host:     "host",
				Port:     1,
				Database: "db",
				Credentials: &CredentialConfig{
					Username: "user",
					Password: "password",
				},
			},
			hasError:   false,
			output:     "mongodb://user:password@host:1",
			stacktrace: "",
		},
		{
			config: ClientConfig{
				Host:     "host",
				Port:     1,
				Database: "db",
				Credentials: &CredentialConfig{
					Username: "user",
					Password: "password",
				},
				Options: &ConnectionOptions{},
			},
			hasError:   false,
			output:     "mongodb://user:password@host:1/?ssl=false&retryWrites=false",
			stacktrace: "",
		},
		{
			config: ClientConfig{
				Host:     "host",
				Port:     1,
				Database: "db",
				Credentials: &CredentialConfig{
					Username: "user",
					Password: "password",
				},
				Options: &ConnectionOptions{
					SSL:            true,
					SSLCert:        "cert-1",
					ReplicaSet:     "replica-1",
					ReadPreference: "preference",
					RetryWrites:    true,
				},
			},
			hasError:   false,
			output:     "mongodb://user:password@host:1/?ssl=true&ssl_ca_certs=cert-1&replicaSet=replica-1&readPreference=preference&retryWrites=true",
			stacktrace: "",
		},
	}

	for _, test := range tests {
		uri, err := test.config.generateURI()
		assert.Equal(t, test.output, uri)
		if test.hasError {
			assert.NotNil(t, err)
			assert.Equal(t, test.stacktrace, err.Error())
		} else {
			assert.Nil(t, err)
		}
	}
}
