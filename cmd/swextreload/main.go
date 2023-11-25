// Copyright (c) 2023  Teddy Wing
//
// This file is part of Swextreload.
//
// Swextreload is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Swextreload is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Swextreload. If not, see <https://www.gnu.org/licenses/>.


package main

import (
	"fmt"
	"log"
	"os"

	"git.sr.ht/~liliace/claw"
	"github.com/dedelala/sysexits"
	swextreload "gopkg.teddywing.com/swextreload/internal"
)

const programVersion = "0.0.1"

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
				LongName:     "reload-current-tab",
				Type:         "bool",
				DefaultValue: false,
				Description:  "pass this to reload the active Chrome tab",
			},
			{
				LongName:    "version",
				ShortName:   'V',
				Type:        "bool",
				DefaultValue: false,
				Description: "show the program version",
			},
			{
				LongName:     "debug",
				Type:         "bool",
				DefaultValue: false,
				Description:  "print debug output",
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
		fmt.Fprintf(
			os.Stderr,
			"error: failed to parse command line arguments: %v\n",
			err,
		)
		os.Exit(sysexits.DataErr)
	}

	version := args["version"].(bool)
	if version {
		fmt.Println(programVersion)
		os.Exit(sysexits.OK)
	}

	socket_url, ok := args["socket-url"].(string)
	if !ok {
		fmt.Fprintln(os.Stderr, "error: '--socket-url' is required")
		os.Exit(sysexits.Usage)
	}

	shouldReloadTab := args["reload-current-tab"].(bool)

	extension_ids := args["extension_id"].([]string)
	if len(extension_ids) == 0 {
		fmt.Fprintln(os.Stderr, "error: missing extension IDs")
		os.Exit(sysexits.Usage)
	}

	isDebug := args["debug"].(bool)
	if isDebug {
		swextreload.SetDebugOn()
	}

	if isDebug {
		log.Printf("args: %#v", args)
	}

	err = swextreload.Reload(
		socket_url,
		extension_ids,
		shouldReloadTab,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: can't reload extension: %v\n", err)
		os.Exit(sysexits.Unavailable)
	}
}
