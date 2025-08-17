package ui

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/tomassantos99/dev-memory-assistant/paste/pkg"
	"syscall"
)

var (
	user32                  = syscall.MustLoadDLL("user32.dll")
	procSetForegroundWindow = user32.MustFindProc("SetForegroundWindow")
)

type HistoryWindow struct {
	Mw              *walk.MainWindow
	lb              *walk.ListBox
	te              *walk.TextEdit
	model           *ClipboardModel
	onItemSelection func(selectedItem string) error
}

type ClipboardModel struct {
	walk.ListModelBase
	items []string
}

type uniformStyler struct {
	h     int
	font  *walk.Font
	model walk.ListModel
}

func (s *uniformStyler) ItemHeightDependsOnWidth() bool  { return false }
func (s *uniformStyler) DefaultItemHeight() int          { return s.h }
func (s *uniformStyler) ItemHeight(index, width int) int { return s.h }
func (s *uniformStyler) StyleItem(st *walk.ListItemStyle) {
	st.DrawBackground()

	if s.model != nil {
		text := fmt.Sprint(s.model.Value(st.Index()))
		st.DrawText(text, st.Bounds(), walk.TextLeft|walk.TextVCenter|walk.TextSingleLine|walk.TextEndEllipsis)
	}
}

func NewEnvModel() *ClipboardModel {
	return &ClipboardModel{items: []string{}}
}

func (m *ClipboardModel) ItemCount() int {
	return len(m.items)
}

func (m *ClipboardModel) Value(index int) any {
	return pkg.CropString(m.items[index], 80)
}

func (m *ClipboardModel) SetItems(items []string) {
	m.items = items
	m.PublishItemsReset()
}

func CreateHistoryWindow(onItemSelection func(selectedItem string) error) *HistoryWindow {

	var window = &HistoryWindow{
		model:           NewEnvModel(),
		onItemSelection: onItemSelection,
	}

	var font, fontErr = walk.NewFont("Segoe UI", 11, 0)
	if fontErr != nil {
		panic(fontErr)
	}

	var err = MainWindow{
		AssignTo: &window.Mw,
		Title:    "Clipboard History",
		Size:     Size{Width: 1000, Height: 1000},
		Layout:   VBox{},
		Visible:  false,
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					ListBox{
						AssignTo:              &window.lb,
						Model:                 window.model,
						OnCurrentIndexChanged: window.onCurrentIndexChanged,
						OnItemActivated:       window.onItemActivated,
						ItemStyler:            &uniformStyler{h: 28, font: font, model: window.model},
					},
					TextEdit{
						AssignTo: &window.te,
						ReadOnly: true,
					},
				},
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

	action := walk.NewAction()
	action.SetShortcut(walk.Shortcut{Key: walk.KeyEscape})
	action.Triggered().Attach(func() {
		window.Mw.SetVisible(false)
	})
	window.Mw.ShortcutActions().Add(action)

	return window
}

func (w *HistoryWindow) ShowHistoryWindow(items []string) {
	if w == nil || w.Mw == nil || w.IsVisible() {
		return
	}

	w.Mw.Synchronize(func() {

		w.model.SetItems(items)

		w.lb.SetCurrentIndex(0)

		w.Mw.Show()

		if btpErr := w.Mw.BringToTop(); btpErr != nil {
			panic(btpErr)
		}

		procSetForegroundWindow.Call(uintptr(w.Mw.Handle()))
	})
}

func (w *HistoryWindow) IsVisible() bool {
	return w.Mw != nil && w.Mw.Visible()
}

func (w *HistoryWindow) onItemActivated() {
	var index = w.lb.CurrentIndex()
	if index < 0 {
		return
	}

	var item = w.model.items[index]
	w.Mw.SetVisible(false)
	err := w.onItemSelection(item)
	if err != nil {
		fmt.Println("Error executing clipboard selection callback:", err)
	}
}

func (w *HistoryWindow) onCurrentIndexChanged() {
	var index = w.lb.CurrentIndex()
	if index < 0 {
		return
	}

	var item = w.model.items[index]
	w.te.SetText(item)
}
