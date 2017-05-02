package main

import (
	"fmt"
	"time"

	"github.com/icza/gowut/gwu"
	"github.com/jinzhu/now"
)

func updateRentalRecord(j *jawaInfo, apt string,
	AptName, TenName, TenKey, MRent, Deposit, DueDay,
	Street, City, State, Zip, RentOwed, BounceOwed, LateOwed, WaterOwed,
	DepositOwed, NextDueDate, RentChargedThrou gwu.TextBox) {

	var day, month, year int
	var dayN, monthN, yearN int

	rec := j.Rental[apt]

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
	rec.Apartment = apt
	tname := TenName.Text()
	if tname != rec.Tenant {
		rec.Tenant = tname
		rec.TenantKey = makeTenantKey(tname)
	}
	ten, ok := j.Tenant[rec.TenantKey]
	if !ok {
		fmt.Printf("Tenant appears to be NEW\n")
		ten = renterRecord{}
		ten.Payment = map[string]payment{}
		ten.NextPaymentDue = now.New(timeNowRental().AddDate(0, 1, 0)).BeginningOfMonth()
		ten.Apartment = rec.Apartment
		ten.Name = rec.Tenant
	}

	fmt.Sscanf(MRent.Text(), "%f", &rec.Rent)
	fmt.Sscanf(Deposit.Text(), "%f", &rec.Deposit)
	fmt.Sscanf(DueDay.Text(), "%d", &rec.DueDay)
	rec.Street = Street.Text()
	rec.City = City.Text()
	rec.State = State.Text()
	rec.Zip = Zip.Text()

	fmt.Sscanf(RentOwed.Text(), "%f", &ten.RentOwed)
	fmt.Sscanf(BounceOwed.Text(), "%f", &ten.BounceOwed)
	fmt.Sscanf(LateOwed.Text(), "%f", &ten.LatePOwed)
	fmt.Sscanf(WaterOwed.Text(), "%f", &ten.WaterOwed)
	fmt.Sscanf(DepositOwed.Text(), "%f", &ten.DepositOwed)

	ten.NextPaymentDue = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	ten.RentChargedThru = time.Date(yearN, time.Month(monthN), dayN, 0, 0, 0, 0, time.Local)

	j.Tenant[rec.TenantKey] = ten
	j.Rental[apt] = rec
}

func updateRentalPage(j *jawaInfo, apt string, e gwu.Event,
	AptName, TenName, TenKey, MRent, Deposit, DueDay,
	Street, City, State, Zip, RentOwed, BounceOwed, LateOwed, WaterOwed,
	DepositOwed, NextDueDate, RentChargedThrou gwu.TextBox) {

	rec := j.Rental[apt]
	AptName.SetText(rec.Apartment)
	TenName.SetText(rec.Tenant)
	TenKey.SetText(rec.TenantKey)
	MRent.SetText(getRentalRent(j, apt))
	Deposit.SetText(getRentalDeposit(j, apt))
	DueDay.SetText(getDayOfMonthDue(j, apt))
	Street.SetText(rec.Street)
	City.SetText(rec.City)
	State.SetText(rec.State)
	Zip.SetText(rec.Zip)

	owed, _ := getTenantRentOwed(j, rec.TenantKey)
	RentOwed.SetText(owed)
	owed, _ = getTenantBounceOwed(j, rec.TenantKey)
	BounceOwed.SetText(owed)
	owed, _ = getTenantLateOwed(j, rec.TenantKey)
	LateOwed.SetText(owed)

	owed, _ = getTenantWaterOwed(j, rec.TenantKey)
	WaterOwed.SetText(owed)
	owed, _ = getTenantDepositOwed(j, rec.TenantKey)
	DepositOwed.SetText(owed)

	NextDueDate.SetText(getTenantRentDueDate(j, rec.TenantKey))
	RentChargedThrou.SetText(getTenantRentChargedThru(j, rec.TenantKey))
	if e != nil {
		//e.MarkDirty(aptlb)
		e.MarkDirty(AptName)
		e.MarkDirty(TenName)
		e.MarkDirty(TenKey)
		e.MarkDirty(MRent)
		e.MarkDirty(Deposit)
		e.MarkDirty(DueDay)
		e.MarkDirty(Street)
		e.MarkDirty(City)
		e.MarkDirty(State)
		e.MarkDirty(Zip)
		e.MarkDirty(RentOwed)
		e.MarkDirty(BounceOwed)
		e.MarkDirty(LateOwed)
		e.MarkDirty(WaterOwed)
		e.MarkDirty(DepositOwed)
		e.MarkDirty(NextDueDate)
		e.MarkDirty(RentChargedThrou)
	}
}

func buildEditRentals(j *jawaInfo) gwu.Panel {
	var AptName, TenName, TenKey, MRent, Deposit, DueDay gwu.TextBox
	var Street, City, State, Zip, RentOwed, BounceOwed, LateOwed gwu.TextBox
	var WaterOwed, DepositOwed gwu.TextBox
	var NextDueDate, RentChargedThrou gwu.TextBox
	var crb gwu.Button
	var cb gwu.CheckBox
	var NewNameLabel, tryAgainLabel gwu.Label
	var NewName gwu.TextBox
	var bN gwu.Button
	var aptlb gwu.ListBox

	c := gwu.NewPanel()

	tableA := gwu.NewTable()
	tableA.SetCellPadding(2)
	tableA.EnsureSize(2, 5)
	tableA.Add(gwu.NewLabel("Edit Rental:"), 0, 0)
	list := getAptList(j)
	apt := list[0]
	aptlb = gwu.NewListBox(list)
	aptlb.Style().SetFullWidth()
	aptlbHandler := func(e gwu.Event) {
		list := getAptList(j)
		apt = list[aptlb.SelectedIdx()]

		updateRentalPage(j, apt, e,
			AptName, TenName, TenKey, MRent, Deposit, DueDay,
			Street, City, State, Zip, RentOwed, BounceOwed, LateOwed,
			WaterOwed, DepositOwed, NextDueDate, RentChargedThrou)

		e.MarkDirty(aptlb)
	}
	aptlb.AddEHandlerFunc(aptlbHandler, gwu.ETypeChange)

	tableA.Add(aptlb, 0, 1)
	tableA.Add(gwu.NewLabel("...."), 0, 3)

	crb = gwu.NewButton("Create Rental Record")
	crb.AddEHandlerFunc(func(e gwu.Event) {
		fmt.Printf("Do something\n")
		tableA.Add(NewNameLabel, 1, 3)
		tableA.Add(NewName, 1, 4)
		tableA.Add(bN, 1, 5)
		NewName.Style().SetBackground(gwu.ClrAqua)
		NewNameLabel.Style().SetBackground(gwu.ClrAqua)
		bN.Style().SetBackground(gwu.ClrAqua)

		// create rental record using "new" for the apartment name and require it be
		//changed
		//????
		e.MarkDirty(NewNameLabel)
		e.MarkDirty(NewName)
		e.MarkDirty(bN)
		e.MarkDirty(tableA)

	}, gwu.ETypeClick)

	tableA.Add(crb, 0, 4)

	NewName = gwu.NewTextBox("new-apt-name")
	NewNameLabel = gwu.NewLabel("Set Name for Apartment")
	tryAgainLabel = gwu.NewLabel("That name already exists, please try again:")

	bN = gwu.NewButton("OK")
	bN.AddEHandlerFunc(func(e gwu.Event) {
		var name string
		fmt.Printf("Verify new name is present and is new then add record")
		fmt.Printf("get new name \n")
		labelToRemove := &NewNameLabel
		name = NewName.Text()
		_, ok := j.Rental[name]
		if ok {
			fmt.Printf("DUP, name[%s] in list\n", name)
			tableA.Remove(NewNameLabel)
			labelToRemove = &tryAgainLabel
			tryAgainLabel.Style().SetBackground(gwu.ClrAqua)
			tableA.Add(tryAgainLabel, 1, 3)
			e.MarkDirty(tryAgainLabel)
			e.MarkDirty(tableA)
			return
		}
		fmt.Printf("name [%s] not in list\n", name)
		j.Rental[name] = rentalRecord{Apartment: name}

		tableA.Remove(*labelToRemove)
		tableA.Remove(NewName)
		tableA.Remove(bN)

		if cb.State() {
			updateRentalRecord(j, name,
				AptName, TenName, TenKey, MRent, Deposit, DueDay,
				Street, City, State, Zip, RentOwed, BounceOwed, LateOwed, WaterOwed,
				DepositOwed, NextDueDate, RentChargedThrou)
		}
		// Zero out the tenant info
		rec := j.Rental[name]
		rec.Tenant = ""
		rec.TenantKey = ""
		j.Rental[name] = rec

		aptlb.ClearSelected()
		tableA.Remove(aptlb)
		list := getAptList(j)
		aptlb = gwu.NewListBox(list)
		aptlb.Style().SetFullWidth()
		nameIndex := strIndex(list, name)
		aptlb.SetSelected(nameIndex, true)
		aptlb.AddEHandlerFunc(aptlbHandler, gwu.ETypeChange)
		tableA.Add(aptlb, 0, 1)

		fmt.Printf("Added New name to list %+v\n", list)

		updateRentalPage(j, name, e,
			AptName, TenName, TenKey, MRent, Deposit, DueDay,
			Street, City, State, Zip, RentOwed, BounceOwed, LateOwed,
			WaterOwed, DepositOwed, NextDueDate, RentChargedThrou)

		e.MarkDirty(aptlb)
		e.MarkDirty(tableA)
	}, gwu.ETypeClick)

	cb = gwu.NewCheckBox("Base Record On Current?")
	cb.AddEHandlerFunc(func(e gwu.Event) { // I'm thinking I don't need to provide this handler
		if cb.State() {
			fmt.Printf("checked\n")
			return
		}
		fmt.Printf("Unchecked\n")
	}, gwu.ETypeChange)
	tableA.Add(cb, 0, 5)

	c.Add(tableA)

	c.AddVSpace(15)

	table := gwu.NewTable()
	table.SetCellPadding(2)
	table.EnsureSize(19, 4)

	table.Add(gwu.NewLabel("Apartment:"), 0, 0)
	table.Add(gwu.NewLabel("Tenant Name:"), 1, 0)
	table.Add(gwu.NewLabel("Tenant Key:"), 2, 0)
	table.Add(gwu.NewLabel("Monthly Rent:"), 3, 0)
	table.Add(gwu.NewLabel("Deposit:"), 4, 0)
	table.Add(gwu.NewLabel("Due day of month:"), 5, 0)
	table.Add(gwu.NewLabel("Street:"), 6, 0)
	table.Add(gwu.NewLabel("City:"), 7, 0)
	table.Add(gwu.NewLabel("State:"), 8, 0)
	table.Add(gwu.NewLabel("Zip:"), 9, 0)

	table.Add(gwu.NewLabel("RentOwed:"), 10, 0)
	table.Add(gwu.NewLabel("Bounce fee Owed:"), 11, 0)
	table.Add(gwu.NewLabel("Late fee Owed:"), 12, 0)
	table.Add(gwu.NewLabel("Water payment Owed:"), 13, 0)
	table.Add(gwu.NewLabel("Deposit Owed:"), 14, 0)
	table.Add(gwu.NewLabel("Next Rent payment due (mm-dd-yyyy):"), 15, 0)
	table.Add(gwu.NewLabel("Rent payment charged throu (mm-dd-yyyy):"), 16, 0)

	AptName = gwu.NewTextBox("")
	AptName.SetReadOnly(true)
	TenName = gwu.NewTextBox("")
	TenKey = gwu.NewTextBox("")
	TenKey.SetReadOnly(true)
	MRent = gwu.NewTextBox("")
	Deposit = gwu.NewTextBox("")
	DueDay = gwu.NewTextBox("")
	Street = gwu.NewTextBox("")
	City = gwu.NewTextBox("")
	State = gwu.NewTextBox("")
	Zip = gwu.NewTextBox("")

	RentOwed = gwu.NewTextBox("0.00")
	BounceOwed = gwu.NewTextBox("0.00")
	LateOwed = gwu.NewTextBox("0.00")
	WaterOwed = gwu.NewTextBox("0.00")
	DepositOwed = gwu.NewTextBox("0.00")

	NextDueDate = gwu.NewTextBox("mm-dd-yyyy")
	RentChargedThrou = gwu.NewTextBox("mm-dd-yyyy")

	AptName.Style().SetWidthPx(260) //??????? WGA
	table.Add(AptName, 0, 1)
	table.Add(TenName, 1, 1)
	table.Add(TenKey, 2, 1)
	table.Add(MRent, 3, 1)
	table.Add(Deposit, 4, 1)
	table.Add(DueDay, 5, 1)
	table.Add(Street, 6, 1)
	table.Add(City, 7, 1)
	table.Add(State, 8, 1)
	table.Add(Zip, 9, 1)
	table.Add(RentOwed, 10, 1)
	table.Add(BounceOwed, 11, 1)
	table.Add(LateOwed, 12, 1)
	table.Add(WaterOwed, 13, 1)
	table.Add(DepositOwed, 14, 1)
	table.Add(NextDueDate, 15, 1)
	table.Add(RentChargedThrou, 16, 1)
	updateRentalPage(j, apt, nil,
		AptName, TenName, TenKey, MRent, Deposit, DueDay,
		Street, City, State, Zip, RentOwed, BounceOwed, LateOwed,
		WaterOwed, DepositOwed, NextDueDate, RentChargedThrou)

	b := gwu.NewButton("Submit")
	b.AddEHandlerFunc(func(e gwu.Event) {

		updateRentalRecord(j, apt,
			AptName, TenName, TenKey, MRent, Deposit, DueDay,
			Street, City, State, Zip, RentOwed, BounceOwed, LateOwed, WaterOwed,
			DepositOwed, NextDueDate, RentChargedThrou)

		updateRentalPage(j, apt, e,
			AptName, TenName, TenKey, MRent, Deposit, DueDay,
			Street, City, State, Zip, RentOwed, BounceOwed, LateOwed,
			WaterOwed, DepositOwed, NextDueDate, RentChargedThrou)

		list := getAptList(j)
		aptlb := gwu.NewListBox(list)

		e.MarkDirty(aptlb)

	}, gwu.ETypeClick)

	tableA.Add(b, 0, 2)
	//table.Add(b, 0, 2)

	cbdTable := gwu.NewTable()
	cbdTable.SetCellPadding(2)
	cbdTable.EnsureSize(1, 2)

	var ays gwu.CheckBox

	cbd := gwu.NewCheckBox("Delete the current Record?")
	cbd.AddEHandlerFunc(func(e gwu.Event) {
		if cbd.State() {
			cbd.Style().SetBackground(gwu.ClrGreen)
			ays.Style().SetBackground(gwu.ClrAqua)
			cbdTable.Add(ays, 0, 1)
		} else {
			cbd.Style().SetBackground("")
			ays.SetState(false)
			cbdTable.Remove(ays)
		}
		e.MarkDirty(cbd)
		e.MarkDirty(ays)
		e.MarkDirty(cbdTable)
	}, gwu.ETypeClick)
	cbdTable.Add(cbd, 0, 0)

	ays = gwu.NewCheckBox("ARE YOU SURE?")
	ays.AddEHandlerFunc(func(e gwu.Event) {
		if ays.State() {
			ays.Style().SetBackground(gwu.ClrGreen)

			// Delete the record here
			aptlb.ClearSelected()
			tableA.Remove(aptlb)
			delete(j.Rental, apt)
			list := getAptList(j)
			apt = list[0]
			aptlb = gwu.NewListBox(list)
			aptlb.Style().SetFullWidth()
			//aptlb.SetSelected(0, true)
			aptlb.AddEHandlerFunc(aptlbHandler, gwu.ETypeChange)
			tableA.Add(aptlb, 0, 1)

			fmt.Printf("Added New name to list %+v\n", list)

			updateRentalPage(j, apt, e,
				AptName, TenName, TenKey, MRent, Deposit, DueDay,
				Street, City, State, Zip, RentOwed, BounceOwed, LateOwed,
				WaterOwed, DepositOwed, NextDueDate, RentChargedThrou)

			cbd.Style().SetBackground("")
			cbd.SetState(false)
			ays.SetState(false)
			cbdTable.Remove(ays)

			e.MarkDirty(cbd)
			e.MarkDirty(ays)
			e.MarkDirty(aptlb)
			e.MarkDirty(tableA)
			e.MarkDirty(cbdTable)

		} else {
			ays.SetState(true)
			//ays.Style().SetBackground("")
		}
		e.MarkDirty(ays)
	}, gwu.ETypeClick)

	c.Add(table)
	c.AddVSpace(15)
	c.Add(cbdTable)

	return c
}
