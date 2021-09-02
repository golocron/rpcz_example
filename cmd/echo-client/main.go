package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/golocron/rpcz"

	"github.com/golocron/rpcz_example/config"
	"github.com/golocron/rpcz_example/proto/echo"
)

var (
	errUnknownMode = errors.New("unknown network mode")
	errInvalidCert = errors.New("invalid tls certificate")
)

func main() {
	lg := log.New(os.Stdout, "", log.LstdFlags)

	if err := setupAndRun(lg, os.Args[0], os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}

		lg.Printf("client error occurred: %s", err)
		os.Exit(1)
	}

	lg.Println("client finished")
	os.Exit(0)
}

type asyncResult struct {
	sent  string
	recvd string
	err   error
}

func setupAndRun(lg *log.Logger, name string, args []string) error {
	cfg, err := config.NewFromArgs(name, args)
	if err != nil {
		return err
	}

	client, err := newClient(cfg)
	if err != nil {
		return err
	}

	if cfg.Async {
		doAsync(cfg, client, lg)
	} else {
		doSync(cfg, client, lg)
	}

	return client.Close()
}

func doSync(cfg *config.Config, c *rpcz.Client, lg *log.Logger) {
	ctx := context.TODO()

	var num int
	for i := 0; i < cfg.NReq; i++ {
		req := &echo.EchoRequest{Msg: strconv.Itoa(i)}
		resp := &echo.EchoResponse{}

		if err := c.SyncDo(ctx, "Echo", "Echo", req, resp); err != nil {
			lg.Printf("failed to make request: %v", err)
			continue
		}

		lg.Printf("sent => %s; received => %s", req.Msg, resp.Msg)
		num++
	}

	lg.Printf("successfully made: %d out of %d requests", num, cfg.NReq)
}

func doAsync(cfg *config.Config, c *rpcz.Client, lg *log.Logger) {
	ctx := context.TODO()

	wg := &sync.WaitGroup{}
	out := make(chan *asyncResult, cfg.NReq/2)

	go func() {
		wg.Wait()
		close(out)
	}()

	for i := 0; i < cfg.NReq; i++ {
		wg.Add(1)

		go func(lctx context.Context, lwg *sync.WaitGroup, m string) {
			defer lwg.Done()

			req := &echo.EchoRequest{Msg: m}
			resp := &echo.EchoResponse{}

			fut := c.Do(ctx, "Echo", "Echo", req, resp)

			select {
			case <-ctx.Done():
				out <- &asyncResult{err: ctx.Err()}
			case err := <-fut.ErrChan():
				out <- &asyncResult{sent: m, recvd: resp.Msg, err: err}
			}
		}(ctx, wg, strconv.Itoa(i))
	}

	var num int
	for v := range out {
		if v.err != nil {
			lg.Printf("failed to make request: %v", v.err)
			continue
		}

		lg.Printf("sent => %s; received => %s", v.sent, v.recvd)
		num++
	}

	lg.Printf("successfully handled: %d out of %d requests", num, cfg.NReq)
}

func newClient(cfg *config.Config) (*rpcz.Client, error) {
	switch cfg.Net {
	case "tcp":
		return newClientTCP(cfg)
	case "unix":
		return newClientUnix(cfg)
	default:
		return nil, errUnknownMode
	}
}

func newClientTCP(cfg *config.Config) (*rpcz.Client, error) {
	if !cfg.TLS {
		return rpcz.NewClient(cfg.Addr)
	}

	tcfg, err := newTLSConfig(cfg)
	if err != nil {
		return nil, err
	}

	return rpcz.NewClientTLS(cfg.Addr, tcfg)
}

func newClientUnix(cfg *config.Config) (*rpcz.Client, error) {
	nc, err := net.Dial(cfg.Net, cfg.Addr)
	if err != nil {
		return nil, err
	}

	if cfg.TLS {
		tcfg, err := newTLSConfig(cfg)
		if err != nil {
			return nil, err
		}

		nc = tls.Client(nc, tcfg)
	}

	return rpcz.NewClientWithConn(&rpcz.ClientOptions{Encoding: rpcz.Protobuf}, nc)
}

func newTLSConfig(cfg *config.Config) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(cfg.CertPath, cfg.KeyPath)
	if err != nil {
		return nil, err
	}

	if len(cert.Certificate) == 0 {
		return nil, errInvalidCert
	}

	xcert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, err
	}

	rcert := x509.NewCertPool()
	rcert.AddCert(xcert)

	return &tls.Config{RootCAs: rcert, ServerName: cfg.ServerName}, nil
}
