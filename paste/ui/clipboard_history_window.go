package ui

import (
	"fmt"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type HistoryWindow struct {
	Mw *walk.MainWindow
	Lb *walk.ListBox
}

func CreateHistoryWindow() *HistoryWindow {
	var window = &HistoryWindow{}

	var err = MainWindow{
		AssignTo: &window.Mw,
		Title:    "Clipboard History",
		Size:     Size{Width: 400, Height: 400},
		Layout:   VBox{},
		Visible:  false,
		Children: []Widget{
			ListBox{
				AssignTo: &window.Lb,
				Model:    []string{},
			},
		},
	}.Create()

	if err != nil {
		fmt.Println("Error creating history window:")
		panic(err)
	}

	window.Mw.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		*canceled = true // prevent destruction on user close
		window.Mw.SetVisible(false)
	})

	return window
}

func (w *HistoryWindow) ShowHistoryWindow(items []string) {
	if w == nil || w.Mw == nil {
		return
	}

	w.Mw.Synchronize(func() {

		var modelErr = w.Lb.SetModel(items)
		if modelErr != nil {
			panic(modelErr)
		}

		w.Mw.Show()

		var btpErr = w.Mw.BringToTop()
		if btpErr != nil {
			panic(btpErr)
		}
	})

}
