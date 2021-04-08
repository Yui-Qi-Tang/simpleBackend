package main

import (

	"simpleBackend/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
