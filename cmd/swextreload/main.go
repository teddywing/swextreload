package main

import (
	"fmt"
	"log"
	"os"

	"git.sr.ht/~liliace/claw"
	"github.com/dedelala/sysexits"
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
				Name:         "extension_id",
				Type:         "string",
				Repeating:    true,
				DefaultValue: []string{},
				Description:  "extensions to reload",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("args: %#v", args)

	socket_url, ok := args["socket-url"].(string)
	if !ok {
		fmt.Fprintln(os.Stderr, "error: '--socket-url' is required")
		os.Exit(sysexits.Usage)
	}

	extension_ids := args["extension_id"].([]string)
	if len(extension_ids) == 0 {
		fmt.Fprintln(os.Stderr, "error: missing extension IDs")
		os.Exit(sysexits.Usage)
	}

	err = swextreload.Reload(
		socket_url,
		extension_ids,
		true,
	)
	if err != nil {
		log.Fatal(err)
	}
}
