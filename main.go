package main

import (
	"crossfs/cmd"
	"log"
)

func main() {
	log.SetFlags(log.Lmsgprefix)
	cmd.Execute()
}
