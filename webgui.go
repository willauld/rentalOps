package main

import (
	"fmt"
	"os"

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

func buildTabPanel(j *jawaInfo) gwu.Comp {

	t := gwu.NewTabPanel()
	t.Style().SetSizePx(800, 400)
	/*
		t.AddEHandlerFunc(func(e gwu.Event) {
			// how to go from t.Selected() to something I can put
			// in
			Notify(fmt.Sprintf("Clicked on tab: %d", t.Selected()), e)
			fmt.Printf("In tabPanel state change handler\n")
			t.CompAt(t.Selected())
		}, gwu.ETypeStateChange)
	*/
	t.SetTabBarPlacement(gwu.TbPlacementLeft)
	t.TabBarFmt().SetHAlign(gwu.HALeft)
	t.TabBarFmt().SetVAlign(gwu.VATop)

	t.AddString("Record Payments", buildRecordPayments(j))
	// &&&& TODO end first panal
	c := buildRemindPage(j)
	t.AddString("Reminders", c)
	// &&& TODO end second panal
	c = gwu.NewPanel()
	c.Add(gwu.NewLabel("You have no sent messages."))
	t.AddString("Status", c)
	// &&& TODO end panal
	erl := gwu.NewLabel("Edit Rentals")
	erl.Style().SetDisplay(gwu.DisplayBlock) // Display: block - so the whole cell of the tab is clickable
	tabc, objToFocus := buildEditRentals(j)
	erl.AddEHandlerFunc(func(e gwu.Event) {
		e.SetFocusedComp(objToFocus)
		fmt.Printf("in tabc lable click handler lunch point\n")
		Notify("Executing tabc handler", e)
	}, gwu.ETypeClick)
	//t.AddString("Edit Rentals", buildEditRentals(j))
	t.Add(erl, tabc)
	// &&& TODO end panal
	t.AddString("Edit Tenants", buildEditTenants(j))
	// &&& TODO end panal
	t.AddString("Edit Common", buildEditCommon(j))
	// &&& TODO end panal
	t.AddString("Display DB", buildDisplayDB(j))
	// &&& TODO end panal

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
