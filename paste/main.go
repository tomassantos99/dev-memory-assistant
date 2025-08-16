package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tomassantos99/dev-memory-assistant/paste/handler"
	"github.com/tomassantos99/dev-memory-assistant/paste/listener"
	"github.com/tomassantos99/dev-memory-assistant/paste/ui"
)

func main() {
	fmt.Println("Starting dev-memory-assistant...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var window = ui.CreateHistoryWindow()
	defer window.Mw.Dispose()

	// Intercept Ctrl+C (SIGINT) and SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println("Received signal:", sig)
		window.Mw.Synchronize(func() {
			window.Mw.Dispose()
		})
		cancel()
	}()

	clipboardHandler := handler.NewClipboardHandler()
	pasteHandler := handler.NewPasteHandler(window)

	go clipboardHandler.ClipboardMessageHandler(ctx)

	go listener.ListenWindowsClipboardUpdates(ctx, clipboardHandler)
	go listener.ListenShortcuts(ctx, pasteHandler)

	fmt.Println("Listening for clipboard updates and shortcuts...")

	window.Mw.Run()

	fmt.Println("Terminating dev-memory-assistant...")
}
