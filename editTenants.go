package main

import (
	"fmt"
	"time"

	"github.com/icza/gowut/gwu"
	//"github.com/jinzhu/now"
)

func updateTenantRecord(j *jawaInfo, ten string,
	AptName, TenName, RentOwed, BounceOwed, LateOwed, WaterOwed,
	DepositOwed, NextDueDate, RentChargedThrou gwu.TextBox) {

	var day, month, year int
	var dayN, monthN, yearN int

	rec := j.Tenant[ten]

	n, err := fmt.Sscanf(NextDueDate.Text(), "%d-%d-%d\n", &month, &day, &year)
	if err != nil || n != 3 {
		fmt.Printf("Due date format is incorrect, please try again\n")
		//TODO: need this not to be pop up for the like // notification
		return
	}
	n, err = fmt.Sscanf(RentChargedThrou.Text(), "%d-%d-%d\n", &monthN, &dayN, &yearN)
	if err != nil || n != 3 {
		fmt.Printf("Rent charged throu date format is incorrect, please try again\n")
		//TODO: need this not to be pop up for the like // notification
		return
	}
	rec.Apartment = AptName.Text() // DON"T UPDATE THIS FIELD, IT
	// is too closely tied to RENTAL record
	rec.Name = TenName.Text()

	fmt.Sscanf(RentOwed.Text(), "%f", &rec.RentOwed)
	fmt.Sscanf(BounceOwed.Text(), "%f", &rec.BounceOwed)
	fmt.Sscanf(LateOwed.Text(), "%f", &rec.LatePOwed)
	fmt.Sscanf(WaterOwed.Text(), "%f", &rec.WaterOwed)
	fmt.Sscanf(DepositOwed.Text(), "%f", &rec.DepositOwed)

	rec.NextPaymentDue = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	rec.RentChargedThru = time.Date(yearN, time.Month(monthN), dayN, 0, 0, 0, 0, time.Local)

	j.Tenant[ten] = rec
}

func updateTenantPage(j *jawaInfo, ten string, e gwu.Event,
	AptName, TenName, RentOwed, BounceOwed, LateOwed, WaterOwed,
	DepositOwed, NextDueDate, RentChargedThrou, paymentRecords gwu.TextBox) {

	rec := j.Tenant[ten]
	AptName.SetText(rec.Apartment)
	TenName.SetText(rec.Name)

	owed, _ := getTenantRentOwed(j, ten)
	RentOwed.SetText(owed)
	owed, _ = getTenantBounceOwed(j, ten)
	BounceOwed.SetText(owed)
	owed, _ = getTenantLateOwed(j, ten)
	LateOwed.SetText(owed)

	owed, _ = getTenantWaterOwed(j, ten)
	WaterOwed.SetText(owed)
	owed, _ = getTenantDepositOwed(j, ten)
	DepositOwed.SetText(owed)

	NextDueDate.SetText(getTenantRentDueDate(j, ten))
	RentChargedThrou.SetText(getTenantRentChargedThru(j, ten))
	paymentRecords.SetText(getTenantNumPayments(j, ten))
	if e != nil {
		e.MarkDirty(AptName)
		e.MarkDirty(TenName)
		e.MarkDirty(RentOwed)
		e.MarkDirty(BounceOwed)
		e.MarkDirty(LateOwed)
		e.MarkDirty(WaterOwed)
		e.MarkDirty(DepositOwed)
		e.MarkDirty(NextDueDate)
		e.MarkDirty(RentChargedThrou)
		e.MarkDirty(paymentRecords)
	}
}

func buildEditTenants(j *jawaInfo) (gwu.Panel, gwu.TextBox) {
	var aptname, tenname gwu.TextBox
	var rentowed, bounceowed, lateowed gwu.TextBox
	var waterowed, depositowed gwu.TextBox
	var nextduedate, rentchargedthrou gwu.TextBox
	var tenlb gwu.ListBox
	var paymentRecords gwu.TextBox
	var tablea gwu.Table
	var tenlbhandler func(e gwu.Event)
	var ten string // The master var for which tenant is active

	c := gwu.NewPanel()
	stb := gwu.NewTextBox("")
	stb.Style().SetWidthPx(1).SetHeightPx(1)
	stb.AddEHandlerFunc(func(e gwu.Event) {

		//meList := getKeyList(j.Rental) //TODO: ***WANT TO MOVE TO THIS not getAptList()
		meList := getTenantList(j)
		tablea.Remove(tenlb)
		err := UpdateListBox(&tenlb, &ten, meList, tenlbhandler)
		if err != nil {
			fmt.Printf("EditTenant update ListBox failed: %v\n", err)
			return
		}
		tablea.Add(tenlb, 0, 1)

		updateTenantPage(j, ten, e,
			aptname, tenname, rentowed, bounceowed, lateowed,
			waterowed, depositowed, nextduedate, rentchargedthrou, paymentRecords)

		e.MarkDirty(tenlb)
		e.MarkDirty(tablea)
		Notify("Focus is on Edit Tenant Tab", e)
	}, gwu.ETypeFocus)
	c.Add(stb)

	tablea = gwu.NewTable()
	tablea.SetCellPadding(2)
	tablea.EnsureSize(2, 5)
	tablea.Add(gwu.NewLabel("Edit Tenant:"), 0, 0)
	list := getTenantList(j)
	ten = list[0]
	tenlb = gwu.NewListBox(list)
	tenlb.Style().SetFullWidth()
	tenlbhandler = func(e gwu.Event) {
		list := getTenantList(j)
		ten = list[tenlb.SelectedIdx()]

		updateTenantPage(j, ten, e,
			aptname, tenname, rentowed, bounceowed, lateowed,
			waterowed, depositowed, nextduedate, rentchargedthrou, paymentRecords)

		e.MarkDirty(tenlb)
	}
	tenlb.AddEHandlerFunc(tenlbhandler, gwu.ETypeChange)

	tablea.Add(tenlb, 0, 1)
	tablea.Add(gwu.NewLabel("...."), 0, 3)

	c.Add(tablea)

	c.AddVSpace(15)

	table := gwu.NewTable()
	table.SetCellPadding(2)
	table.EnsureSize(10, 4)

	table.Add(gwu.NewLabel("Apartment:"), 0, 0)
	table.Add(gwu.NewLabel("Tenant name:"), 1, 0)

	table.Add(gwu.NewLabel("Rent owed:"), 2, 0)
	table.Add(gwu.NewLabel("Bounce fee owed:"), 3, 0)
	table.Add(gwu.NewLabel("Late fee owed:"), 4, 0)
	table.Add(gwu.NewLabel("Water payment owed:"), 5, 0)
	table.Add(gwu.NewLabel("Deposit owed:"), 6, 0)
	table.Add(gwu.NewLabel("Next rent payment due (mm-dd-yyyy):"), 7, 0)
	table.Add(gwu.NewLabel("Rent payment charged throu (mm-dd-yyyy):"), 8, 0)
	table.Add(gwu.NewLabel("Payment Records:"), 9, 0)

	aptname = gwu.NewTextBox("")
	aptname.SetReadOnly(true)
	tenname = gwu.NewTextBox("")
	tenname.SetReadOnly(true)

	rentowed = gwu.NewTextBox("0.00")
	bounceowed = gwu.NewTextBox("0.00")
	lateowed = gwu.NewTextBox("0.00")
	waterowed = gwu.NewTextBox("0.00")
	depositowed = gwu.NewTextBox("0.00")

	nextduedate = gwu.NewTextBox("mm-dd-yyyy")
	rentchargedthrou = gwu.NewTextBox("mm-dd-yyyy")
	paymentRecords = gwu.NewTextBox("0")
	paymentRecords.SetReadOnly(true)

	aptname.Style().SetWidthPx(260) //??????? wga
	table.Add(aptname, 0, 1)
	table.Add(tenname, 1, 1)
	table.Add(rentowed, 2, 1)
	table.Add(bounceowed, 3, 1)
	table.Add(lateowed, 4, 1)
	table.Add(waterowed, 5, 1)
	table.Add(depositowed, 6, 1)
	table.Add(nextduedate, 7, 1)
	table.Add(rentchargedthrou, 8, 1)
	table.Add(paymentRecords, 9, 1)

	updateTenantPage(j, ten, nil,
		aptname, tenname, rentowed, bounceowed, lateowed,
		waterowed, depositowed, nextduedate, rentchargedthrou, paymentRecords)

	b := gwu.NewButton("submit")
	b.AddEHandlerFunc(func(e gwu.Event) {

		updateTenantRecord(j, ten,
			aptname, tenname, rentowed, bounceowed, lateowed,
			waterowed, depositowed, nextduedate, rentchargedthrou)

		updateTenantPage(j, ten, e,
			aptname, tenname, rentowed, bounceowed, lateowed,
			waterowed, depositowed, nextduedate, rentchargedthrou, paymentRecords)

		list := getTenantList(j)
		tenlb := gwu.NewListBox(list)

		e.MarkDirty(tenlb)

	}, gwu.ETypeClick)

	tablea.Add(b, 0, 2)
	//table.Add(b, 0, 2)
	u := gwu.NewButton("Update Page")
	u.AddEHandlerFunc(func(e gwu.Event) {

		updateTenantPage(j, ten, e,
			aptname, tenname, rentowed, bounceowed, lateowed,
			waterowed, depositowed, nextduedate, rentchargedthrou, paymentRecords)

		list := getTenantList(j)
		tenlb := gwu.NewListBox(list)

		e.MarkDirty(tenlb)

	}, gwu.ETypeClick)
	tablea.Add(u, 0, 4)

	cbdtable := gwu.NewTable()
	cbdtable.SetCellPadding(2)
	cbdtable.EnsureSize(1, 2)

	var ays gwu.CheckBox

	cbd := gwu.NewCheckBox("delete the current record?")
	cbd.AddEHandlerFunc(func(e gwu.Event) {
		if cbd.State() {
			cbd.Style().SetBackground(gwu.ClrGreen)
			ays.Style().SetBackground(gwu.ClrAqua)
			cbdtable.Add(ays, 0, 1)
		} else {
			cbd.Style().SetBackground("")
			ays.SetState(false)
			cbdtable.Remove(ays)
		}
		e.MarkDirty(cbd)
		e.MarkDirty(ays)
		e.MarkDirty(cbdtable)
	}, gwu.ETypeClick)
	cbdtable.Add(cbd, 0, 0)

	ays = gwu.NewCheckBox("are you sure?")
	ays.AddEHandlerFunc(func(e gwu.Event) {
		if ays.State() {
			ays.Style().SetBackground(gwu.ClrGreen)

			// delete the record here
			tenlb.ClearSelected()
			tablea.Remove(tenlb)
			delete(j.Tenant, ten)
			list := getTenantList(j)
			ten = list[0]
			tenlb = gwu.NewListBox(list)
			tenlb.Style().SetFullWidth()
			//aptlb.setselected(0, true)
			tenlb.AddEHandlerFunc(tenlbhandler, gwu.ETypeChange)
			tablea.Add(tenlb, 0, 1)

			fmt.Printf("added new name to list %+v\n", list)

			updateTenantPage(j, ten, e,
				aptname, tenname, rentowed, bounceowed, lateowed,
				waterowed, depositowed, nextduedate, rentchargedthrou, paymentRecords)

			cbd.Style().SetBackground("")
			cbd.SetState(false)
			ays.SetState(false)
			cbdtable.Remove(ays)

			e.MarkDirty(cbd)
			e.MarkDirty(ays)
			e.MarkDirty(tenlb)
			e.MarkDirty(tablea)
			e.MarkDirty(cbdtable)

		} else {
			ays.SetState(true)
			//ays.Style().SetBackGround("")
		}
		e.MarkDirty(ays)
	}, gwu.ETypeClick)

	c.Add(table)
	c.AddVSpace(15)
	c.Add(cbdtable)

	return c, stb
}
