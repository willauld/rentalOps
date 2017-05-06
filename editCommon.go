package main

import (
	"fmt"

	"github.com/icza/gowut/gwu"
)

func updateCommonRecord(j *jawaInfo,
	OfficeStreet, OfficeCity, OfficeState, OfficeZip,
	OfficeManager, CTitle, LateFee, LateDays, BounceFee gwu.TextBox) {

	rec := j.CI

	rec.OfficeManager = OfficeManager.Text()
	rec.OfficeStreet = OfficeStreet.Text()
	rec.OfficeCity = OfficeCity.Text()
	rec.OfficeState = OfficeState.Text()
	rec.OfficeZip = OfficeZip.Text()
	rec.CorrespondenceTitle = CTitle.Text()
	fmt.Sscanf(LateFee.Text(), "%f", &rec.LateFee)
	fmt.Sscanf(LateDays.Text(), "%d", &rec.IncurrAfter)
	fmt.Sscanf(BounceFee.Text(), "%f", &rec.BounceFee)
	j.CI = rec
}

func updateCommonPage(j *jawaInfo, e gwu.Event,
	OfficeStreet, OfficeCity, OfficeState, OfficeZip,
	OfficeManager, CTitle, LateFee, LateDays,
	BounceFee gwu.TextBox) {

	rec := j.CI
	OfficeStreet.SetText(rec.OfficeStreet)
	OfficeCity.SetText(rec.OfficeCity)
	OfficeState.SetText(rec.OfficeState)
	OfficeZip.SetText(rec.OfficeZip)
	OfficeManager.SetText(rec.OfficeManager)
	CTitle.SetText(rec.CorrespondenceTitle)
	LateFee.SetText(fmt.Sprintf("%-6.2f", rec.LateFee))
	LateDays.SetText(fmt.Sprintf("%d", rec.IncurrAfter))
	BounceFee.SetText(fmt.Sprintf("%-6.2f", rec.BounceFee))

	if e != nil {
		e.MarkDirty(OfficeStreet)
		e.MarkDirty(OfficeCity)
		e.MarkDirty(OfficeState)
		e.MarkDirty(OfficeZip)
		e.MarkDirty(OfficeManager)
		e.MarkDirty(CTitle)
		e.MarkDirty(LateFee)
		e.MarkDirty(LateDays)
		e.MarkDirty(BounceFee)
	}
}

func buildEditCommon(j *jawaInfo) gwu.Panel {
	var OfficeStreet, OfficeCity, OfficeState, OfficeZip gwu.TextBox
	var OfficeManager, CTitle, LateFee, LateDays, BounceFee gwu.TextBox

	c := gwu.NewPanel()
	c.AddEHandlerFunc(func(e gwu.Event) {

		updateCommonPage(j, e,
			OfficeStreet, OfficeCity, OfficeState, OfficeZip,
			OfficeManager, CTitle, LateFee, LateDays, BounceFee)

		Notify("EditCommon Focus Event happened", e)
		fmt.Printf("inside the editCommon state change handler\n")
	}, gwu.ETypeStateChange /*gwu.ETypeFocus*/ /*gwu.ETypeClick*/)

	tableA := gwu.NewTable()
	tableA.SetCellPadding(2)
	tableA.EnsureSize(2, 5)
	tableA.Add(gwu.NewLabel("Edit Rental:"), 0, 0)
	tableA.Add(gwu.NewLabel("...."), 0, 3)

	c.Add(tableA)

	c.AddVSpace(15)

	table := gwu.NewTable()
	table.SetCellPadding(2)
	table.EnsureSize(9, 4)

	table.Add(gwu.NewLabel("Office Street:"), 0, 0)
	table.Add(gwu.NewLabel("Office City:"), 1, 0)
	table.Add(gwu.NewLabel("Office State:"), 2, 0)
	table.Add(gwu.NewLabel("Office Zip:"), 3, 0)
	table.Add(gwu.NewLabel("Office Manager:"), 4, 0)
	table.Add(gwu.NewLabel("Correspondence Title:"), 5, 0)
	table.Add(gwu.NewLabel("Late Fee"), 6, 0)
	table.Add(gwu.NewLabel("Late after ? days:"), 7, 0)
	table.Add(gwu.NewLabel("Bounce Fee:"), 8, 0)

	OfficeStreet = gwu.NewTextBox("")
	OfficeStreet.Style().SetWidthPx(260)
	OfficeCity = gwu.NewTextBox("")
	OfficeCity.Style().SetWidthPx(260)
	OfficeState = gwu.NewTextBox("")
	OfficeState.Style().SetWidthPx(260)
	OfficeZip = gwu.NewTextBox("")
	OfficeZip.Style().SetWidthPx(260)
	OfficeManager = gwu.NewTextBox("")
	OfficeManager.Style().SetWidthPx(260)
	CTitle = gwu.NewTextBox("")
	CTitle.Style().SetWidthPx(260)
	LateFee = gwu.NewTextBox("")
	LateFee.Style().SetWidthPx(260)
	LateDays = gwu.NewTextBox("")
	LateDays.Style().SetWidthPx(260)
	BounceFee = gwu.NewTextBox("")
	BounceFee.Style().SetWidthPx(260)

	table.Add(OfficeStreet, 0, 1)
	table.Add(OfficeCity, 1, 1)
	table.Add(OfficeState, 2, 1)
	table.Add(OfficeZip, 3, 1)
	table.Add(OfficeManager, 4, 1)
	table.Add(CTitle, 5, 1)
	table.Add(LateFee, 6, 1)
	table.Add(LateDays, 7, 1)
	table.Add(BounceFee, 8, 1)

	updateCommonPage(j, nil,
		OfficeStreet, OfficeCity, OfficeState, OfficeZip,
		OfficeManager, CTitle, LateFee, LateDays, BounceFee)

	b := gwu.NewButton("Submit")
	b.AddEHandlerFunc(func(e gwu.Event) {

		updateCommonRecord(j,
			OfficeStreet, OfficeCity, OfficeState, OfficeZip,
			OfficeManager, CTitle, LateFee, LateDays, BounceFee)

		updateCommonPage(j, e,
			OfficeStreet, OfficeCity, OfficeState, OfficeZip,
			OfficeManager, CTitle, LateFee, LateDays, BounceFee)

	}, gwu.ETypeClick)

	tableA.Add(b, 0, 2)
	//table.Add(b, 0, 2)

	c.Add(table)

	return c
}
