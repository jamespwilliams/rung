# rung

A code generating utility to make Go CLI flag parsing less painful.

## Example

Firstly, the Go code (`main.go`):

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

Then we just need a script to call the generation code (in `flags.sh`):

```sh
go run github.com/jamespwilliams/rung/rung_gen \
    -foo:int fooflag 10 "the foo flag" \
    -bar:bool barflag true "the bar flag" \
    -quux:float64 quuxflag 10.241 "the quux flag"
```

That's it.

Then, we can run for example:

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

Note that I'm using a separate `flags.sh` script because `go:generate`
directives don't support multiline commands.

See [this article by Mat
Ryer](https://pace.dev/blog/2020/02/12/why-you-shouldnt-use-func-main-in-golang-by-mat-ryer.html)
for an explanation of why abstracting `main` out into a separate `run` function
is useful.

## Usage

Really, this is more of a proof-of-concept, but anyway: the `rung_gen` binary
accepts argument quadruples of the following form:

```
-flag:type name default "usage string"
```

## Future Ideas

The expected `run` function signature could be changed such that an `error` is
expected, which could then be `log.Fatal`'d, for example.
