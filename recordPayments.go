package main

import (
	"fmt"
	"time"

	"github.com/icza/gowut/gwu"
	//"github.com/willauld/gowut/gwu"
)

func getAptList(ji *jawaInfo) []string {
	return *ji.CI.Apartments
}
func getRentalRent(j *jawaInfo, apt string) string {
	return fmt.Sprintf("%-6.2f", j.Rental[apt].Rent)
}
func getTenantRentDue(j *jawaInfo, apt string) (string, bool) {
	val := j.Tenant[j.Rental[apt].TenantKey].RentOwed
	return fmt.Sprintf("%-6.2f", val), !(val < 0)
}
func getTenantLateDue(j *jawaInfo, apt string) (string, bool) {
	val := j.Tenant[j.Rental[apt].TenantKey].LatePOwed
	return fmt.Sprintf("%-6.2f", val), !(val < 0)
}
func getTenantBounceDue(j *jawaInfo, apt string) (string, bool) {
	val := j.Tenant[j.Rental[apt].TenantKey].BounceOwed
	return fmt.Sprintf("%-6.2f", val), !(val < 0)
}
func getTenantDeposit(j *jawaInfo, apt string) (string, bool) {
	val := j.Tenant[j.Rental[apt].TenantKey].Deposit
	return fmt.Sprintf("%-6.2f", val), !(val < 0)
}
func getTenantRentDueDate(j *jawaInfo, apt string) string {
	ten := j.Tenant[j.Rental[apt].TenantKey]
	return fmt.Sprintf("%d-%d-%d", ten.NextPaymentDue.Month(), ten.NextPaymentDue.Day(), ten.NextPaymentDue.Year())
}

/* obsolete
func getTenantPaidThrou(j *jawaInfo, apt string) string {
	ten := j.Tenant[j.Rental[apt].TenantKey]
	return fmt.Sprintf("%d-%d-%d", ten.PaidThrough.Month(), ten.PaidThrough.Day(), ten.PaidThrough.Year())
}
*/
func slamToInitTenantState(j *jawaInfo, apt string) {
	//fmt.Printf("Current Month: %s\n", time.Now().Month())
	rr := j.Rental[apt]
	tr := j.Tenant[rr.TenantKey]
	tr.Rent = map[string]payment{}
	tr.Water = map[string]payment{}
	tr.LatePenalty = map[string]payment{}
	tr.BouncePenalty = map[string]payment{}
	tr.RentOwed = 0
	tr.LatePOwed = 0
	tr.BounceOwed = 0
	tr.WaterOwed = 0
	tr.Deposit = 0
	j.Tenant[rr.TenantKey] = tr
}

func submitPayment(j *jawaInfo, apt string, date, next time.Time,
	rent, late, bounce, deposit float32, initPay bool) {

	ten := j.Tenant[j.Rental[apt].TenantKey]
	//aptUnit := j.Rental[apt]
	if !initPay {
		ten.RentOwed -= rent
	}
	if date.After(ten.NextPaymentDue.AddDate(0, 0, j.CI.IncurrAfter)) {
		ten.LatePOwed += j.CI.LateFee
	}
	ten.LatePOwed -= late
	ten.BounceOwed -= bounce
	ten.Deposit += deposit
	total := rent + late + bounce + deposit
	ten.Rent[getUniqueDateKey(date)] = payment{total, date}
	ten.NextPaymentDue = next
	j.Tenant[j.Rental[apt].TenantKey] = ten
}

func getPaymentAllocation(payRent, payLate, payBoun, payDepo gwu.TextBox) (rent, late, boun, depo float32, tryAgain bool) {
	//var rent, late, boun, depo float32
	tryAgain = false
	//var payDate, payRent, payLate, payBoun, payDepo gwu.TextBox
	n, err := fmt.Sscanf(payRent.Text(), "%f", &rent)
	if err != nil || n != 1 {
		fmt.Printf("Rent value is incorrect format, please try again\n") //TODO: notify
		tryAgain = true
	}
	n, err = fmt.Sscanf(payLate.Text(), "%f", &late)
	if err != nil || n != 1 {
		fmt.Printf("Late fee value is incorrect format, please try again\n")
		tryAgain = true
	}
	n, err = fmt.Sscanf(payBoun.Text(), "%f", &boun)
	if err != nil || n != 1 {
		fmt.Printf("Bounce fee value is incorrect format, please try again\n")
		tryAgain = true
	}
	n, err = fmt.Sscanf(payDepo.Text(), "%f", &depo)
	if err != nil || n != 1 {
		fmt.Printf("Deposit payment is incorrect format, please try again\n")
		tryAgain = true
	}
	return rent, late, boun, depo, tryAgain
}

func setBackgroudIfPos(cur gwu.Label, color string, positive bool) {
	if positive {
		cur.Style().SetBackground(color)
		return
	}
	cur.Style().SetBackground("")
}
func updateRecordPaymentPage(ji *jawaInfo, apt string, e gwu.Event,
	cur1, cur2, cur3, cur4, monthly, payDue, nextDueDate gwu.Label,
	cb, ays gwu.CheckBox, cbTable gwu.Table,
	payDate, payRent, payLate, payBoun, payDepo, totalSub gwu.TextBox) {

	payDate.SetText("mm-dd-yyyy")
	payRent.SetText("0.00")
	payLate.SetText("0.00")
	payBoun.SetText("0.00")
	payDepo.SetText("0.00")
	totalSub.SetText("")

	val, pos := getTenantRentDue(ji, apt)
	cur1.SetText(val)
	setBackgroudIfPos(cur1, gwu.ClrRed, pos)
	val, pos = getTenantLateDue(ji, apt)
	cur2.SetText(val)
	setBackgroudIfPos(cur2, gwu.ClrRed, pos)
	val, pos = getTenantBounceDue(ji, apt)
	cur3.SetText(val)
	setBackgroudIfPos(cur3, gwu.ClrRed, pos)
	val, pos = getTenantDeposit(ji, apt)
	cur4.SetText(val)
	setBackgroudIfPos(cur4, gwu.ClrRed, !pos)
	monthly.SetText(getRentalRent(ji, apt))
	payDue.SetText(getTenantRentDueDate(ji, apt))
	nextDueDate.SetText(getTenantRentDueDate(ji, apt))
	cb.Style().SetBackground("")
	cb.SetState(false)
	ays.SetState(false)
	cbTable.Remove(ays)

	e.MarkDirty(cur1)
	e.MarkDirty(cur2)
	e.MarkDirty(cur3)
	e.MarkDirty(cur4)
	e.MarkDirty(monthly)
	e.MarkDirty(payDue)
	e.MarkDirty(payDate)
	e.MarkDirty(nextDueDate)
	e.MarkDirty(cb)
	e.MarkDirty(ays)
	e.MarkDirty(cbTable)
	e.MarkDirty(payRent)
	e.MarkDirty(payLate)
	e.MarkDirty(payBoun)
	e.MarkDirty(payDepo)
	e.MarkDirty(totalSub)

}
func buildRecordPayments( /*t gwu.TabPanel,*/ ji *jawaInfo) gwu.Panel {
	var apt string
	var table, tableb, cbTable gwu.Table
	var cb, ays gwu.CheckBox
	var cur1, cur2, cur3, cur4, monthly, payDue gwu.Label
	var payDate, nextDueDate, payRent, payLate, payBoun, payDepo, totalSub gwu.TextBox

	c := gwu.NewPanel()
	tableA := gwu.NewTable()
	tableA.SetCellPadding(2)
	tableA.EnsureSize(1, 5)
	tableA.Add(gwu.NewLabel("Payment for apartment:"), 0, 0)
	list := getAptList(ji)
	apt = list[0]
	//aptlb := gwu.NewListBox([]string{"Top", "Middle", "Bottom"})
	aptlb := gwu.NewListBox(list)
	aptlb.Style().SetFullWidth()
	aptlb.AddEHandlerFunc(func(e gwu.Event) {
		apt = list[aptlb.SelectedIdx()]

		updateRecordPaymentPage(ji, apt, e, cur1, cur2, cur3, cur4,
			monthly, payDue, nextDueDate, cb, ays, cbTable,
			payDate, payRent, payLate, payBoun, payDepo, totalSub)
	}, gwu.ETypeChange)

	tableA.Add(aptlb, 0, 1)
	tableA.Add(gwu.NewLabel("............."), 0, 2)
	tableA.Add(gwu.NewLabel("Monthly Rent:"), 0, 3)
	monthly = gwu.NewLabel(getRentalRent(ji, apt))
	tableA.Add(monthly, 0, 4)
	c.Add(tableA)

	c.AddVSpace(15)

	table = gwu.NewTable()
	table.SetCellPadding(2)
	table.EnsureSize(7, 4)
	table.Add(gwu.NewLabel("Enter payment date (mm-dd-yyyy):"), 0, 0)
	table.Add(gwu.NewLabel("Enter new due date (mm-dd-yyyy):"), 1, 0)
	table.Add(gwu.NewLabel("Enter rent payment amount:"), 2, 0)
	table.Add(gwu.NewLabel("Enter Late fee payment amount:"), 3, 0)
	table.Add(gwu.NewLabel("Enter Bounce fee payment amount:"), 4, 0)
	table.Add(gwu.NewLabel("Enter Deposit amount:"), 5, 0)
	table.Add(gwu.NewLabel("Total Submitted:"), 6, 0)

	payDate = gwu.NewTextBox("mm-dd-yyyy")
	nextDueDate = gwu.NewTextBox(getTenantRentDueDate(ji, apt))
	payRent = gwu.NewTextBox("0.00")
	payLate = gwu.NewTextBox("0.00")
	payBoun = gwu.NewTextBox("0.00")
	payDepo = gwu.NewTextBox("0.00")
	totalSub = gwu.NewTextBox("")

	table.Add(payDate, 0, 1)
	table.Add(nextDueDate, 1, 1)
	table.Add(payRent, 2, 1)
	table.Add(payLate, 3, 1)
	table.Add(payBoun, 4, 1)
	table.Add(payDepo, 5, 1)
	table.Add(totalSub, 6, 1)

	val, pos := getTenantRentDue(ji, apt)
	cur1 = gwu.NewLabel(val)
	setBackgroudIfPos(cur1, gwu.ClrRed, pos)
	val, pos = getTenantLateDue(ji, apt)
	cur2 = gwu.NewLabel(val)
	setBackgroudIfPos(cur2, gwu.ClrRed, pos)
	val, pos = getTenantBounceDue(ji, apt)
	cur3 = gwu.NewLabel(val)
	setBackgroudIfPos(cur3, gwu.ClrRed, pos)
	val, pos = getTenantDeposit(ji, apt)
	cur4 = gwu.NewLabel(val)
	setBackgroudIfPos(cur4, gwu.ClrRed, !pos)

	table.Add(gwu.NewLabel("Current balance"), 1, 3)
	table.Add(cur1, 2, 3)
	table.Add(cur2, 3, 3)
	table.Add(cur3, 4, 3)
	table.Add(cur4, 5, 3)

	c.Add(table)

	cbTable = gwu.NewTable()
	cbTable.SetCellPadding(2)
	cbTable.EnsureSize(1, 2)

	cb = gwu.NewCheckBox("Initial payment?")
	cb.AddEHandlerFunc(func(e gwu.Event) {
		if cb.State() {
			cb.Style().SetBackground(gwu.ClrGreen)
			//ays.Style().SetFontStyle(gwu.ClrAqua)
			ays.Style().SetBackground(gwu.ClrAqua)
			cbTable.Add(ays, 0, 1)
		} else {
			cb.Style().SetBackground("")
			ays.SetState(false)
			cbTable.Remove(ays)
		}
		e.MarkDirty(cb)
		e.MarkDirty(ays)
		e.MarkDirty(cbTable)
	}, gwu.ETypeClick)
	cbTable.Add(cb, 0, 0)

	ays = gwu.NewCheckBox("ARE YOU SURE?")
	ays.AddEHandlerFunc(func(e gwu.Event) {
		if ays.State() {
			ays.Style().SetBackground(gwu.ClrGreen)

			slamToInitTenantState(ji, apt)

			val, pos := getTenantRentDue(ji, apt)
			cur1.SetText(val)
			setBackgroudIfPos(cur1, gwu.ClrRed, pos)
			val, pos = getTenantLateDue(ji, apt)
			cur2.SetText(val)
			setBackgroudIfPos(cur2, gwu.ClrRed, pos)
			val, pos = getTenantBounceDue(ji, apt)
			cur3.SetText(val)
			setBackgroudIfPos(cur3, gwu.ClrRed, pos)
			val, pos = getTenantDeposit(ji, apt)
			cur4.SetText(val)
			setBackgroudIfPos(cur4, gwu.ClrRed, !pos)
		} else {
			ays.SetState(true)
			//ays.Style().SetBackground("")
		}
		e.MarkDirty(ays)
		e.MarkDirty(cur1)
		e.MarkDirty(cur2)
		e.MarkDirty(cur3)
		e.MarkDirty(cur4)
	}, gwu.ETypeClick)

	//cbTable.Add(ays, 0, 1) // Only added by cb when checked!
	c.Add(cbTable)

	c.AddVSpace(10)
	hp := gwu.NewHorizontalPanel()
	hp.Style().SetFullWidth()
	hp.SetHAlign(gwu.HACenter)
	b := gwu.NewButton("Submit")
	b.AddEHandlerFunc(func(e gwu.Event) {
		var day, month, year int
		var dayN, monthN, yearN int
		var tryAgain = false

		n, err := fmt.Sscanf(payDate.Text(), "%d-%d-%d\n", &month, &day, &year)
		if err != nil || n != 3 {
			fmt.Printf("Payment date format is incorrect, please try again\n")
			tryAgain = true
			//TODO: need this not to be pop up for the like // notification
		}
		n, err = fmt.Sscanf(nextDueDate.Text(), "%d-%d-%d\n", &monthN, &dayN, &yearN)
		if err != nil || n != 3 {
			fmt.Printf("Next due date format is incorrect, please try again\n")
			tryAgain = true
			//TODO: need this not to be pop up for the like // notification
		}
		rent, late, bounce, deposit, tryAgain2 :=
			getPaymentAllocation(payRent, payLate, payBoun, payDepo)
		if !tryAgain && !tryAgain2 {
			date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
			//fmt.Printf("Month: %d, Day: %d, Year: %d\n", month, day, year)
			//fmt.Printf("Update in progress at time: %v\n", date)
			newDueDate := time.Date(yearN, time.Month(monthN), dayN, 0, 0, 0, 0, time.Local)
			submitPayment(ji, apt, date, newDueDate, rent, late, bounce, deposit, ays.State())

			updateRecordPaymentPage(ji, apt, e, cur1, cur2, cur3, cur4,
				monthly, payDue, nextDueDate, cb, ays, cbTable,
				payDate, payRent, payLate, payBoun, payDepo, totalSub)

			/*
				payDate.SetText("mm-dd-yyyy")

				val, pos := getTenantRentDue(ji, apt)
				cur1.SetText(val)
				setBackgroudIfPos(cur1, gwu.ClrRed, pos)
				val, pos = getTenantLateDue(ji, apt)
				cur2.SetText(val)
				setBackgroudIfPos(cur2, gwu.ClrRed, pos)
				val, pos = getTenantBounceDue(ji, apt)
				cur3.SetText(val)
				setBackgroudIfPos(cur3, gwu.ClrRed, pos)
				val, pos = getTenantDeposit(ji, apt)
				cur4.SetText(val)
				setBackgroudIfPos(cur4, gwu.ClrRed, !pos)

				nextDueDate.SetText(getTenantRentDueDate(ji, apt)) // Why is this not working
				* /

				cb.SetState(false)
				cb.Style().SetBackground("")
				ays.SetState(false)
				cb.Style().SetBackground("")
				cbTable.Remove(ays)
				e.MarkDirty(payDate)
				e.MarkDirty(nextDueDate)
				e.MarkDirty(cur1)
				e.MarkDirty(cur2)
				e.MarkDirty(cur3)
				e.MarkDirty(cur4)
				e.MarkDirty(cb)
				e.MarkDirty(ays)
				e.MarkDirty(cbTable)
			*/
		}
	}, gwu.ETypeClick)
	hp.Add(b)

	//c.AddVSpace(15)

	bt := gwu.NewButton("Total up")
	bt.AddEHandlerFunc(func(e gwu.Event) {
		rent, late, boun, depo, tryAgain :=
			getPaymentAllocation(payRent, payLate, payBoun, payDepo)
		if !tryAgain {
			total := rent + late + boun + depo
			totalSub.SetText(fmt.Sprintf("%-6.2f", total))
			e.MarkDirty(totalSub)
		}
	}, gwu.ETypeClick)
	hp.Add(bt)
	c.Add(hp)
	c.AddVSpace(10)

	tableb = gwu.NewTable()
	tableb.SetCellPadding(2)
	tableb.EnsureSize(4, 2)

	tableb.Add(gwu.NewLabel("Date:"), 1, 0)
	tableb.Add(gwu.NewLabel("Next payment due:"), 2, 0)

	date := gwu.NewLabel(fmt.Sprintf("%d-%d-%d", time.Now().Month(), time.Now().Day(), time.Now().Year()))
	payDue = gwu.NewLabel(getTenantRentDueDate(ji, apt))

	tableb.Add(date, 1, 1)
	tableb.Add(payDue, 2, 1)

	c.Add(tableb)
	return c
}
