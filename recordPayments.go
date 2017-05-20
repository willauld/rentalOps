package main

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/icza/gowut/gwu"
	//"github.com/willauld/gowut/gwu"
)

/*
func getAptListOld(ji *jawaInfo) []string {
	return *ji.CI.Apartments
}
*/
func strIndex(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

//TODO:GetKeys getAptList and GetTenantList should be moved to Main

//GetKeys returns a []string containing the map keys (keys must be type string)
func GetKeys(m interface{}) []string {
	//TODO: NEED CHECKING THAT KEYS ARE STRINGS
	temp := reflect.ValueOf(m).MapKeys()

	keys := make([]string, len(temp))
	for i, elem := range temp {
		v := elem.String()
		keys[i] = v
	}
	sort.Strings(keys)
	return keys
}
func getAptList(ji *jawaInfo) []string {
	if len(ji.Rental) < 1 {
		fmt.Printf("There are no apartments yet defined\n")
		keys := make([]string, 1)
		keys[0] = "undefined"
		return keys
	}
	/*
		keys := make([]string, len(ji.Rental))

		i := 0
		for k := range ji.Rental {
			keys[i] = k
			i++
		}
		sort.Strings(keys)
		return keys
	*/
	return GetKeys(ji.Rental)
}
func getTenantList(ji *jawaInfo) []string {
	if len(ji.Rental) < 1 {
		fmt.Printf("There are no tenants yet defined\n")
		keys := make([]string, 1)
		keys[0] = "undefined"
		return keys
	}
	/*
		keys := make([]string, len(ji.Tenant))

		i := 0
		for k := range ji.Tenant {
			keys[i] = k
			i++
		}
		sort.Strings(keys)
		return keys
	*/
	return GetKeys(ji.Tenant)
}
func makeTenantKey(nameStr string) string {
	return strings.Replace(nameStr, " ", "", -1)
}
func getRentalRent(j *jawaInfo, apt string) string {
	return fmt.Sprintf("%-6.2f", j.Rental[apt].Rent)
}
func getDayOfMonthDue(j *jawaInfo, apt string) string {
	return fmt.Sprintf("%d", j.Rental[apt].DueDay)
}
func getRentalDeposit(j *jawaInfo, apt string) string {
	return fmt.Sprintf("%-6.2f", j.Rental[apt].Deposit)
}
func getRentalTenant(j *jawaInfo, apt string) string {
	return j.Rental[apt].Tenant
}
func getTenantKey(j *jawaInfo, apt string) string {
	return j.Rental[apt].TenantKey
}
func getTenantRentOwed(j *jawaInfo, tenkey string) (string, bool) {
	val := j.Tenant[tenkey].RentOwed
	return fmt.Sprintf("%-6.2f", val), !(val < 0)
}
func getTenantLateOwed(j *jawaInfo, tenkey string) (string, bool) {
	val := j.Tenant[tenkey].LatePOwed
	return fmt.Sprintf("%-6.2f", val), !(val < 0)
}
func getTenantBounceOwed(j *jawaInfo, tenkey string) (string, bool) {
	val := j.Tenant[tenkey].BounceOwed
	return fmt.Sprintf("%-6.2f", val), !(val < 0)
}
func getTenantDepositOwed(j *jawaInfo, tenkey string) (string, bool) {
	val := j.Tenant[tenkey].DepositOwed
	return fmt.Sprintf("%-6.2f", val), !(val < 0)
}
func getTenantWaterOwed(j *jawaInfo, tenkey string) (string, bool) {
	val := j.Tenant[tenkey].WaterOwed
	return fmt.Sprintf("%-6.2f", val), !(val < 0)
}
func getTenantRentDueDate(j *jawaInfo, tenkey string) string {
	ten := j.Tenant[tenkey]
	return fmt.Sprintf("%d-%d-%d", ten.NextPaymentDue.Month(), ten.NextPaymentDue.Day(), ten.NextPaymentDue.Year())
}
func getTenantRentChargedThru(j *jawaInfo, tenkey string) string {
	ten := j.Tenant[tenkey]
	return fmt.Sprintf("%d-%d-%d", ten.RentChargedThru.Month(), ten.RentChargedThru.Day(), ten.RentChargedThru.Year())
}
func getTenantNumPayments(j *jawaInfo, tenkey string) string {
	fmt.Printf("tenkey: [%s]\n", tenkey)
	ten := j.Tenant[tenkey]
	return fmt.Sprintf("%d", len(ten.Payment))
}

/* obsolete
func getTenantPaidThrou(j *jawaInfo, apt string) string {
	ten := j.Tenant[j.Rental[apt].TenantKey]
	return fmt.Sprintf("%d-%d-%d", ten.PaidThrough.Month(), ten.PaidThrough.Day(), ten.PaidThrough.Year())
}
*/
func slamToInitTenantState(j *jawaInfo, apt string) {
	//fmt.Printf("Current Month: %s\n", timeNowRental().Month())
	rr := j.Rental[apt]
	tr := j.Tenant[rr.TenantKey]
	tr.Payment = map[string]payment{}
	tr.RentOwed = 0
	tr.LatePOwed = 0
	tr.BounceOwed = 0
	tr.WaterOwed = 0
	tr.DepositOwed = 0
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
	ten.DepositOwed -= deposit
	total := rent + late + bounce + deposit
	if ten.Payment == nil {
		ten.Payment = map[string]payment{}
	}
	key, err := getUniqueKey(ten.Payment, date)
	if err != nil {
		fmt.Printf("getUniqueKey() failed: %s", err)
	}
	ten.Payment[key] = payment{total, rent, late, bounce, 0, deposit, date}
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

func setTextFontIfPos(cur gwu.Label, color string, positive bool) {
	if positive {
		cur.Style().SetColor(color)
		//cur.Style().SetBackground(color)
		return
	}
	cur.Style().SetBackground("")
}
func updateRecordPaymentPage(ji *jawaInfo, apt string, e gwu.Event,
	cur1, cur2, cur3, cur4, monthly, rentalDeposit, payDue,
	nextDueDate, Tenant gwu.Label,
	cb, ays gwu.CheckBox, cbTable, tableb gwu.Table,
	payDate, payRent, payLate, payBoun, payDepo, totalSub gwu.TextBox) {

	/*
			payDate.SetText("mm-dd-yyyy")
				payRent.SetText("0.00")
				payLate.SetText("0.00")
				payBoun.SetText("0.00")
				payDepo.SetText("0.00")
				totalSub.SetText("")
		payDue.SetText(getTenantRentDueDate(ji, getTenantKey(ji, apt)))
			nextDueDate.SetText(getTenantRentDueDate(ji, apt))
	*/

	val, pos := getTenantRentOwed(ji, getTenantKey(ji, apt))
	cur1.SetText(val)
	setTextFontIfPos(cur1, gwu.ClrRed, pos)
	val, pos = getTenantLateOwed(ji, getTenantKey(ji, apt))
	cur2.SetText(val)
	setTextFontIfPos(cur2, gwu.ClrRed, pos)
	val, pos = getTenantBounceOwed(ji, getTenantKey(ji, apt))
	cur3.SetText(val)
	setTextFontIfPos(cur3, gwu.ClrRed, pos)
	val, pos = getTenantDepositOwed(ji, getTenantKey(ji, apt))
	cur4.SetText(val)
	setTextFontIfPos(cur4, gwu.ClrRed, pos)

	monthly.SetText(getRentalRent(ji, apt))
	rentalDeposit.SetText(getRentalDeposit(ji, apt))
	Tenant.SetText(getRentalTenant(ji, apt))

	payDue.SetText(getTenantRentDueDate(ji, getTenantKey(ji, apt)))
	fmt.Printf("payDue: %+v\n", payDue.Text())
	/*
		nextDueDate.SetText(getTenantRentDueDate(ji, apt))
	*/
	cb.Style().SetBackground("")
	cb.SetState(false)
	ays.SetState(false)
	cbTable.Remove(ays)

	if e != nil {
		e.MarkDirty(cur1)
		e.MarkDirty(cur2)
		e.MarkDirty(cur3)
		e.MarkDirty(cur4)
		e.MarkDirty(monthly)
		e.MarkDirty(rentalDeposit)
		e.MarkDirty(Tenant)
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
		e.MarkDirty(tableb)
		Notify("tableb marked dirty ", e)
	}
}
func buildRecordPayments(ji *jawaInfo) (gwu.Panel, gwu.TextBox) {
	var apt string //This is the KEY for the active apartment
	var aptlb gwu.ListBox
	var table, tableA, tableb, cbTable gwu.Table
	var cb, ays gwu.CheckBox
	var aptlbHandler func(e gwu.Event)
	var cur1, cur2, cur3, cur4, monthly, rentalDeposit, payDue, rentalTenant gwu.Label
	var payDate, nextDueDate, payRent, payLate, payBoun, payDepo, totalSub gwu.TextBox

	c := gwu.NewPanel()
	stb := gwu.NewTextBox("")
	stb.Style().SetWidthPx(1).SetHeightPx(1)
	stb.AddEHandlerFunc(func(e gwu.Event) {
		//meList := getKeyList(j.Rental) //TODO: ***WANT TO MOVE TO THIS not getAptList()
		meList := getAptList(ji)
		tableA.Remove(aptlb)
		err := UpdateListBox(&aptlb, &apt, meList, aptlbHandler)
		if err != nil {
			fmt.Printf("EditTenant update ListBox failed: %v\n", err)
			return
		}
		tableA.Add(aptlb, 0, 1)

		rent, late, boun, depo, tryAgain :=
			getPaymentAllocation(payRent, payLate, payBoun, payDepo)
		if !tryAgain {
			total := rent + late + boun + depo
			totalSub.SetText(fmt.Sprintf("%-6.2f", total))
			e.MarkDirty(totalSub)
		}
		updateRecordPaymentPage(ji, apt, e, cur1, cur2, cur3, cur4,
			monthly, rentalDeposit, payDue, nextDueDate, rentalTenant,
			cb, ays, cbTable, tableb, payDate, payRent, payLate,
			payBoun, payDepo, totalSub)
		e.MarkDirty(aptlb)
		e.MarkDirty(tableA)

		Notify("Focus is on Record Payments Tab", e)
	}, gwu.ETypeFocus)
	c.Add(stb)

	tableA = gwu.NewTable()
	tableA.SetCellPadding(2)
	tableA.EnsureSize(2, 8)
	tableA.Add(gwu.NewLabel("Payment for apartment:"), 0, 0)
	list := getAptList(ji)
	apt = list[0]
	aptlb = gwu.NewListBox(list)
	aptlb.Style().SetFullWidth()
	aptlbHandler = func(e gwu.Event) {
		list := getAptList(ji)
		indx := aptlb.SelectedIdx()
		if indx < 0 {
			indx = 0
		}
		apt = list[indx]

		updateRecordPaymentPage(ji, apt, e, cur1, cur2, cur3, cur4,
			monthly, rentalDeposit, payDue, nextDueDate, rentalTenant,
			cb, ays, cbTable, tableb, payDate, payRent, payLate,
			payBoun, payDepo, totalSub)
	}
	aptlb.AddEHandlerFunc(aptlbHandler, gwu.ETypeChange)

	tableA.Add(aptlb, 0, 1)
	tableA.Add(gwu.NewLabel("..."), 0, 2)
	tableA.Add(gwu.NewLabel("Monthly Rent:"), 0, 3)
	monthly = gwu.NewLabel(getRentalRent(ji, apt))
	tableA.Add(monthly, 0, 4)
	tableA.Add(gwu.NewLabel("..."), 0, 5)
	tableA.Add(gwu.NewLabel("Deposit:"), 0, 6)
	rentalDeposit = gwu.NewLabel(getRentalDeposit(ji, apt))
	tableA.Add(rentalDeposit, 0, 7)
	tableA.Add(gwu.NewLabel("Tenant:"), 1, 3)
	rentalTenant = gwu.NewLabel(getRentalTenant(ji, apt))
	tableA.Add(rentalTenant, 1, 4)
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
	totalSub = gwu.NewTextBox("0.00")

	table.Add(payDate, 0, 1)
	table.Add(nextDueDate, 1, 1)
	table.Add(payRent, 2, 1)
	table.Add(payLate, 3, 1)
	table.Add(payBoun, 4, 1)
	table.Add(payDepo, 5, 1)
	table.Add(totalSub, 6, 1)

	val, pos := getTenantRentOwed(ji, apt)
	cur1 = gwu.NewLabel(val)
	setTextFontIfPos(cur1, gwu.ClrRed, pos)
	val, pos = getTenantLateOwed(ji, apt)
	cur2 = gwu.NewLabel(val)
	setTextFontIfPos(cur2, gwu.ClrRed, pos)
	val, pos = getTenantBounceOwed(ji, apt)
	cur3 = gwu.NewLabel(val)
	setTextFontIfPos(cur3, gwu.ClrRed, pos)
	val, pos = getTenantDepositOwed(ji, apt)
	cur4 = gwu.NewLabel(val)
	setTextFontIfPos(cur4, gwu.ClrRed, pos)

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

			val, pos := getTenantRentOwed(ji, apt)
			cur1.SetText(val)
			setTextFontIfPos(cur1, gwu.ClrRed, pos)
			val, pos = getTenantLateOwed(ji, apt)
			cur2.SetText(val)
			setTextFontIfPos(cur2, gwu.ClrRed, pos)
			val, pos = getTenantBounceOwed(ji, apt)
			cur3.SetText(val)
			setTextFontIfPos(cur3, gwu.ClrRed, pos)
			val, pos = getTenantDepositOwed(ji, apt)
			cur4.SetText(val)
			setTextFontIfPos(cur4, gwu.ClrRed, pos)
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

			payDate.SetText("mm-dd-yyyy")
			payRent.SetText("0.00")
			payLate.SetText("0.00")
			payBoun.SetText("0.00")
			payDepo.SetText("0.00")
			totalSub.SetText("0.00")
			nextDueDate.SetText(getTenantRentDueDate(ji, getTenantKey(ji, apt)))
			payDue.SetText(getTenantRentDueDate(ji, getTenantKey(ji, apt)))

			updateRecordPaymentPage(ji, apt, e, cur1, cur2, cur3, cur4,
				monthly, rentalDeposit, payDue, nextDueDate, rentalTenant,
				cb, ays, cbTable, tableb, payDate, payRent, payLate,
				payBoun, payDepo, totalSub)
		}

		e.MarkDirty(payDate, payRent, payLate, payBoun, payDepo)
		e.MarkDirty(totalSub, payDue, nextDueDate)
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
		updateRecordPaymentPage(ji, apt, e, cur1, cur2, cur3, cur4,
			monthly, rentalDeposit, payDue, nextDueDate, rentalTenant,
			cb, ays, cbTable, tableb, payDate, payRent, payLate,
			payBoun, payDepo, totalSub)
	}, gwu.ETypeClick)
	hp.Add(bt)
	c.Add(hp)
	c.AddVSpace(10)

	tableb = gwu.NewTable()
	tableb.SetCellPadding(2)
	tableb.EnsureSize(4, 2)

	tableb.Add(gwu.NewLabel("Date:"), 1, 0)
	tableb.Add(gwu.NewLabel("Next payment due:"), 2, 0)

	date := gwu.NewLabel(fmt.Sprintf("%d-%d-%d", timeNowRental().Month(), timeNowRental().Day(), timeNowRental().Year()))
	payDue = gwu.NewLabel(getTenantRentDueDate(ji, getTenantKey(ji, apt)))

	tableb.Add(date, 1, 1)
	tableb.Add(payDue, 2, 1)

	c.Add(tableb)
	updateRecordPaymentPage(ji, apt, nil, cur1, cur2, cur3, cur4,
		monthly, rentalDeposit, payDue, nextDueDate, rentalTenant,
		cb, ays, cbTable, tableb, payDate, payRent, payLate,
		payBoun, payDepo, totalSub)
	return c, stb
}
