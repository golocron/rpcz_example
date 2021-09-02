package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"log"
	"net"
	"os"

	"github.com/golocron/daemon"
	"github.com/golocron/rpcz"

	"github.com/golocron/rpcz_example/config"
	"github.com/golocron/rpcz_example/service/echo"
)

var (
	errUnknownMode = errors.New("unknown network mode")
)

func main() {
	lg := log.New(os.Stdout, "", log.LstdFlags)

	if err := setupAndRun(lg, os.Args[0], os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}

		lg.Printf("server error occurred: %s", err)
		os.Exit(1)
	}

	lg.Println("server finished")
	os.Exit(0)
}

func setupAndRun(lg *log.Logger, name string, args []string) error {
	cfg, err := config.NewFromArgs(name, args)
	if err != nil {
		return err
	}

	srv, err := newServer(cfg)
	if err != nil {
		return err
	}

	svc := echo.New()
	if err := srv.Register(svc); err != nil {
		return err
	}

	lg.Printf("starting server")

	d := daemon.New(srv)

	return d.Start()
}

func newServer(cfg *config.Config) (*rpcz.Server, error) {
	switch cfg.Net {
	case "tcp":
		return newServerTCP(cfg)
	case "unix":
		return newServerUnix(cfg)
	default:
		return nil, errUnknownMode
	}
}

func newServerTCP(cfg *config.Config) (*rpcz.Server, error) {
	if !cfg.TLS {
		return rpcz.NewServer(cfg.Addr)
	}

	cert, err := tls.LoadX509KeyPair(cfg.CertPath, cfg.KeyPath)
	if err != nil {
		return nil, err
	}

	tcfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	return rpcz.NewServerTLS(cfg.Addr, tcfg)
}

func newServerUnix(cfg *config.Config) (*rpcz.Server, error) {
	ln, err := net.Listen(cfg.Net, cfg.Addr)
	if err != nil {
		return nil, err
	}

	if cfg.TLS {
		cert, err := tls.LoadX509KeyPair(cfg.CertPath, cfg.KeyPath)
		if err != nil {
			return nil, err
		}

		tcfg := &tls.Config{Certificates: []tls.Certificate{cert}}
		ln = tls.NewListener(ln, tcfg)
	}

	return rpcz.NewServerWithListener(&rpcz.ServerOptions{Encoding: rpcz.Protobuf}, ln)
}
