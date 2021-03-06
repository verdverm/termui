// Panels Layout is a Flex widget with
// hidable panels and a main content.
// Can be vert or horz oriendted throught the Flex widget.
// Panels can be jumped to with <focuskey> and hidden with <shift>-<focuskey>
// Recommend making the <focuskey> an: '<alt>-<key>' and hidden will be '<shift>-<alt>-<key>'
// Normal movement and interaction keys within the focussed panel.
//
// main (middle) panel, can be anything, including...
// - the router (when this is the root view)
// - another DashAndPanels
// - a pager, grid, or any other primitive
package panels

import (
	"fmt"
	"sort"

	"github.com/verdverm/tview"
	"github.com/verdverm/vermui"
	"github.com/verdverm/vermui/events"
)

type Panel struct {
	Name       string
	Item       tview.Primitive
	FixedSize  int
	Proportion int
	Focus      int
	FocusKey   string
	Hidden     bool
	HiddenKey  string
}

type Layout struct {
	*tview.Flex

	// first (left/top) panels, can be almost anything and hidden.
	fPanels map[string]*Panel

	// main (middle) panel, can be anything, I think.
	mPanel *Panel

	// last (right/bottom) panels, can be almost anything and hidden.
	lPanels map[string]*Panel
}

func New() *Layout {
	L := &Layout{
		Flex:    tview.NewFlex(),
		fPanels: map[string]*Panel{},
		lPanels: map[string]*Panel{},
	}

	return L
}

// AddFirstPanel adds a Panel to the left or top, depending on orientation.
func (L *Layout) AddFirstPanel(name string, item tview.Primitive, fixedSize, proportion,
	focus int, focuskey string, hidden bool, hiddenkey string) {
	panel := &Panel{
		Name:       name,
		Item:       item,
		FixedSize:  fixedSize,
		Proportion: proportion,
		Focus:      focus,
		FocusKey:   focuskey,
		Hidden:     hidden,
		HiddenKey:  hiddenkey,
	}

	L.fPanels[name] = panel
}

// AddLastPanel adds a Panel to the right or bottom, depending on orientation.
func (L *Layout) AddLastPanel(name string, item tview.Primitive, fixedSize, proportion,
	focus int, focuskey string, hidden bool, hiddenkey string) {
	panel := &Panel{
		Name:       name,
		Item:       item,
		FixedSize:  fixedSize,
		Proportion: proportion,
		Focus:      focus,
		FocusKey:   focuskey,
		Hidden:     hidden,
		HiddenKey:  hiddenkey,
	}

	L.lPanels[name] = panel
}

func (L *Layout) SetMainPanel(name string, item tview.Primitive, fixedSize, proportion, focus int, focuskey string) {
	panel := &Panel{
		Name:       name,
		Item:       item,
		FixedSize:  fixedSize,
		Proportion: proportion,
		Focus:      focus,
		FocusKey:   focuskey,
	}

	L.mPanel = panel
}

func (L *Layout) Mount(context map[string]interface{}) error {
	err := L.build()
	if err != nil {
		return err
	}

	// Setup focuskeys
	for _, panel := range L.fPanels {
		panel.Item.Mount(context)
		if panel.FocusKey != "" {
			localPanel := panel
			vermui.AddWidgetHandler(L, "/sys/key/"+localPanel.FocusKey, func(e events.Event) {
				go events.SendCustomEvent("/console/trace", "Focus: "+localPanel.Name)
				vermui.SetFocus(localPanel.Item)
			})
		}
		if panel.HiddenKey != "" {
			localPanel := panel
			vermui.AddWidgetHandler(L, "/sys/key/"+localPanel.HiddenKey, func(e events.Event) {
				localPanel.Hidden = !localPanel.Hidden
				go events.SendCustomEvent("/console/trace", fmt.Sprintf("Hidden: %s (%v)", localPanel.Name, localPanel.Hidden))
				L.build()
				if localPanel.Hidden {
					vermui.SetFocus(L.mPanel.Item)
				} else {
					vermui.SetFocus(localPanel.Item)
				}
				vermui.Draw()
			})
		}
	}
	if L.mPanel.FocusKey != "" {
		L.mPanel.Item.Mount(context)
		localPanel := L.mPanel
		vermui.AddWidgetHandler(L, "/sys/key/"+localPanel.FocusKey, func(e events.Event) {
			go events.SendCustomEvent("/console/trace", "Focus: "+localPanel.Name)
			vermui.SetFocus(localPanel.Item)
		})
	}
	for _, panel := range L.lPanels {
		panel.Item.Mount(context)
		if panel.FocusKey != "" {
			localPanel := panel
			vermui.AddWidgetHandler(L, "/sys/key/"+localPanel.FocusKey, func(e events.Event) {
				go events.SendCustomEvent("/console/trace", "Focus: "+localPanel.Name)
				vermui.SetFocus(localPanel.Item)
			})
		}
		if panel.HiddenKey != "" {
			localPanel := panel
			vermui.AddWidgetHandler(L, "/sys/key/"+localPanel.HiddenKey, func(e events.Event) {
				localPanel.Hidden = !localPanel.Hidden
				go events.SendCustomEvent("/console/trace", fmt.Sprintf("Hidden: %s (%v)", localPanel.Name, localPanel.Hidden))
				L.build()
				if localPanel.Hidden {
					vermui.SetFocus(L.mPanel.Item)
				} else {
					vermui.SetFocus(localPanel.Item)
				}
				vermui.Draw()
			})
		}
	}

	return nil
}

func (L *Layout) build() error {
	// get and order the fPanels
	fPs := []*Panel{}
	for _, panel := range L.fPanels {
		if panel.Hidden {
			continue
		}
		fPs = append(fPs, panel)
	}
	sort.Slice(fPs, func(i, j int) bool {
		return fPs[i].Focus < fPs[j].Focus
	})

	// get and order the lPanels
	lPs := []*Panel{}
	for _, panel := range L.lPanels {
		if panel.Hidden {
			continue
		}
		lPs = append(lPs, panel)
	}
	sort.Slice(lPs, func(i, j int) bool {
		return lPs[i].Focus < lPs[j].Focus
	})

	// Start a fresh Flex item
	orient := L.GetDirection()
	L.Flex = tview.NewFlex().SetDirection(orient)

	for _, p := range fPs {
		L.AddItem(p.Item, p.FixedSize, p.Proportion, false)
	}

	p := L.mPanel
	L.AddItem(p.Item, p.FixedSize, p.Proportion, true)

	for _, p := range lPs {
		L.AddItem(p.Item, p.FixedSize, p.Proportion, false)
	}

	return nil
}
