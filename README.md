<!--
SPDX-FileCopyrightText: 2022 Winni Neessen <winni@neessen.dev>

SPDX-License-Identifier: CC0-1.0
-->

# fileperm - An easy way to file permissions for the current user in Go

[![Go Reference](https://pkg.go.dev/badge/github.com/wneessen/go-fileperm.svg)](https://pkg.go.dev/github.com/wneessen/go-fileperm) 
[![Go Report Card](https://goreportcard.com/badge/github.com/wneessen/go-fileperm)](https://goreportcard.com/report/github.com/wneessen/go-fileperm)
[![codecov](https://codecov.io/gh/wneessen/go-fileperm/branch/main/graph/badge.svg?token=48AX0B6W7L)](https://codecov.io/gh/wneessen/go-fileperm)
[![REUSE status](https://api.reuse.software/badge/github.com/wneessen/go-fileperm)](https://api.reuse.software/info/github.com/wneessen/go-fileperm)
<a href="https://ko-fi.com/D1D24V9IX"><img src="https://uploads-ssl.webflow.com/5c14e387dab576fe667689cf/5cbed8a4ae2b88347c06c923_BuyMeACoffee_blue.png" height="20" alt="buy ma a coffee"></a>

In Go it is not trivial to check what permissions a specific file has to the current user. This module provides a small
library to work around this.

## Usage

First create a new `UserPerm` struct:

```go
package main

import (
	"fmt"
	"os"

	"github.com/wneessen/go-fileperm"
)

func main() {
	up, err := fileperm.New("/var/tmp/foo.txt")
	if err != nil {
		fmt.Print("failed to create new filepe:", err)
		os.Exit(1)
	}
}
```
Once the struct is ready at your hands, you can use different method check the permissions of the file for the
user of the current running process.

```go
package main

import (
	"fmt"
	"os"

	"github.com/wneessen/go-fileperm"
)

func main() {
	up, err := fileperm.New("/var/tmp/foo.txt")
	if err != nil {
		fmt.Print("ERROR:", err)
		os.Exit(1)
	}
	fmt.Printf("User can write to file: %t", up.UserWritable())
}
```