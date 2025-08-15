package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/tomassantos99/dev-memory-assistant/paste/handler"
	"github.com/tomassantos99/dev-memory-assistant/paste/listener"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Intercept Ctrl+C (SIGINT) and SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println("Received signal:", sig)
		cancel()
	}()

	clipboardHandler := handler.NewClipboardHandler()
	go clipboardHandler.ClipboardMessageHandler(ctx)

	go listener.ListenWindowsClipboardUpdates(ctx, clipboardHandler.ClipboardMessages)
	go listener.ListenShortcuts(ctx)

	<- ctx.Done() // Wait for context cancellation
	fmt.Println("Exiting application...")
}
