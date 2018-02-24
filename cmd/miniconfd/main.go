// Copyright 2018 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a Apache-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"

	_ "github.com/chai2010/libconfd"
)

func main() {
	flag.Parse()

	if err := Main(); err != nil {
		logger.Fatal(err)
	}

	fmt.Println("Done")
}

func Main() error {
	return fmt.Errorf("TODO")
}