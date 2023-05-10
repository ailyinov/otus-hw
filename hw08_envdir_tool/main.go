package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("not enough args")
		return
	}
	args := os.Args[1:]

	env, err := ReadDir(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	os.Exit(RunCmd(args[1:], env))
}
