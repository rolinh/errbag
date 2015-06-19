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
