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

	ctx, cancel := chromedp.NewContext(allocatorContext)
	defer cancel()

	// TODO: I think get targets once first, and reload all extensions using those targets. Rather than getting targets for each extension reload.
	targets, err := chromedp.Targets(ctx)
	if err != nil {
		return fmt.Errorf("swextreload: can't get targets: %v", err)
	}

	if isDebug {
		log.Printf("Targets: %#v", targets)
	}

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

		// TODO: Do the reload of the current page after reloading all
		// extensions. The current system doesn't work well with multiple
		// extensions.
	}

	if shouldReloadTab {
		time.Sleep(200 * time.Millisecond)

		extensionURL := "chrome-extension://" + extensionIDs[0] + "/"

		var firstExtensionTarget *target.Info
		for _, target := range targets {
			if strings.HasPrefix(target.URL, extensionURL) {
				firstExtensionTarget = target

				break
			}
		}

		err = reloadTab(
			allocatorContext,
			extensionIDs[0],
			firstExtensionTarget,
			isExtensionManifestV2(firstExtensionTarget),
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
			if isDebug {
				log.Printf("Target: %#v", target)
			}

			targetCtx, cancel := chromedp.NewContext(ctx, chromedp.WithTargetID(target.TargetID))
			defer cancel()

			log.Printf("Connected to target '%s'", target.TargetID)

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

			log.Printf("Reloaded extension")

			if isDebug {
				log.Printf("Runtime: %v", string(runtimeResp))
			}
		}
	}

	return nil
}

func reloadTab(
	ctx context.Context,
	extensionID string,
	letarget *target.Info,
	isExtensionManifestV2 bool,
) error {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	if isDebug {
		log.Printf("Reload tab (Manifest V2: %t)", isExtensionManifestV2)
	}

	if !isExtensionManifestV2 {
		// TODO: If MV2, then don't re-attach, only do it if "service_worker"
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

		ctx, cancel = chromedp.NewContext(ctx, chromedp.WithTargetID(targetID))
		defer cancel()
	}

	if isDebug {
		log.Printf("Connecting to target %s", letarget.TargetID)
	}

	ctx, cancel = chromedp.NewContext(ctx, chromedp.WithTargetID(letarget.TargetID))
	// defer cancel()

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

	if isDebug {
		log.Printf("Tabs: %v", string(tabsResp))
	}

	return nil
}

func isExtensionManifestV2(target *target.Info) bool {
	return target.Type == "background_page"
}
