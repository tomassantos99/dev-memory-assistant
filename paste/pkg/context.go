package pkg

import (
	"context"
	"fmt"
	"syscall"
)

const (
	WM_QUIT = 0x0012
)

var (
	user32            = syscall.MustLoadDLL("user32.dll")
	postThreadMessage = user32.MustFindProc("PostThreadMessageW")

	kernel32           = syscall.MustLoadDLL("kernel32.dll")
	getCurrentThreadId = kernel32.MustFindProc("GetCurrentThreadId")
)

func HandleContext(ctx context.Context) {
	threadID, _, _ := getCurrentThreadId.Call()
	go func(threadId uintptr) {
		<-ctx.Done()
		postThreadMessage.Call(threadId, uintptr(WM_QUIT), 0, 0)
		fmt.Println("Posted WM_QUIT to thread:", threadId)
	}(threadID)
}
