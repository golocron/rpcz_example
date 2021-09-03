// Package config provides a way to configure the cli app examples.
package config

import (
	"flag"
)

// Config holds configuration for a server or client.
type Config struct {
	Net        string
	Addr       string
	CertPath   string
	KeyPath    string
	ServerName string
	NReq       int
	TLS        bool
	Async      bool
}

// NewFromArgs creates a config with the values contained in args.
func NewFromArgs(name string, args []string) (*Config, error) {
	fset := flag.NewFlagSet(name, flag.ContinueOnError)
	result := &Config{}

	fset.StringVar(&result.Net, "net", "tcp", "network type: tcp or unix")
	fset.StringVar(&result.Addr, "addr", "localhost:10217", "listen on address: ip:port or path to socket")
	fset.StringVar(&result.CertPath, "cert", "", "tls: path to ssl certificate")
	fset.StringVar(&result.KeyPath, "key", "", "tls: path to ssl key")
	fset.StringVar(&result.ServerName, "ssl-server-name", "localhost", "tls client: name of the server (should match the one specified in the cert")
	fset.IntVar(&result.NReq, "nreq", 100, "client: number of requests to be made")
	fset.BoolVar(&result.TLS, "tls", false, "tls: use tls or not. Requeries cert and key")
	fset.BoolVar(&result.Async, "async", false, "client: run async mode")

	if err := fset.Parse(args); err != nil {
		return nil, err
	}

	if result.TLS {
		result.ensureCertKey()
	}

	return result, nil
}

func (c *Config) ensureCertKey() {
	if !c.TLS {
		return
	}

	if c.CertPath == "" && c.KeyPath == "" {
		c.CertPath = "misc/localhost.crt"
		c.KeyPath = "misc/localhost.key"
	}
}
