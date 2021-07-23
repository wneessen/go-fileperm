# fileperm - Check file permissions based on the current user in Go

[![Go Reference](https://pkg.go.dev/badge/github.com/wneessen/go-fileperm.svg)](https://pkg.go.dev/github.com/wneessen/go-fileperm) [![Go Report Card](https://goreportcard.com/badge/github.com/wneessen/go-fileperm)](https://goreportcard.com/report/github.com/wneessen/go-fileperm) 

In Go it is not trivial to check what permissions a specific file has to the current user. This module provides a small
library to work around this.

## Usage

First create a new `FileUserPerm` struct:

```go
package main

import (
	"fmt"
	"os"
	"github.com/wneessen/go-fileperm"
)

func main() {
	fup, err := fileperm.NewFileUserPerm("/var/tmp/foo.txt")
	if err != nil {
		fmt.Print("ERROR:", err)
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
	fup, err := fileperm.NewFileUserPerm("/var/tmp/foo.txt")
	if err != nil {
		fmt.Print("ERROR:", err)
		os.Exit(1)
	}
	fmt.Printf("User can write to file: %v", fup.UserWritable())
}
```
