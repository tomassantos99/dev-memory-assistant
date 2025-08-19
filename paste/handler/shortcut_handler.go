package handler

import (
	"fmt"
	"github.com/tomassantos99/dev-memory-assistant/paste/storage"
	"github.com/tomassantos99/dev-memory-assistant/paste/ui"
)



type ShortcutHandler struct {
	window *ui.HistoryWindow
}

func NewShortcutHandler(window *ui.HistoryWindow) *ShortcutHandler {
	return &ShortcutHandler{
		window: window,
	}
}

func (p *ShortcutHandler) HandleClipboardHistoryWindowShortcut() {
	if (p.window.IsVisible()){
		return
	}

	messages, err := storage.GetLastClipboardMessages(100)
	if err != nil {
		fmt.Println("Error fetching clipboard messages:", err)
		return
	}

	p.window.ShowHistoryWindow(messages)
}






