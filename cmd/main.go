package main

import (
	"os"
	"vspb"
)

func main() {
	confPath := "config.yml"
	if len(os.Args) > 1 {
		confPath = os.Args[1]
	}

	if err := vspb.Run(confPath); err != nil {
		panic(err)
	}
}
