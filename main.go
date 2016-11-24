package main

import (
	"fmt"

	"github.com/vdemeester/clasp/cmd"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		fmt.Printf("%v\n", err)
	}
}
