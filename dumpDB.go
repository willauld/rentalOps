package main

import "fmt"
import "github.com/icza/gowut/gwu"

func updateDBPage(j *jawaInfo, e gwu.Event, p gwu.Panel) {

	p.Clear()
	p.Add(gwu.NewLabel("Rental Ops DB:"))
	p.Add(gwu.NewLabel("=============="))
	p.Add(gwu.NewLabel("Common Info:"))
	p.Add(gwu.NewLabel("=============="))
	v := j.CI
	p.Add(gwu.NewLabel(fmt.Sprintf("Office manager:       %s", v.OfficeManager)))
	p.Add(gwu.NewLabel(fmt.Sprintf("Office street:        %s", v.OfficeStreet)))
	p.Add(gwu.NewLabel(fmt.Sprintf("Office city:          %s", v.OfficeCity)))
	p.Add(gwu.NewLabel(fmt.Sprintf("Office state:         %s", v.OfficeState)))
	p.Add(gwu.NewLabel(fmt.Sprintf("Office state:         %s", v.OfficeZip)))
	p.Add(gwu.NewLabel(fmt.Sprintf("Correspondence title: %s", v.CorrespondenceTitle)))
	p.Add(gwu.NewLabel(fmt.Sprintf("Late fee:             %-6.2f", v.LateFee)))
	p.Add(gwu.NewLabel(fmt.Sprintf("Incurr fee after:     %d", v.IncurrAfter)))
	p.Add(gwu.NewLabel(fmt.Sprintf("Bounce fee:           %-6.2f", v.BounceFee)))
	p.Add(gwu.NewLabel("=================="))
	p.Add(gwu.NewLabel("Apartment Records:"))
	p.Add(gwu.NewLabel("=================="))
	for k, v := range j.Rental {
		p.Add(gwu.NewLabel("=================="))
		p.Add(gwu.NewLabel(fmt.Sprintf("Apartment %s:", k)))
		p.Add(gwu.NewLabel("=================="))
		p.Add(gwu.NewLabel(fmt.Sprintf("Apartment: %s", v.Apartment)))
		p.Add(gwu.NewLabel(fmt.Sprintf("Tenant:    %s", v.Tenant)))
		p.Add(gwu.NewLabel(fmt.Sprintf("TenantKey: %s", v.TenantKey)))
		p.Add(gwu.NewLabel(fmt.Sprintf("Rent:      %-6.2f", v.Rent)))
		p.Add(gwu.NewLabel(fmt.Sprintf("Deposit:   %-6.2f", v.Deposit)))
		p.Add(gwu.NewLabel(fmt.Sprintf("Due day:   %d", v.DueDay)))
		p.Add(gwu.NewLabel(fmt.Sprintf("Street:    %s", v.Street)))
		p.Add(gwu.NewLabel(fmt.Sprintf("City:      %s", v.City)))
		p.Add(gwu.NewLabel(fmt.Sprintf("State:     %s", v.State)))
		p.Add(gwu.NewLabel(fmt.Sprintf("Zip:       %s", v.Zip)))
	}
	p.Add(gwu.NewLabel("=================="))
	p.Add(gwu.NewLabel("Tenant Records:"))
	p.Add(gwu.NewLabel("=================="))
	for k, v := range j.Tenant {
		p.Add(gwu.NewLabel("=================="))
		p.Add(gwu.NewLabel(fmt.Sprintf("Tenant %s:", k)))
		p.Add(gwu.NewLabel("=================="))
		p.Add(gwu.NewLabel(fmt.Sprintf("Apartment:         %s", v.Apartment)))
		p.Add(gwu.NewLabel(fmt.Sprintf("Name:              %s", v.Name)))
		p.Add(gwu.NewLabel(fmt.Sprintf("Rent owed:         %-6.2f", v.RentOwed)))
		p.Add(gwu.NewLabel(fmt.Sprintf("Bounce fee owed:   %-6.2f", v.BounceOwed)))
		p.Add(gwu.NewLabel(fmt.Sprintf("Late fee owed:     %-6.2f", v.LatePOwed)))
		p.Add(gwu.NewLabel(fmt.Sprintf("Water fee owed:    %-6.2f", v.WaterOwed)))
		p.Add(gwu.NewLabel(fmt.Sprintf("Deposit owed:      %-6.2f", v.DepositOwed)))
		p.Add(gwu.NewLabel(fmt.Sprintf("Next payment due:  %v", v.NextPaymentDue)))
		p.Add(gwu.NewLabel(fmt.Sprintf("Rent charged thru: %v", v.RentChargedThru)))
		p.Add(gwu.NewLabel("        =================="))
		p.Add(gwu.NewLabel("        Payment Records:"))
		p.Add(gwu.NewLabel("        =================="))
		for k, v2 := range v.Payment {
			p.Add(gwu.NewLabel("        =================="))
			p.Add(gwu.NewLabel(fmt.Sprintf("        Payment: %s", k)))
			p.Add(gwu.NewLabel("        =================="))
			p.Add(gwu.NewLabel(fmt.Sprintf("        Check Amount: %-6.2f", v2.Amount)))
			p.Add(gwu.NewLabel(fmt.Sprintf("        Check Date:   %v", v2.Date)))
			p.Add(gwu.NewLabel(fmt.Sprintf("                Rent:       %-6.2f", v2.Rent)))
			p.Add(gwu.NewLabel(fmt.Sprintf("                Late:       %-6.2f", v2.Late)))
			p.Add(gwu.NewLabel(fmt.Sprintf("                Bounce:     %-6.2f", v2.Bounce)))
			p.Add(gwu.NewLabel(fmt.Sprintf("                Water:      %-6.2f", v2.Water)))
			p.Add(gwu.NewLabel(fmt.Sprintf("                Deposit:    %-6.2f", v2.Deposit)))

		}

	}
	if e != nil {
		e.MarkDirty(p)
	}
}

func buildDisplayDB(j *jawaInfo) (gwu.Panel, gwu.TextBox) {
	// display Common
	c := gwu.NewPanel()
	stb := gwu.NewTextBox("")
	stb.Style().SetWidthPx(1).SetHeightPx(1)
	stb.AddEHandlerFunc(func(e gwu.Event) {
		// same as update button
		updateDBPage(j, e, c)
		e.MarkDirty(c)
		Notify("Focus is on Dump DB Tab", e)
	}, gwu.ETypeFocus)
	c.Add(stb)

	tableA := gwu.NewTable()
	tableA.SetCellPadding(2)
	tableA.EnsureSize(1, 3)
	np := gwu.NewPanel()
	tableA.Add(gwu.NewLabel("Update:"), 0, 0)
	updb := gwu.NewButton("Update Me!")
	updb.AddEHandlerFunc(func(e gwu.Event) {
		updateDBPage(j, e, np)
		e.MarkDirty(np)
		Notify("dumpDB Update complete", e)
	}, gwu.ETypeChange)
	tableA.Add(updb, 0, 1)

	c.Add(tableA)
	//np := gwu.NewPanal()
	c.Add(np)

	updateDBPage(j, nil, np)

	return c, stb
}
