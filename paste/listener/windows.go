package listener

import (
	"errors"
	"syscall"
	"unsafe"
	"github.com/tomassantos99/dev-memory-assistant/paste/handler"
)

// I have no idea wtf is going on here, but this is the code I found with the help of my two best friends, Google and ChatGPT.
// Not sure how it works but it does. It basicly sets up a listener for clipboard updates on Windows and execute a callback when the clipboard is updated.
// I should probably read on how this works, but for now it is too low-level and Windows specific. Not one of my priorities, don't judge me.

const (
	WM_CLIPBOARDUPDATE = 0x031D
)

var globalClipboardHandler *handler.ClipboardHandler

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
	user32 := syscall.MustLoadDLL("user32.dll")
	if msg == WM_CLIPBOARDUPDATE {
		openClipboard := user32.MustFindProc("OpenClipboard")
		getClipboardData := user32.MustFindProc("GetClipboardData")
		closeClipboard := user32.MustFindProc("CloseClipboard")

		const CF_UNICODETEXT = 13

		openClipboard.Call(0)
		h, _, _ := getClipboardData.Call(CF_UNICODETEXT)
		if h != 0 {
			ptr := uintptr(h)
			text := syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(ptr))[:])
			globalClipboardHandler.ClipboardMessages <- text
		}
		closeClipboard.Call()
	}
	// Call default window procedure

	defProc := user32.MustFindProc("DefWindowProcW")
	ret, _, _ := defProc.Call(hwnd, uintptr(msg), wParam, lParam)
	return ret
}

func ListenWindowsClipboardUpdates(clipboardHandler *handler.ClipboardHandler) {
	globalClipboardHandler = clipboardHandler

	user32 := syscall.MustLoadDLL("user32.dll")
	defer user32.Release()

	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	defer kernel32.Release()

	registerClassEx := user32.MustFindProc("RegisterClassExW")
	createWindowEx := user32.MustFindProc("CreateWindowExW")
	addClipboardFormatListener := user32.MustFindProc("AddClipboardFormatListener")
	getMessage := user32.MustFindProc("GetMessageW")
	dispatchMessage := user32.MustFindProc("DispatchMessageW")

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
	for {
		ret, _, _ := getMessage.Call(uintptr(unsafe.Pointer(&msg[0])), hwnd, 0, 0)
		if ret == 0 {
			break // WM_QUIT
		}
		dispatchMessage.Call(uintptr(unsafe.Pointer(&msg[0])))
	}
}
