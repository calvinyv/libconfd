// Copyright 2018 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a Apache-style
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/chai2010/libconfd"
)

var logger libconfd.Logger = libconfd.NewStdLogger(os.Stderr, "", "", 0)

func SetLogger(l libconfd.Logger) {
	logger = l
}
