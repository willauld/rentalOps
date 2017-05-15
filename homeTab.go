package main

import "github.com/icza/gowut/gwu"

func buildHomeTab(j *jawaInfo) (gwu.Panel, gwu.TextBox) {

	c := gwu.NewPanel()
	stb := gwu.NewTextBox("")
	stb.Style().SetWidthPx(1).SetHeightPx(1)
	stb.AddEHandlerFunc(func(e gwu.Event) {
		Notify("Focus is on Home Tab", e)
	}, gwu.ETypeFocus)
	c.Add(stb)

	c.Add(gwu.NewLabel("Yuli, What do you want to see on this page?"))

	return c, stb
}
