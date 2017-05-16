package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/icza/gowut/gwu"
	//"github.com/willauld/gowut/gwu"
)

var nl gwu.Label

// Notify posts a notification above Rental Ops Table
func Notify(n string, e gwu.Event) {
	nl.SetText(n)
	if e != nil {
		e.MarkDirty(nl)
	}
}

var tabUpdateFuncs [20]func()

func registerTabUpdateFunc(idx int, f func()) {
	if idx < 0 || idx >= 20 {
		return
	}
	tabUpdateFuncs[idx] = f
}

func updateTab(idx int) {
	if tabUpdateFuncs[idx] != nil {
		tabUpdateFuncs[idx]()
	}
}

func getKeyList(o /*map[string]*/ interface{}) []string {
	kvm, ok := o.(map[string]interface{})
	if !ok {
		fmt.Printf("HELP\n")
	}
	if len(kvm) < 1 {
		fmt.Printf("There are no keys yet defined\n")
		keys := make([]string, 1)
		keys[0] = "undefined"
		return keys
	}
	keys := make([]string, len(kvm))

	i := 0
	for k := range kvm {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

// UpdateListBox updates a gwu.ListBox item with a new list and handler
func UpdateListBox(lb *gwu.ListBox, name *string, l []string,
	lbHandler func(gwu.Event)) error {

	if len(l) < 1 {
		return fmt.Errorf("UpdateListBox() list is empty")
	}
	(*lb).ClearSelected()
	(*lb) = gwu.NewListBox(l)
	(*lb).Style().SetFullWidth()
	nameIndex := strIndex(l, *name)
	if nameIndex == -1 {
		// try to do next best thing
		nameIndex = 0
		*name = l[0]
		//return nil, fmt.Errorf("UpdateListBox() [%s] is not in %+v", name, l)
	}
	(*lb).SetSelected(nameIndex, true)
	(*lb).AddEHandlerFunc(lbHandler, gwu.ETypeChange)

	return nil
}

func buildPanelTab(s string, t *gwu.TabPanel,
	f func(*jawaInfo) (gwu.Panel, gwu.TextBox), j *jawaInfo) {
	erl := gwu.NewLabel(s)
	erl.Style().SetDisplay(gwu.DisplayBlock) // Display: block - so the whole cell of the tab is clickable
	tabc, objToFocus := f(j)
	erl.AddEHandlerFunc(func(e gwu.Event) {
		e.SetFocusedComp(objToFocus)
		//Notify("Executing tabc handler", e)
	}, gwu.ETypeClick)
	(*t).Add(erl, tabc)
}

func buildTabPanel(j *jawaInfo) gwu.Comp {

	t := gwu.NewTabPanel()
	t.Style().SetSizePx(800, 400)
	t.SetTabBarPlacement(gwu.TbPlacementLeft)
	t.TabBarFmt().SetHAlign(gwu.HALeft)
	t.TabBarFmt().SetVAlign(gwu.VATop)

	buildPanelTab("Home Tab", &t, buildHomeTab, j)
	buildPanelTab("Record Payments", &t, buildRecordPayments, j)
	buildPanelTab("Reminders", &t, buildRemindPage, j)
	buildPanelTab("Edit Rentals", &t, buildEditRentals, j)
	buildPanelTab("Edit Tenants", &t, buildEditTenants, j)
	buildPanelTab("Edit Common Info", &t, buildEditCommon, j)
	buildPanelTab("Dump DB", &t, buildDisplayDB, j)
	return t
}

func establishWindow(j *jawaInfo) {
	var server gwu.Server
	addr := "localhost:48991"
	// Create and build a window
	win := gwu.NewWindow("main", "Rental Ops tool")
	win.Style().SetFullWidth()
	win.SetHAlign(gwu.HACenter)
	win.SetCellPadding(2)

	p := gwu.NewPanel()
	hp := gwu.NewHorizontalPanel()
	hp.Style().SetFullWidth()
	hp.SetHAlign(gwu.HARight)

	bs := gwu.NewButton("SAVE current state")
	bs.AddEHandlerFunc(func(e gwu.Event) {
		err := Save(dataTarget, j)
		if err != nil {
			fmt.Printf("Save Failed: %v\n", err) //TODO: replace with notification
			return
		}
		fmt.Printf("Data Saved\n") //TODO: replace with notification
	}, gwu.ETypeClick)
	hp.Add(bs)
	fbs := gwu.NewButton("SAVE and Exit the Application")
	fbs.AddEHandlerFunc(func(e gwu.Event) {
		err := Save(dataTarget, j)
		if err != nil {
			fmt.Printf("Save Failed - not exiting: %v\n", err) // TODO: replace with notification
			return
		}
		fmt.Printf("Data Saved\n") // TODO: replace with notification
		os.Exit(2)
	}, gwu.ETypeClick)
	hp.Add(fbs)
	fb := gwu.NewButton("NO SAVE and Exit the Application")
	fb.AddEHandlerFunc(func(e gwu.Event) {
		//server.Close()
		os.Exit(2)
	}, gwu.ETypeClick)
	hp.Add(fb)
	p.Add(hp)

	np := gwu.NewHorizontalPanel() //notification pannel
	np.Style().SetFullWidth()
	np.SetHAlign(gwu.HACenter)
	nl = gwu.NewLabel("")
	np.Add(nl)

	p.Add(np)

	p.AddVSpace(15)

	t := buildTabPanel(j)
	p.Add(t)
	win.Add(p)

	// Create and start a GUI server (omitting error check)
	server = gwu.NewServer("testgui", addr /*"localhost:48991"*/)
	server.SetText("Rental Ops WIP")
	server.AddWin(win)
	server.Start("main") // Also opens windows list in browser if param is ""
}
