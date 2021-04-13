**This repository is no longer maintained**

Instead, I recommend for others to use [grpcurl](https://github.com/fullstorydev/grpcurl) or [evans](https://github.com/ktr0731/evans).

# grpcurl

[![CircleCI](https://circleci.com/gh/kazegusuri/grpcurl.svg?style=svg)](https://circleci.com/gh/kazegusuri/grpcurl)

super experimental

# Installation

```
go get -u github.com/kazegusuri/grpcurl
```

# Usage

```
$ grpcurl
A handy and universal gRPC command line client

Usage:
  grpcurl [flags]
  grpcurl [command]

Available Commands:
  call          Call gRPC method with JSON
  help          Help about any command
  list_services List services and methods provided by gRPC server

Flags:
  -h, --help       help for grpcurl
  -k, --insecure   with insecure
  -v, --verbose    verbose output

Use "grpcurl [command] --help" for more information about a command.
```

### List services

```
$ grpcurl -k ls localhost:8080
test.EchoService
grpc.reflection.v1alpha.ServerReflection

$ grpcurl -k ls localhost:8080 test.EchoService
test.Test.Echo
```

### Call gRPC method

```
$ echo '{"Message": "hello"} | grpcurl -k call localhost:8080 test.EchoService.Echo
{"Message":"hello"}
```
