// Package swextreload enables reloading Chrome extensions.
package swextreload

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
)

// isDebug controls whether debug printing is enabled.
var isDebug = false

// SetDebugOn turns on debug printing.
func SetDebugOn() {
	isDebug = true
}

// Reload reloads the extensions in extensionIDs.
func Reload(
	url string,
	extensionIDs []string,
	shouldReloadTab bool,
) error {
	var err error

	allocatorContext, cancel := chromedp.NewRemoteAllocator(
		context.Background(),
		url,
	)
	defer cancel()

	for _, extensionID := range extensionIDs {
		err = reloadExtension(
			allocatorContext,
			extensionID,
			shouldReloadTab,
		)
		if err != nil {
			return err
		}

		// TODO: Do the reload of the current page after reloading all
		// extensions. The current system doesn't work well with multiple
		// extensions.
	}

	if shouldReloadTab {
		time.Sleep(200 * time.Millisecond)

		err = reloadTab(allocatorContext, extensionIDs[0])
		if err != nil {
			return err
		}
	}

	return nil
}

// reloadExtension reloads the extension extensionID. If shouldReloadTab is
// true, also reload the current tab.
func reloadExtension(
	ctx context.Context,
	extensionID string,
	shouldReloadTab bool,
) error {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	targets, err := chromedp.Targets(ctx)
	if err != nil {
		return fmt.Errorf("swextreload: can't get targets: %v", err)
	}

	if isDebug {
		log.Printf("Targets: %#v", targets)
	}

	extensionURL := "chrome-extension://" + extensionID + "/"

	var targetID target.ID
	for _, target := range targets {
		if strings.HasPrefix(target.URL, extensionURL) {
			if isDebug {
				log.Printf("Target: %#v", target)
			}

			targetID = target.TargetID
			break
		}
	}

	targetCtx, cancel := chromedp.NewContext(ctx, chromedp.WithTargetID(targetID))
	defer cancel()

	var runtimeResp []byte
	err = chromedp.Run(
		targetCtx,
		chromedp.Evaluate(`chrome.runtime.reload();`, nil),
	)
	if err != nil {
		return fmt.Errorf(
			"swextreload: error reloading extension '%s': %v",
			extensionID,
			err,
		)
	}

	if isDebug {
		log.Printf("Runtime: %v", string(runtimeResp))
	}

	return nil
}

func reloadTab(ctx context.Context, extensionID string) error {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	targets, err := chromedp.Targets(ctx)
	if err != nil {
		return fmt.Errorf(
			"swextreload: can't get targets for '%s' tab reload: %v",
			extensionID,
			err,
		)
	}

	if isDebug {
		log.Printf("Targets: %#v", targets)
	}

	extensionURL := "chrome-extension://" + extensionID + "/"

	var targetID target.ID
	for _, target := range targets {
		if strings.HasPrefix(target.URL, extensionURL) {
			if isDebug {
				log.Printf("Target: %#v", target)
			}

			targetID = target.TargetID
			break
		}
	}

	targetCtx, cancel := chromedp.NewContext(ctx, chromedp.WithTargetID(targetID))
	defer cancel()

	var tabsResp []byte
	err = chromedp.Run(
		targetCtx,
		chromedp.Evaluate(`chrome.tabs.reload();`, nil),
	)
	if err != nil {
		return fmt.Errorf(
			"swextreload: error reloading tab '%s': %v",
			extensionID,
			err,
		)
	}

	if isDebug {
		log.Printf("Tabs: %v", string(tabsResp))
	}

	return nil
}
