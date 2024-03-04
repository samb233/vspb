package main

import (
	"fmt"
	"vspb"
)

func main() {
	conf, err := vspb.ReadConfig("")
	if err != nil {
		panic(err)
	}

	fmt.Println(conf.Packages[0].Env)

	fmt.Printf("%v", conf)
}
