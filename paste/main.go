package main

import (
	"github.com/tomassantos99/dev-memory-assistant/paste/listener"
	"github.com/tomassantos99/dev-memory-assistant/paste/handler"
)

func main() {
	clipboardHandler := handler.NewClipboardHandler()

	defer close(clipboardHandler.ClipboardMessages)
	go clipboardHandler.Run()
	
	listener.ListenWindowsClipboardUpdates(clipboardHandler)
}
