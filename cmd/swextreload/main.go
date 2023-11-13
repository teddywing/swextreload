package main

import (
	"log"

	swextreload "gopkg.teddywing.com/swextreload/internal"
)

func main() {
	err := swextreload.Reload(
		"ws://127.0.0.1:55755/devtools/browser/4536efdf-6ddf-40b6-9a16-258a1935d866",
		"imcibeelfmccdpnnlemllnepgbfdbkgo",
		true,
	)
	if err != nil {
		log.Fatal(err)
	}
}
