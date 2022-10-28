package main

import (
	"fmt"
	"huffman/src/flow"
	"os"
)

func main() {
	programFlow := flow.NewFlow()

	err := programFlow.Init()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
