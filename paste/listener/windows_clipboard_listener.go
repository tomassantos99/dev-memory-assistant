package listener

import (
	"context"
	"errors"
	"github.com/tomassantos99/dev-memory-assistant/paste/pkg"
	"syscall"
	"unsafe"
)

// I have no idea wtf is going on here, but this is the code I found with the help of my two best friends, Google and ChatGPT.
// Not sure how it works but it does. It basicly sets up a listener for clipboard updates on Windows and execute a callback when the clipboard is updated.
// I should probably read on how this works, but for now it is too low-level and Windows specific. Not one of my priorities, don't judge me.

const (
	WM_CLIPBOARDUPDATE = 0x031D
)

var gobalClipboardHandlerChannel chan string

type WNDCLASSEX struct {
	cbSize        uint32
	style         uint32
	lpfnWndProc   uintptr
	cbClsExtra    int32
	cbWndExtra    int32
	hInstance     uintptr
	hIcon         uintptr
	hCursor       uintptr
	hbrBackground uintptr
	lpszMenuName  *uint16
	lpszClassName *uint16
	hIconSm       uintptr
}

func onSysMessage(hwnd uintptr, msg uint32, wParam, lParam uintptr) uintptr {
	if msg == WM_CLIPBOARDUPDATE {
		const CF_UNICODETEXT = 13

		openClipboard.Call(0)
		h, _, _ := getClipboardData.Call(CF_UNICODETEXT)
		if h != 0 {
			ptr := uintptr(h)
			text := syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(ptr))[:])
			gobalClipboardHandlerChannel <- text
		}
		closeClipboard.Call()
	}
	// Call default window procedure

	ret, _, _ := defProc.Call(hwnd, uintptr(msg), wParam, lParam)
	return ret
}

func ListenWindowsClipboardUpdates(ctx context.Context, clipboardHandlerChannel chan string) {
	gobalClipboardHandlerChannel = clipboardHandlerChannel

	className, err := syscall.UTF16PtrFromString("WindowsClipboardListener")
	if err != nil {
		panic(errors.New("failed to convert class name to UTF16"))
	}

	// Register window class
	wndClass := WNDCLASSEX{
		cbSize:        uint32(unsafe.Sizeof(WNDCLASSEX{})),
		style:         0,
		lpfnWndProc:   syscall.NewCallback(onSysMessage),
		cbClsExtra:    0,
		cbWndExtra:    0,
		hInstance:     0,
		hIcon:         0,
		hCursor:       0,
		hbrBackground: 0,
		lpszMenuName:  nil,
		lpszClassName: className,
		hIconSm:       0,
	}
	registerClassEx.Call(uintptr(unsafe.Pointer(&wndClass)))

	// Create message-only window
	HWND_MESSAGE := ^uintptr(2) // (HWND_MESSAGE = (HWND)-3)
	hwnd, _, _ := createWindowEx.Call(
		0, uintptr(unsafe.Pointer(className)), 0, 0, 0, 0, 0, 0,
		HWND_MESSAGE, 0, 0, 0,
	)

	// Register for clipboard updates
	addClipboardFormatListener.Call(hwnd)

	var msg [56]byte // MSG struct is 56 bytes on 64-bit Windows

	pkg.HandleContext(ctx) // Post quit message to current thread on context cancellation
	for {
		ret, _, _ := getMessage.Call(uintptr(unsafe.Pointer(&msg[0])), hwnd, 0, 0)
		if ret == 0 {
			return // WM_QUIT
		}
		dispatchMessage.Call(uintptr(unsafe.Pointer(&msg[0])))
	}
}
