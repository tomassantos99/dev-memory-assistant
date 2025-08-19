package ui

import (
	"fmt"
	// "os"
	// "path/filepath"
	"strings"
	"syscall"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/tomassantos99/dev-memory-assistant/paste/pkg"
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
	textSearch      *walk.LineEdit
	onItemSelection func(selectedItem string) error
}

type ClipboardModel struct {
	walk.ListModelBase
	items        []string
	displayItems []string
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

		bounds := st.Bounds()
		bounds.X += 8
		bounds.Width -= 8

		st.DrawText(text, bounds, walk.TextLeft|walk.TextVCenter|walk.TextSingleLine|walk.TextEndEllipsis)
	}
}

func NewEnvModel() *ClipboardModel {
	return &ClipboardModel{
		items:        []string{},
		displayItems: []string{},
	}
}

func (m *ClipboardModel) ItemCount() int {
	return len(m.displayItems)
}

func (m *ClipboardModel) Value(index int) any {
	return pkg.CropString(m.displayItems[index], 80)
}

func (m *ClipboardModel) SetItems(items []string) {
	m.displayItems = items
	m.PublishItemsReset()
}

func (m *ClipboardModel) SetOriginalItems(items []string) {
	m.items = items
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

	var icon, iconErr = walk.NewIconFromResource("CLIPBOARD_ICON")
	if iconErr != nil {
		panic(iconErr)
	}

	var err = MainWindow{
		AssignTo: &window.Mw,
		Title:    "Clipboard History",
		Size:     Size{Width: 1000, Height: 1000},
		Layout:   VBox{},
		Visible:  false,
		Icon: icon,
		Children: []Widget{
			Label{
				Text: "Search:",
			},
			LineEdit{
				AssignTo:      &window.textSearch,
				Font:          Font{Family: "Segoe UI", PointSize: 11},
				StretchFactor: 0,
				OnTextChanged: window.onSearchInput,
				OnKeyDown:     window.onSearchKeyDown,
			},
			HSplitter{
				StretchFactor: 1,
				Children: []Widget{
					ListBox{
						AssignTo:              &window.lb,
						Model:                 window.model,
						OnCurrentIndexChanged: window.onCurrentIndexChanged,
						OnItemActivated:       window.onItemActivated,
						OnKeyDown:             window.onLbKeyDown,
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

		w.model.SetOriginalItems(items) // originalItems
		w.model.SetItems(items)         // displayItems

		if btpErr := w.Mw.BringToTop(); btpErr != nil {
			panic(btpErr)
		}

		procSetForegroundWindow.Call(uintptr(w.Mw.Handle()))

		w.lb.SetCurrentIndex(0)

		if err := w.lb.SetFocus(); err != nil {
			panic(err)
		}

		w.Mw.Show()
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

func (w *HistoryWindow) onSearchInput() {
	var filteredItems []string

	for _, item := range w.model.items {
		if strings.Contains(item, w.textSearch.Text()) {
			filteredItems = append(filteredItems, item)
		}
	}

	w.model.SetItems(filteredItems)
	w.model.PublishItemsReset()
}

func (w *HistoryWindow) onSearchKeyDown(key walk.Key) {
	if key == walk.KeyDown {
		w.lb.SetFocus()
		w.lb.SetCurrentIndex(0)
	}
}

func (w *HistoryWindow) onLbKeyDown(key walk.Key) {
	if (walk.ModifiersDown() == walk.ModControl && key == walk.KeyS) || (key == walk.KeyUp && w.lb.CurrentIndex() == 0) {
		w.textSearch.SetFocus()
	}
}
