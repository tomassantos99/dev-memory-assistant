package listener

import (
	"context"
	"unsafe"
	"github.com/tomassantos99/dev-memory-assistant/paste/handler"
	"github.com/tomassantos99/dev-memory-assistant/paste/pkg"
)

const (
	MOD_CONTROL = 0x0002
	MOD_SHIFT   = 0x0004
	VK_V        = 0x56
	WM_HOTKEY   = 0x0312
	WM_QUIT     = 0
)

func ListenShortcuts(ctx context.Context, pasteHander *handler.PasteHandler) error {
	// Register Ctrl+Shift+V as a global hotkey (id=1)
	ret, _, err := registerHotKey.Call(0, 1, MOD_CONTROL|MOD_SHIFT, VK_V)
	if ret == 0 {
		return err
	}

	pkg.HandleContext(ctx) // Post quit message to current thread on context cancellation

	var msg [56]byte // MSG struct is 56 bytes on 64-bit Windows
	for {
		r, _, _ := getMessage.Call(uintptr(unsafe.Pointer(&msg[0])), 0, 0, 0)
		if r == WM_QUIT {
			break
		}
		msgType := *(*uint32)(unsafe.Pointer(&msg[8]))
		if msgType == WM_HOTKEY {
			pasteHander.HandlePaste()
		}
	}
	return nil
}
