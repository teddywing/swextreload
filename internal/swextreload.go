package swextreload


import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
)

// TODO
func Reload(
	url string,
	extensionID string,
	shouldReloadTab bool,
) error {
	allocatorContext, cancel := chromedp.NewRemoteAllocator(
		context.Background(),
		url,
	)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocatorContext)
	defer cancel()

	targets, err := chromedp.Targets(ctx)
	if err != nil {
		log.Fatalf("error: targets: %v", err)
	}

	log.Printf("Targets: %#v", targets)
	println()

	extensionURL := "chrome-extension://" + extensionID + "/"

	var targetID target.ID
	for _, target := range targets {
		if strings.HasPrefix(target.URL, extensionURL) {
			log.Printf("Target: %#v", target)
			targetID = target.TargetID
			break
		}
	}

	targetCtx, cancel := chromedp.NewContext(ctx, chromedp.WithTargetID(targetID))
	defer cancel()

	var runtimeResp []byte
	err = chromedp.Run(
		targetCtx,
		// chromedp.Evaluate(`chrome.runtime.reload();`, &runtimeResp),
		// chromedp.Evaluate(`chrome.tabs.reload();`, &tabsResp),
		// chromedp.Evaluate(`chrome.runtime.reload();`, nil),
		// chromedp.EvaluateAsDevTools(`chrome.runtime.reload();`, nil),
		chromedp.Evaluate(`chrome.runtime.reload();`, nil),
		// chromedp.Evaluate(`chrome.tabs.reload();`, nil),
	)
	if err != nil {
		log.Fatalf("error: run: %v", err)
	}

	// var tabsResp []byte
	// err = chromedp.Run(
	// 	targetCtx,
	// 	// chromedp.Evaluate(`chrome.tabs.reload();`, &tabsResp),
	// 	// chromedp.Evaluate(`chrome.tabs.reload();`, nil),
	// 	chromedp.EvaluateAsDevTools(`chrome.tabs.reload();`, nil),
	// )
	// if err != nil {
	// 	log.Fatalf("error: run tabs: %v", err)
	// }

	log.Printf("Runtime: %v", string(runtimeResp))
	// log.Printf("Tabs: %v", string(tabsResp))

	time.Sleep(200 * time.Millisecond)

	targets, err = chromedp.Targets(ctx)
	if err != nil {
		log.Fatalf("error: targets2: %v", err)
	}

	log.Printf("Targets: %#v", targets)
	println()

	for _, target := range targets {
		if strings.HasPrefix(target.URL, extensionURL) {
			log.Printf("Target: %#v", target)
			targetID = target.TargetID
			break
		}
	}

	targetCtx, cancel = chromedp.NewContext(ctx, chromedp.WithTargetID(targetID))
	defer cancel()

	var tabsResp []byte
	err = chromedp.Run(
		targetCtx,
		chromedp.Evaluate(`chrome.tabs.reload();`, nil),
	)
	if err != nil {
		log.Fatalf("error: run tabs: %v", err)
	}

	log.Printf("Tabs: %v", string(tabsResp))

	return nil
}
