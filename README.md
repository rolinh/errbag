# errbag - An error throttler for Go

![License](https://img.shields.io/badge/license-BSD-blue.svg)
[![Build Status](https://travis-ci.org/Rolinh/errbag.png?branch=master)](https://travis-ci.org/Rolinh/errbag)
[![GoDoc](http://godoc.org/github.com/Rolinh/errbag?status.svg)](http://godoc.org/github.com/Rolinh/errbag)

Is your Go (golang) application running for long period of times or as a daemon?
If the answer is yes, chances are that you might find `errbag` useful.

`Errbag` is a Go package that allows Go programs to pause for a predefined
amount of time when too many errors happen.

Errors are usually unavoidable for programs or applications which are supposed
to run continuously for long period of times. In these case, errors are dealt
with and eventually logged.  What if, all of a sudden, tens or hundreds of
errors are raised each second?  This usually indicates something troublesome is
happening to the machine running the application. It could be the network being
down or a disk partition being full but this obviously depends on the
application. Anyway, in this case it would be advised not to log errors like
crazy and keep on going but pause the application in order to let the sysadmin
deal with the underlying problem. This is exactly to solve this problem that
`errbag` has been created.

## Installation and usage

Get the package first: `go get github.com/Rolinh/errbag`
And import it into your project.

First thing needed is to create a new `Errbag` by invoking `New()` and then
inflate it using `Inflate()`. From this point, you can record errors using
`Record()`. This method optionally takes a callback function as argument. When
too many errors are recorded, it will block for a defined amount of time before
returning.
When done, do not forget to `Deflate()` the `Errbag` (I would advise to call
`Deflate()` in a `defer` statement).

Here is a full (stupid) example:

```go
package main

import (
	"fmt"
	"os"

	"github.com/Rolinh/errbag"
)

func main() {
	var waitTime, errBagSize, leakInterval uint
	waitTime, errBagSize, leakInterval = 5, 60, 1000

	bag, err := errbag.New(waitTime, errBagSize, leakInterval)
	if err != nil {
		fmt.Fprintln(os.Stderr, "impossible to create a new errbag:", err)
		os.Exit(1)
	}
	bag.Inflate()
	defer bag.Deflate()

	// some function which potentially returns a non nil error
	// (use your imagination...)
	foo := func() error {
		var err error
		return err
	}

	// we do a lots of calls to foo() later on...
	for i := 0; i < 10000; i++ {
		bag.Record(foo(), func(status errbag.Status) {
			if status.State != errbag.StatusOK {
				// Damn, we are throttling for status.WaitTime time.
				// Take the appropriate action, like sending an email to the
				// sysadmin or whatever :)
				// This is obviously optional. If you don't care, simply pass
				// nil as a second argument to Record().
			}
		})
	}
}
```
