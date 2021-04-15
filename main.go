package main

import (
	"simpleBackend/cmd"
	"syscall"
)

// set the open file to ulimit(just for test!!)
func setUlimit() error {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		return err
	}
	rLimit.Cur = rLimit.Max
	// rLimit.Cur = 1024

	return syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
}

func main() {

	/*  just test
	if err := setUlimit(); err != nil {
		panic(err)
	}
	*/

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
