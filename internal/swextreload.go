// Package swextreload enables reloading Chrome extensions.
package swextreload

import (
	"context"
	"errors"
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

	allocatorContext, _ := chromedp.NewRemoteAllocator(
		context.Background(),
		url,
	)
	allocatorContext, cancel := context.WithTimeout(
		allocatorContext,
		5*time.Second,
	)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocatorContext)
	defer cancel()

	targets, err := chromedp.Targets(ctx)
	if err != nil {
		return fmt.Errorf("swextreload: can't get targets: %v", err)
	}

	logDebugf("Targets: %#v", targets)

	for _, extensionID := range extensionIDs {
		err = reloadExtension(
			ctx,
			targets,
			extensionID,
			shouldReloadTab,
		)
		if err != nil {
			return err
		}
	}

	if shouldReloadTab {
		extensionURL := "chrome-extension://" + extensionIDs[0] + "/"

		var firstExtensionTarget *target.Info
		for _, target := range targets {
			if strings.HasPrefix(target.URL, extensionURL) {
				firstExtensionTarget = target

				logDebugf(
					"firstExtensionTarget %s: %#v",
					extensionURL,
					firstExtensionTarget,
				)

				break
			}
		}

		if firstExtensionTarget == nil {
			// TODO: continue loop until target != null
			return errors.New("swextreload: can't reload tab, no target available")
		}

		// In Manifest V3, we need to wait until the service worker reinstalls
		// before we can re-attach to it.
		if !isExtensionManifestV2(firstExtensionTarget) {
			time.Sleep(200 * time.Millisecond)
		}

		err = reloadTab(
			allocatorContext,
			extensionIDs[0],
			firstExtensionTarget,
		)
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
	targets []*target.Info,
	extensionID string,
	shouldReloadTab bool,
) error {
	extensionURL := "chrome-extension://" + extensionID + "/"

	for _, target := range targets {
		if strings.HasPrefix(target.URL, extensionURL) {
			logDebugf("Target: %#v", target)

			targetCtx, cancel := chromedp.NewContext(ctx, chromedp.WithTargetID(target.TargetID))
			defer cancel()

			logDebugf("Connected to target '%s'", target.TargetID)

			var runtimeResp []byte
			err := chromedp.Run(
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

			logDebugf("Reloaded extension")

			logDebugf("Runtime: %v", string(runtimeResp))
		}
	}

	return nil
}

// reloadTab reloads the current Chrome tab using the background console of
// either extensionID or reloadTarget.
func reloadTab(
	ctx context.Context,
	extensionID string,
	reloadTarget *target.Info,
) error {
	// Don't cancel the context. Otherwise, the background page DevTools
	// window closes.
	ctx, cancel := chromedp.NewContext(ctx)

	isMV2 := isExtensionManifestV2(reloadTarget)
	logDebugf("Reload tab (Manifest V2: %t)", isMV2)

	// If the extension is Manifest V3, its `targetId` reset after we reloaded
	// the extension from the service worker, presumably because it was
	// reinstalled. In that case, we need to get targets again, find the new
	// `targetId`, and connect to it.
	//
	// If the extension is Manifest V2, we can just reconnect to the existing
	// target.
	if !isMV2 {
		targets, err := chromedp.Targets(ctx)
		if err != nil {
			return fmt.Errorf(
				"swextreload: can't get targets for '%s' tab reload: %v",
				extensionID,
				err,
			)
		}

		logDebugf("Targets: %#v", targets)

		extensionURL := "chrome-extension://" + extensionID + "/"

		var targetID target.ID
		for _, target := range targets {
			if strings.HasPrefix(target.URL, extensionURL) {
				logDebugf("Target: %#v", target)

				targetID = target.TargetID
				break
			}
		}

		ctx, cancel = chromedp.NewContext(ctx, chromedp.WithTargetID(targetID))
		defer cancel()
	} else {
		logDebugf("Connecting to target %s", reloadTarget.TargetID)

		// Don't cancel the context. Otherwise, the background page DevTools
		// window closes.
		ctx, _ = chromedp.NewContext(
			ctx,
			chromedp.WithTargetID(reloadTarget.TargetID),
		)
	}

	var tabsResp []byte
	err := chromedp.Run(
		ctx,
		chromedp.Evaluate(`chrome.tabs.reload();`, nil),
	)
	if err != nil {
		return fmt.Errorf(
			"swextreload: error reloading tab '%s': %v",
			extensionID,
			err,
		)
	}

	logDebugf("Tabs: %v", string(tabsResp))

	return nil
}

// isExtensionManifestV2 returns true if target is a Manifest V2 extension.
func isExtensionManifestV2(target *target.Info) bool {
	return target.Type == "background_page"
}

// logDebugf prints a debug log if isDebug is on.
func logDebugf(format string, v ...any) {
	if !isDebug {
		return
	}

	log.Printf(format, v...)
}
