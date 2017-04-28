package main

import (
	"fmt"
	"os"

	"github.com/icza/gowut/gwu"
	//"github.com/willauld/gowut/gwu"
)

func buildTabPanel(j *jawaInfo) gwu.Comp {

	t := gwu.NewTabPanel()
	t.Style().SetSizePx(800, 400) // was 500 x 300

	t.SetTabBarPlacement(gwu.TbPlacementLeft)
	t.TabBarFmt().SetHAlign(gwu.HALeft)
	t.TabBarFmt().SetVAlign(gwu.VATop)

	t.AddString("Record Payments", buildRecordPayments( /*t,*/ j))
	// &&&& TODO end first panal
	c := gwu.NewPanel()
	c.Add(gwu.NewLabel("You have no new messages."))
	t.AddString("Reminders", c)
	// &&& TODO end second panal
	c = gwu.NewPanel()
	c.Add(gwu.NewLabel("You have no sent messages."))
	t.AddString("Status", c)
	// &&& TODO end panal
	t.AddString("Edit Rentals", buildEditRentals(j))
	// &&& TODO end panal
	c = gwu.NewPanel()
	tb := gwu.NewTextBox("Click to edit this comment.")
	tb.SetRows(10)
	tb.SetCols(40)
	c.Add(tb)
	t.AddString("Comment", c)

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
