# rpcz_example

RPCz Example is a small project that demonstrates a way of using the [`rpcz`](https://github.com/golocron/rpcz) library for serving and talking to an RPC service.

The code in this repository is intended to be self-explanatory, and shows a few basic concepts:

- [./proto](./proto) contains the Protobuf definitions of the request and response types used by the Echo service
- [./service/echo](./service/echo) is an example of an RPC service which makes use of the type definitions above
- [./cmd/echo-server](./cmd/echo-server) creates and runs an RPCz server with an instance of the Echo service
- [./cmd/echo-client](./cmd/echo-client) shows an RPCz client that makes requests to the Echo service that is run by the RPCz server.


## Server

To build and run the server app, use the following commands:

```bash
make echo-server

./bin/echo-server
```


## Client

The client app can be built and run as simple as:

```bash
make echo-client

./bin/echo-client
```


## More Options

There are a few more options for running both client and server apps:

```bash
$ ./bin/echo-client -help

Usage of ./bin/echo-client:
  -addr string
    	listen on address: ip:port or path to socket (default "localhost:10217")
  -async
    	client: run in async mode
  -cert string
    	tls: path to ssl certificate
  -key string
    	tls: path to ssl key
  -net string
    	network type: tcp or unix (default "tcp")
  -nreq int
    	client: number of requests to be made (default 100)
  -ssl-server-name string
    	tls client: name of the server (should match the one specified in the cert (default "localhost")
  -tls
    	tls: use tls or not. Requeries cert and key
```
