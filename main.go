package main

import (
	"log"

	swextreload "gopkg.teddywing.com/swextreload/internal"
)

func main() {
	err := swextreload.Reload("", "", true)
	if err != nil {
		log.Fatal(err)
	}
}
