package main

import (
	"log"

	"git.sr.ht/~liliace/claw"
	swextreload "gopkg.teddywing.com/swextreload/internal"
)

func main() {
	args, err := claw.Parse(&claw.Options{
		Name:        "swextreload",
		Description: "Reload Google Chrome extensions.",
		Flags: []claw.Flag{
			{
				LongName:    "socket-url",
				Type:        "string",
				Description: "DevTools protocol WebSocket URL",
			},
			{
				LongName:    "reload-current-tab",
				Type:        "bool",
				Description: "pass this to reload the active Chrome tab",
			},
			{
				LongName:    "version",
				ShortName:   'V',
				Type:        "bool",
				Description: "show the program version",
			},
		},
		Positionals: []claw.Positional{
			{
				Name:        "extension_id",
				Type:        "string",
				Repeating:   true,
				Description: "extensions to reload",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("args: %#v", args)
	return

	err = swextreload.Reload(
		"ws://127.0.0.1:55755/devtools/browser/4536efdf-6ddf-40b6-9a16-258a1935d866",
		"imcibeelfmccdpnnlemllnepgbfdbkgo",
		true,
	)
	if err != nil {
		log.Fatal(err)
	}
}
