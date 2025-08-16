package handler

import (
	"runtime"
	"time"
	"github.com/atotto/clipboard"
	"github.com/micmonay/keybd_event"
)

func PasteContent(content string) error {
    oldClip, err := clipboard.ReadAll()
    if err != nil {
        return err
    }

	if err := clipboard.WriteAll(content); err != nil {
		return err
	}

	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		return err
	}

	kb.SetKeys(keybd_event.VK_V)

	switch runtime.GOOS {
	case "darwin":
		kb.HasSuper(true)
	default:
		kb.HasCTRL(true)
	}

	kb.Launching()

    time.Sleep(50 *time.Millisecond)

    if err := clipboard.WriteAll(oldClip); err != nil {
		return err
	}

	return nil
}
