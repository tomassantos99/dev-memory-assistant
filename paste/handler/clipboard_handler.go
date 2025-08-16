package handler

import (
	"context"
	"github.com/tomassantos99/dev-memory-assistant/paste/storage"
)

type ClipboardHandler struct {
	ClipboardMessages chan string
}

func NewClipboardHandler() *ClipboardHandler {
	return &ClipboardHandler{
		ClipboardMessages: make(chan string, 256),
	}
}

func (clipboardHandler *ClipboardHandler) OnMessageHandler(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(clipboardHandler.ClipboardMessages)
			println("Clipboard message handler stopped")
			return
		case message, ok := <-clipboardHandler.ClipboardMessages:
			if !ok {
				return
			}
			err := storage.SaveClipboardMessage(message) // direct ref for now, could be abstracted later
			if err != nil {
				println("Error saving clipboard message:", err.Error())
			}
		}
	}
}
