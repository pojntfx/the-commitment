package main

import "github.com/pojntfx/the-commitment/cmd/pledge/cmd"

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
