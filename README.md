<!--
SPDX-FileCopyrightText: 2022 Winni Neessen <winni@neessen.dev>

SPDX-License-Identifier: CC0-1.0
-->

# fileperm - file permission tests for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/wneessen/go-fileperm.svg)](https://pkg.go.dev/github.com/wneessen/go-fileperm)
[![Go Report Card](https://goreportcard.com/badge/github.com/wneessen/go-fileperm)](https://goreportcard.com/report/github.com/wneessen/go-fileperm)
[![codecov](https://codecov.io/gh/wneessen/go-fileperm/branch/main/graph/badge.svg?token=48AX0B6W7L)](https://codecov.io/gh/wneessen/go-fileperm)
[![REUSE status](https://api.reuse.software/badge/github.com/wneessen/go-fileperm)](https://api.reuse.software/info/github.com/wneessen/go-fileperm)
<a href="https://ko-fi.com/D1D24V9IX"><img src="https://uploads-ssl.webflow.com/5c14e387dab576fe667689cf/5cbed8a4ae2b88347c06c923_BuyMeACoffee_blue.png" height="20" alt="buy ma a coffee"></a>

## The problem
Go does not offer a trivial way to get the permissions of a specific file in the context of the current user.
While we get tools like [os.Stat](https://pkg.go.dev/os#Stat) to get general file permissions, this does not 
help to identify how our current user is affected by these permissions.

Other languages like the Shell or Perl provide a simple test mechanism, to allow the program to validate if the
user, that the program is executed as does have permissions on the specific file.

Sh:
```shell
$ chmod 400 test.txt && test -r test.txt && echo "We can read the file" || echo "We cannot read the file"
We can read the file

$ chmod 100 test.txt && test -r test.txt && echo "We can read the file" || echo "We cannot read the file"
We cannot read the file
```

Perl:
```perl
#!/usr/bin/env perl

if(-r "./test.txt") {
    print "We can"
}
else {
    print "We cannot"
}
print " read the file\n"';
```

go-fileperm aims to provide a similar mechanism in Go. The module offers methods to check if a provided file is
accessible to the user (readable, writable, executable and combinations of those).

## Usage

First, we create a new `UserPerm` instance by calling the `New` method together with the file in question as an 
argument:

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

Once the `UserPerm` instance is ready in your hands, we can use different methods to check the permissions of the file 
for the context of the user that our program or process is running in:

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
	fmt.Printf("User can write to file: %t", up.UserWritable())
}
```

## Performance

go-fileperm is quite fast and works allocation free. The single checks like `UserReadable`, `UserWritable` and
`UserExecutable` works in the 140-150ns range, while the combined ones take a little longer (180-210ns).

```
goos: darwin
goarch: arm64
pkg: github.com/wneessen/go-fileperm
BenchmarkPermUser_UserReadable
BenchmarkPermUser_UserReadable-8                 7364846               143.6 ns/op             0 B/op          0 allocs/op
BenchmarkPermUser_UserWritable
BenchmarkPermUser_UserWritable-8                 7803267               154.9 ns/op             0 B/op          0 allocs/op
BenchmarkPermUser_UserExecutable
BenchmarkPermUser_UserExecutable-8               7922624               149.2 ns/op             0 B/op          0 allocs/op
BenchmarkPermUser_UserWriteReadable
BenchmarkPermUser_UserWriteReadable-8            6494815               186.1 ns/op             0 B/op          0 allocs/op
BenchmarkPermUser_UserWriteExecutable
BenchmarkPermUser_UserWriteExecutable-8          6590229               181.0 ns/op             0 B/op          0 allocs/op
BenchmarkPermUser_UserReadExecutable
BenchmarkPermUser_UserReadExecutable-8           6190532               184.7 ns/op             0 B/op          0 allocs/op
BenchmarkPermUser_UserWriteReadExecutable
BenchmarkPermUser_UserWriteReadExecutable-8      5728713               208.8 ns/op             0 B/op          0 allocs/op
PASS
```