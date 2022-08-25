package main

import (
	"crossfs/cmd"
	"log"
	"os"
)

func main() {
	err := os.Mkdir("build", 0666)
	if err != nil && err != os.ErrExist {
		log.Fatal(err)
	}

	err = cmd.GenerateDocs("build")
	if err != nil {
		log.Fatal(err)
	}
}
