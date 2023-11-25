swextreload
===========

[![GoDoc](https://godocs.io/gopkg.teddywing.com/swextreload?status.svg)][Documentation]

Reload Chrome extensions from the command line. Facilitates Chrome extension
development.

Communicates with Chrome over the [DevTools Protocol].

This program is a rewrite of [Extreload] to add support for Manifest V3
extensions. Swextreload doesn’t work reliably with Manifest V2 extensions, so
you may want to continue using Extreload in that case.


## Usage
Chrome must be started with the `--remote-debugging-port` flag to enable the
DevTools Protocol, and the `--silent-debugger-extension-api` flag to allow debug
access to extensions. On Mac OS X:

	$ /Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome \
		--silent-debugger-extension-api \
		--remote-debugging-port=0 &

When Chrome is launched, the DevTools Protocol socket URL will be printed to the
console. That WebSocket URL must be passed to `swextreload` with the
`--socket-url` argument. For example:

	$ swextreload \
		--socket-url ws://127.0.0.1:55755/devtools/browser/208ae571-d691-4c98-ad41-3a15d507b656 \
		--reload-current-tab \
		ooeilikhhbbkljfdhbglpalaplegfcmj


## Install

	$ go install gopkg.teddywing.com/swextreload/cmd/swextreload@latest


## License
Copyright © 2023 Teddy Wing. Licensed under the GNU GPLv3+ (see the included
COPYING file).


[Documentation]: https://godocs.io/gopkg.teddywing.com/swextreload
[DevTools Protocol]: https://chromedevtools.github.io/devtools-protocol/
[Extreload]: https://github.com/teddywing/extreload/
