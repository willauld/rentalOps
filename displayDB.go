package main

import (
	"fmt"

	"github.com/icza/gowut/gwu"
)

func updateDBBox(j *jawaInfo, e gwu.Event, p gwu.TextBox) {

	curText := "Rental Ops DB:\n"
	curText += "==============\n"
	curText += "Common Info:\n"
	curText += "==============\n"
	v := j.CI
	curText += fmt.Sprintf("Office manager:       %s\n", v.OfficeManager)
	curText += fmt.Sprintf("Office street:        %s\n", v.OfficeStreet)
	curText += fmt.Sprintf("Office city:          %s\n", v.OfficeCity)
	curText += fmt.Sprintf("Office state:         %s\n", v.OfficeState)
	curText += fmt.Sprintf("Office state:         %s\n", v.OfficeZip)
	curText += fmt.Sprintf("Correspondence title: %s\n", v.CorrespondenceTitle)
	curText += fmt.Sprintf("Late fee:             %-6.2f\n", v.LateFee)
	curText += fmt.Sprintf("Incurr fee after:     %d\n", v.IncurrAfter)
	curText += fmt.Sprintf("Bounce fee:           %-6.2f\n", v.BounceFee)
	curText += "==================\n"
	curText += "Apartment Records:\n"
	curText += "==================\n"
	for k, v := range j.Rental {
		curText += "==================\n"
		curText += fmt.Sprintf("Apartment %s:", k)
		curText += "==================\n"
		curText += fmt.Sprintf("Apartment: %s\n", v.Apartment)
		curText += fmt.Sprintf("Tenant:    %s\n", v.Tenant)
		curText += fmt.Sprintf("TenantKey: %s\n", v.TenantKey)
		curText += fmt.Sprintf("Rent:      %-6.2f\n", v.Rent)
		curText += fmt.Sprintf("Deposit:   %-6.2f\n", v.Deposit)
		curText += fmt.Sprintf("Due day:   %d\n", v.DueDay)
		curText += fmt.Sprintf("Street:    %s\n", v.Street)
		curText += fmt.Sprintf("City:      %s\n", v.City)
		curText += fmt.Sprintf("State:     %s\n", v.State)
		curText += fmt.Sprintf("Zip:       %s\n", v.Zip)
	}
	curText += "==================\n"
	curText += "Tenant Records:\n"
	curText += "==================\n"
	for k, v := range j.Tenant {
		curText += "==================\n"
		curText += fmt.Sprintf("Tenant %s:\n", k)
		curText += "==================\n"
		curText += fmt.Sprintf("Apartment:         %s\n", v.Apartment)
		curText += fmt.Sprintf("Name:              %s\n", v.Name)
		curText += fmt.Sprintf("Rent owed:         %-6.2f\n", v.RentOwed)
		curText += fmt.Sprintf("Bounce fee owed:   %-6.2f\n", v.BounceOwed)
		curText += fmt.Sprintf("Late fee owed:     %-6.2f\n", v.LatePOwed)
		curText += fmt.Sprintf("Water fee owed:    %-6.2f\n", v.WaterOwed)
		curText += fmt.Sprintf("Deposit owed:      %-6.2f\n", v.DepositOwed)
		curText += fmt.Sprintf("Next payment due:  %v\n", v.NextPaymentDue)
		curText += fmt.Sprintf("Rent charged thru: %v\n", v.RentChargedThru)
		curText += "        ==================\n"
		curText += "        Payment Records:\n"
		curText += "        ==================\n"
		payments := GetKeys(v.Payment)
		for i := 0; i < len(payments); i++ {
			k := payments[i]
			v2 := v.Payment[k]

			curText += "        ==================\n"
			curText += fmt.Sprintf("        Payment: %s\n", k)
			curText += "        ==================\n"
			curText += fmt.Sprintf("        Check Amount: %-6.2f\n", v2.Amount)
			curText += fmt.Sprintf("        Check Date:   %v\n", v2.Date)
			curText += fmt.Sprintf("                Rent:       %-6.2f\n", v2.Rent)
			curText += fmt.Sprintf("                Late:       %-6.2f\n", v2.Late)
			curText += fmt.Sprintf("                Bounce:     %-6.2f\n", v2.Bounce)
			curText += fmt.Sprintf("                Water:      %-6.2f\n", v2.Water)
			curText += fmt.Sprintf("                Deposit:    %-6.2f\n", v2.Deposit)

		}
	}
	p.SetText(curText)
	if e != nil {
		e.MarkDirty(p)
	}
}

func buildDisplayDB(j *jawaInfo) (gwu.Panel, gwu.TextBox) {
	// display Common
	var DisplayBox gwu.TextBox

	c := gwu.NewPanel()
	stb := gwu.NewTextBox("")
	stb.Style().SetWidthPx(1).SetHeightPx(1)
	stb.AddEHandlerFunc(func(e gwu.Event) {
		// same as update button
		updateDBBox(j, e, DisplayBox)
		e.MarkDirty(DisplayBox)
		Notify("Focus is on Display DB Tab", e)
	}, gwu.ETypeFocus)
	c.Add(stb)

	tableA := gwu.NewTable()
	tableA.SetCellPadding(2)
	tableA.EnsureSize(1, 3)
	tableA.Add(gwu.NewLabel("Update:"), 0, 0)
	updb := gwu.NewButton("Update Me!")
	updb.AddEHandlerFunc(func(e gwu.Event) {
		updateDBBox(j, e, DisplayBox)
		e.MarkDirty(DisplayBox)
		Notify("dumpDB Update complete", e)
	}, gwu.ETypeChange)
	tableA.Add(updb, 0, 1)

	c.Add(tableA)

	DisplayBox = gwu.NewTextBox("")
	DisplayBox.Style().SetWidthPx(500).SetHeightPx(500)
	DisplayBox.SetRows(10)
	DisplayBox.SetCols(70)
	c.Add(DisplayBox)

	updateDBBox(j, nil, DisplayBox)

	return c, stb
}
