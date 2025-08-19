package main

import (
	"context"
	"fmt"
	"github.com/tomassantos99/dev-memory-assistant/paste/handler"
	"github.com/tomassantos99/dev-memory-assistant/paste/listener"
	"github.com/tomassantos99/dev-memory-assistant/paste/storage"
	"github.com/tomassantos99/dev-memory-assistant/paste/ui"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Starting dev-memory-assistant...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	window := ui.CreateHistoryWindow(handler.PasteContent)
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

	storage.SetupDB()

	clipboardHandler := handler.NewClipboardHandler()
	pasteHandler := handler.NewShortcutHandler(window)

	go clipboardHandler.OnMessageHandler(ctx)

	go listener.ListenWindowsClipboardUpdates(ctx, clipboardHandler)
	go listener.ListenShortcuts(ctx, pasteHandler)

	fmt.Println("Listening for clipboard updates and shortcuts...")

	window.Mw.Run()

	fmt.Println("Terminating dev-memory-assistant...")
}
