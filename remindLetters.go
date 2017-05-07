package main

import (
	"fmt"
	"strings"

	"github.com/icza/gowut/gwu"
)

func sprintIfGraterThanZero(amt float32, name string,
	str string, rows, cols int) (string, int, int) {
	if amt > 0 {
		str1 := fmt.Sprintf("\t\t%-6.2f\t%s\n", amt, name)
		rows++
		if len(str) > cols {
			cols = len(str)
		}
		str += str1
	}
	return str, rows, cols
}
func getAmountDetails(j *jawaInfo, apt string) (rows, cols int, str string) {
	r := j.Rental[apt]
	tr := j.Tenant[r.TenantKey]
	procRec := []struct {
		owed float32
		s    string
	}{
		//{r.Rent, "Next Rent"},
		{tr.RentOwed, "Rent"},
		{tr.LatePOwed, "Late Fee"},
		{tr.BounceOwed, "Bounce Fee"},
		{tr.WaterOwed, "Water Fee"},
	}
	sum := tr.RentOwed + tr.LatePOwed + tr.BounceOwed + tr.WaterOwed
	for _, v := range procRec {
		str, rows, cols =
			sprintIfGraterThanZero(v.owed, v.s, str, rows, cols)
	}
	str += "\n\t\t========\n"
	str += fmt.Sprintf("\t\t%-6.2f\t%s\n", sum, "Balance due")
	//fmt.Printf("\n")
	if cols < 60 {
		cols = 60
	}
	rows += 4
	return rows, cols, str
}
func updateRemindPage(j *jawaInfo, apt string, e gwu.Event,
	OStreetL, OCityStateZip, TenName, TStreet,
	TCityStateZip, LDate, DearT, FirstP gwu.Label, Details gwu.TextBox,
	SecondP, Salutation, Correspondent gwu.Label) {

	arec := j.Rental[apt]
	//trec := j.Tenant
	ci := j.CI
	OStreetL.SetText(ci.OfficeStreet)
	str := fmt.Sprintf("%s, %s %s", ci.OfficeCity, ci.OfficeState, ci.OfficeZip)
	OCityStateZip.SetText(str)
	TenName.SetText(arec.Tenant)
	TStreet.SetText(arec.Street)
	str = fmt.Sprintf("%s, %s %s", arec.City, arec.State, arec.Zip)
	TCityStateZip.SetText(str)
	t := timeNowRental()
	str = fmt.Sprintf("%s %d, %d", t.Month(), t.Day(), t.Year())
	LDate.SetText(str)
	str = fmt.Sprintf("Dear %s,\n", strings.Fields(arec.Tenant)[0])
	DearT.SetText(str)
	FirstP.SetText("You currently owe, based on the following details:")
	rows, cols, dtext := getAmountDetails(j, apt)
	Details.SetRows(rows)
	Details.SetCols(cols)
	Details.SetText(dtext)
	SecondP.SetText("Thank you for your prompt payment.")
	Salutation.SetText("Thank you,")
	Correspondent.SetText(ci.CorrespondenceTitle)

	if e != nil {
		e.MarkDirty(OStreetL)
		e.MarkDirty(OCityStateZip)
		e.MarkDirty(TenName)
		e.MarkDirty(TStreet)
		e.MarkDirty(TCityStateZip)
		e.MarkDirty(LDate)
		e.MarkDirty(DearT)
		e.MarkDirty(FirstP)
		e.MarkDirty(Details)
		e.MarkDirty(SecondP)
		e.MarkDirty(Salutation)
		e.MarkDirty(Correspondent)
	}
}

func buildRemindPage(j *jawaInfo) gwu.Panel {
	var OStreetL, OCityStateZip, TenName, TStreet gwu.Label
	var TCityStateZip, LDate, DearT, FirstP gwu.Label
	var Details gwu.TextBox
	var SecondP, Salutation, Correspondent gwu.Label
	var aptlb gwu.ListBox
	var apt string

	c := gwu.NewPanel()
	c.AddEHandlerFunc(func(e gwu.Event) {
		list := getAptList(j)
		apt = list[aptlb.SelectedIdx()]

		updateRemindPage(j, apt, e,
			OStreetL, OCityStateZip, TenName, TStreet,
			TCityStateZip, LDate, DearT, FirstP, Details,
			SecondP, Salutation, Correspondent)

		e.MarkDirty(aptlb)
		Notify("Remind Leter Focus Event happened", e)
		fmt.Printf("inside the remind page state change handler\n")
	}, gwu.ETypeStateChange /*gwu.ETypeFocus*/ /*gwu.ETypeClick*/)

	tableA := gwu.NewTable()
	tableA.SetCellPadding(2)
	tableA.EnsureSize(2, 5)
	tableA.Add(gwu.NewLabel("Letter for Rental:"), 0, 0)
	list := getAptList(j)
	apt = list[0]
	aptlb = gwu.NewListBox(list)
	aptlb.Style().SetFullWidth()
	aptlbHandler := func(e gwu.Event) {
		list := getAptList(j)
		apt = list[aptlb.SelectedIdx()]

		updateRemindPage(j, apt, e,
			OStreetL, OCityStateZip, TenName, TStreet,
			TCityStateZip, LDate, DearT, FirstP, Details,
			SecondP, Salutation, Correspondent)

		e.MarkDirty(aptlb)
	}
	aptlb.AddEHandlerFunc(aptlbHandler, gwu.ETypeChange)

	tableA.Add(aptlb, 0, 1)
	tableA.Add(gwu.NewLabel("...."), 0, 3)

	c.Add(tableA)

	c.AddVSpace(15)

	//hp := gwu.NewHorizontalPanel()

	table := gwu.NewTable()
	table.SetCellPadding(2)
	table.EnsureSize(19, 4)

	OStreetL = gwu.NewLabel("")
	OCityStateZip = gwu.NewLabel("")

	TenName = gwu.NewLabel("")
	TStreet = gwu.NewLabel("")
	TCityStateZip = gwu.NewLabel("")

	LDate = gwu.NewLabel("")
	DearT = gwu.NewLabel("")

	FirstP = gwu.NewLabel("")

	hp := gwu.NewHorizontalPanel()
	hp.Add(NewHSpace(45))
	Details = gwu.NewTextBox("")
	hp.Add(Details)

	SecondP = gwu.NewLabel("")

	Salutation = gwu.NewLabel("")
	Correspondent = gwu.NewLabel("")

	//AptName.Style().SetWidthPx(260) //??????? WGA
	table.Add(OStreetL, 0, 1)
	table.Add(OCityStateZip, 1, 1)
	table.Add(NewVSpace(55), 2, 1)
	table.Add(TenName, 3, 1)
	table.Add(TStreet, 4, 1)
	table.Add(TCityStateZip, 5, 1)
	table.Add(NewVSpace(15), 6, 1)
	table.Add(LDate, 7, 1)
	table.Add(NewVSpace(15), 8, 1)
	table.Add(DearT, 9, 1)
	table.Add(NewVSpace(15), 10, 1)
	table.Add(FirstP, 11, 1)
	table.Add(NewVSpace(15), 12, 1)
	table.Add(hp, 13, 1)
	table.Add(NewVSpace(15), 14, 1)
	table.Add(SecondP, 15, 1)
	table.Add(NewVSpace(15), 16, 1)
	table.Add(Salutation, 17, 1)
	table.Add(Correspondent, 18, 1)
	updateRemindPage(j, apt, nil,
		OStreetL, OCityStateZip, TenName, TStreet,
		TCityStateZip, LDate, DearT, FirstP, Details,
		SecondP, Salutation, Correspondent)

	b := gwu.NewButton("Print")
	b.AddEHandlerFunc(func(e gwu.Event) {

		updateRemindPage(j, apt, e,
			OStreetL, OCityStateZip, TenName, TStreet,
			TCityStateZip, LDate, DearT, FirstP, Details,
			SecondP, Salutation, Correspondent)

		list := getAptList(j)
		aptlb := gwu.NewListBox(list)

		Notify("I don't know how to Print yet", e)
		e.MarkDirty(aptlb)

	}, gwu.ETypeClick)

	tableA.Add(b, 0, 2)
	//table.Add(b, 0, 2)

	c.Add(table)
	c.AddVSpace(15)

	return c
}
