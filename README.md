# grpcurl

super experimental

# Installation

```
go get -u github.com/kazegusuri/grpcurl
```

# Usage

```
$ echo '{"Message": "hello"} | grpcurl -k call localhost:5000 test.EchoService Echo
{"Message":"hello"}
```
