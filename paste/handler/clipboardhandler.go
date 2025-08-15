package handler

import "github.com/tomassantos99/dev-memory-assistant/paste/storage"

type ClipboardHandler struct {
	ClipboardMessages chan string
}

func NewClipboardHandler() *ClipboardHandler {
	return &ClipboardHandler{
		ClipboardMessages: make(chan string, 256),
	}
}

func (clipboardHandler *ClipboardHandler) Run() {
	for message := range clipboardHandler.ClipboardMessages {
		err := storage.SaveClipboardMessage(message) // direct ref for now, could be abstracted later
		if err != nil {
			println("Error saving clipboard message:", err.Error())
		}
	}
}
