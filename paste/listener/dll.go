package listener

import "syscall"

var (
	user32                     = syscall.MustLoadDLL("user32.dll")
	openClipboard              = user32.MustFindProc("OpenClipboard")
	getClipboardData           = user32.MustFindProc("GetClipboardData")
	closeClipboard             = user32.MustFindProc("CloseClipboard")
	defProc                    = user32.MustFindProc("DefWindowProcW")
	registerClassEx            = user32.MustFindProc("RegisterClassExW")
	createWindowEx             = user32.MustFindProc("CreateWindowExW")
	addClipboardFormatListener = user32.MustFindProc("AddClipboardFormatListener")
	getMessage                 = user32.MustFindProc("GetMessageW")
	dispatchMessage            = user32.MustFindProc("DispatchMessageW")
	registerHotKey             = user32.MustFindProc("RegisterHotKey")
)
