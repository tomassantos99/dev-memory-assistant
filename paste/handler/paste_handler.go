package handler

import (
	"github.com/tomassantos99/dev-memory-assistant/paste/ui"
)



type PasteHandler struct {
	window *ui.HistoryWindow
}

func NewPasteHandler(window *ui.HistoryWindow) *PasteHandler {
	return &PasteHandler{
		window: window,
	}
}

func (p *PasteHandler) HandlePaste() {
	items := []string{"First", "Second", "Third"} // TODO: Replace with actual clipboard data retrieval logic
	p.window.ShowHistoryWindow(items) // TODO: abstract windows/unix ui
}
