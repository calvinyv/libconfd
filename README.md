# libconfd

[![Build Status](https://travis-ci.org/openpitrix/libconfd.svg)](https://travis-ci.org/openpitrix/libconfd)
[![Go Report Card](https://goreportcard.com/badge/openpitrix.io/libconfd)](https://goreportcard.com/report/openpitrix.io/libconfd)
[![GoDoc](https://godoc.org/github.com/openpitrix/libconfd?status.svg)](https://godoc.org/github.com/openpitrix/libconfd)
[![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/openpitrix/libconfd/blob/master/LICENSE)

mini confd lib, based on [confd](https://github.com/kelseyhightower/confd)/[memkv](https://github.com/kelseyhightower/memkv)/[secconf](https://github.com/xordataexchange/crypt)/[logger](https://github.com/chai2010/logger).


## Example

```go
package main

import (
	"openpitrix.io/libconfd"
)

func main() {
	cfg := libconfd.MustLoadConfig("./confd.toml")
	client := libconfd.MustNewBackendClient("./confd-backend.toml")

	libconfd.NewProcessor().Run(cfg, client)
}
```

## miniconfd (only support toml backend)

```
$ go run miniconfd.go -h
```

See [miniconfd.go](miniconfd.go)
