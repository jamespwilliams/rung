# rung

A code generating utility to make writing Golang CLI apps easier.

## Example

`main.go`:

```golang
package main

import (
	"fmt"
	"io"
)

//go:generate sh flags.sh

func run(in io.Reader, out io.Writer, flags flagSet) {
	fmt.Println(flags.fooflag.Value)
	if flags.barflag.WasSet {
		fmt.Printf("do something with %v\n", flags.barflag.Value)
	}
}
```

`flags.sh`:

```sh
go run github.com/jamespwilliams/rung/rung_gen \
    -foo/-f:int fooflag 10 "the foo flag" \
    -bar/-b:bool barflag true "the bar flag" \
    -quux/-q:float64 quuxflag 10.241 "the quux flag"
```

(I'm using a separate script here because `go:generate` directives don't support
multiline commands.)

Then, for example:

```console
[jpw@xyz]$ go generate .

[jpw@xyz]$ go run . -foo 10000
10000

[jpw@xyz]$ go run . -foo 10000 -bar false
10000
do something with true

[jpw@xyz]$ go run . -quux 3.1415
10

[jpw@xyz]$ go run . -invalid
flag provided but not defined: -invalid
Usage of /run/user/1000/go-build723213117/b001/exe/rungtest:
  -bar
        the bar flag (default true)
  -foo int
        the foo flag (default 10)
  -quux float
        the quux flag (default 10.241)
```
