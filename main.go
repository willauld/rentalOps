package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/jinzhu/now"
	"github.com/spf13/pflag"
)

var version = struct {
	major int
	minor int
}{0, 2}

type commonInfo struct {
	OfficeManager       string
	OfficeStreet        string
	OfficeCity          string
	OfficeState         string
	OfficeZip           string
	CorrespondenceTitle string
	LateFee             float32
	IncurrAfter         int
	BounceFee           float32
	Apartments          *[]string
}
type rentalRecord struct {
	Apartment string
	Tenant    string
	TenantKey string
	Rent      float32
	Deposit   float32
	DueDay    int
	Street    string
	City      string
	State     string
	Zip       string
}
type payment struct {
	Amount  float32
	Rent    float32
	Late    float32
	Bounce  float32
	Water   float32
	Deposit float32
	Date    time.Time
}
type renterRecord struct {
	Apartment   string
	Name        string
	Payment     map[string]payment
	RentOwed    float32
	BounceOwed  float32
	LatePOwed   float32
	WaterOwed   float32
	DepositOwed float32

	NextPaymentDue  time.Time
	RentChargedThru time.Time
}
type jawaInfo struct {
	CI     commonInfo
	Rental map[string]rentalRecord
	Tenant map[string]renterRecord
}

func displayCurrentAptList(a *[]string) {
	fmt.Printf("==========================================\n")
	fmt.Printf("Current List of Apartments:\n")
	if a != nil {
		for j, apt := range *a {
			fmt.Printf("Unit[%d]:  %s\n", j+1, apt)
		}
	}
	fmt.Printf("==========================================\n")
}
func displayRental2(r rentalRecord) { //TODO delete if dup
	fmt.Printf("Apartment: %s\n", r.Apartment)
	fmt.Printf("Tenant: %s\n", r.Tenant)
	fmt.Printf("TenantKey: %s\n", r.TenantKey)
	fmt.Printf("Address:\n\t%s\n\t%s, %s %s\n", r.Street, r.City, r.State, r.Zip)
	fmt.Printf("\n")
	fmt.Printf("Rent: %-6.2f\n", r.Rent)
	fmt.Printf("Due Day: %d\n", r.DueDay)
}
func isKnownApt(apt string, i *jawaInfo) bool {
	for _, v := range *i.CI.Apartments {
		if apt == v {
			return true
		}
	}
	return false
}
func editRental(i *jawaInfo, apt string) {
	var tenantKey string
	var tr renterRecord
	//apt := (*i.CI.Apartments)[indx-1]
	if !isKnownApt(apt, i) {
		fmt.Printf("Unknown Apartment name: %s\n", apt)
		return
	}
	if i.Rental == nil {
		i.Rental = map[string]rentalRecord{}
	}
	if i.Tenant == nil {
		i.Tenant = map[string]renterRecord{}
	}
	displayRental(i, apt, true)
	rr, ok := i.Rental[apt]
	if !ok {
		rr = rentalRecord{}
		rr.Apartment = apt
	}
	nameStr := getInputQuotedStr("Enter tenant's name double quoted: ")
	if nameStr != "" {
		tenantKey = strings.Replace(nameStr, " ", "", -1)
		tr, ok = i.Tenant[tenantKey]
		if !ok {
			fmt.Printf("Tenant appears to be NEW\n")
			tr = renterRecord{}
		}
		tr.NextPaymentDue = now.New(timeNowRental().AddDate(0, 1, 0)).BeginningOfMonth()
		rr.TenantKey = tenantKey
		tr.Apartment = apt
		tr.Name = nameStr
		rr.Tenant = nameStr
	}
	rent := getInputFloat("Enter Monthly Rent: ")
	if rent != 0 {
		rr.Rent = rent
	}
	dueDay := getInputInt("Enter day of the month the rent will be due: ")
	if dueDay != 0 {
		rr.DueDay = dueDay
	}
	str := getInputQuotedStr("Enter street address double quoted: ")
	if str != "" {
		rr.Street = str
	}
	str = getInputStr("Enter City: ")
	if str != "" {
		rr.City = str
	}
	str = getInputStr("Enter State: ")
	if str != "" {
		rr.State = str
	}
	str = getInputStr("Enter Zip: ")
	if str != "" {
		rr.Zip = str
	}
	i.Rental[apt] = rr
	i.Tenant[tenantKey] = tr
	if nameStr != "" {
		indx := getInputInt("Enter an initial Payment?\n0)No\n1)Yes\n\tSelect the number:\n")
		if indx == 1 {
			recordPayments(i.Rental[apt], i, true)
		}
	}
	fmt.Printf("=================================\n")
	fmt.Printf("Updated Record is now as follows:\n")
	displayRental(i, apt, true)
}

func inputAndDeleteString(display string, a *[]string) *[]string {
	var input string

	displayCurrentAptList(a)
	fmt.Printf("%s", display)
	fmt.Scanf("%s\n", &input)
	if a == nil {
		a = new([]string)
		fmt.Printf("There are no apartements to remove %s from\n", input)
	}
	*a = remove(*a, input)
	//fmt.Printf("new appartments is\n%+v\n", a)
	return a
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func inputAndAddString(display string, a *[]string) *[]string {
	var input string
	displayCurrentAptList(a)
	fmt.Printf("%s", display)
	fmt.Scanf("%s\n", &input)
	if a == nil {
		a = new([]string)
	}
	*a = append(*a, input)
	//fmt.Printf("new appartments is\n%+v\n", a)
	return a
}

func getInputFloat(display string) float32 {
	var input float32
	fmt.Printf("%s", display)
	fmt.Scanf("%f\n", &input)
	return input
}

func getInputInt(display string) int {
	var input int
	fmt.Printf("%s", display)
	fmt.Scanf("%d\n", &input)
	return input
}

func getInputPaidThrough(datePaid time.Time) time.Time {
	var date time.Time
	//var str string
	guess1 := now.New(datePaid).EndOfMonth()
	guess1Str := guess1.Format("Jan 2")
	guess2 := now.New(guess1.AddDate(0, 0, 1)).EndOfMonth()
	guess2Str := guess2.Format("Jan 2")
	str := fmt.Sprintf("This payment will pay through to:\n1)%v\n2)%v\n0)Niether\n\tSelect an index: ", guess1Str, guess2Str)
escape:
	for {
		index := getInputInt(str)
		switch index {
		case 0:
			date = getInputDate()
			break escape
		case 1:
			date = guess1
			break escape
		case 2:
			date = guess2
			break escape
		default:
			fmt.Printf("incorrect choice of index, please try again\n")
		}
	}
	return date
}

func getInputDate() time.Time {
	var day, month, year int
	fmt.Printf("%s", "Enter payment date as mm-dd-yy (no spaces): ")
	fmt.Scanf("%d-%d-%d\n", &month, &day, &year)
	return time.Date(2000+year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}
func getInputStr(display string) string {
	var input string
	fmt.Printf("%s", display)
	fmt.Scanf("%s\n", &input)
	//fmt.Printf("getInputStr() found [%s]\n", input)
	return input
}

func getInputQuotedStr(display string) string {
	var input string
	fmt.Printf("%s", display)
	fmt.Scanf("%q\n", &input)
	//fmt.Printf("getInputStr() found [%s]\n", input)
	return input
}

func displayCommon(i *jawaInfo) {
	fmt.Printf("==========================================\n")
	fmt.Printf("Common Rental Information\n")
	fmt.Printf("------------------------------------------\n")
	if i == nil {
		return
	}
	fmt.Printf("Late Fee:    %v\n", i.CI.LateFee)
	fmt.Printf("Incur after: %d days\n", i.CI.IncurrAfter)
	fmt.Printf("Bounce Fee:  %v\n", i.CI.BounceFee)
	if i.CI.Apartments != nil {
		for j, apt := range *i.CI.Apartments {
			fmt.Printf("Unit[%d]:  %s\n", j+1, apt)
		}
	}
	fmt.Printf("==========================================\n")
}

func displayTenant(tr renterRecord, withHistory bool) {
	fmt.Printf("---------------------\n")
	fmt.Printf("Tenant (renterRecord)\n")
	fmt.Printf("---------------------\n")
	fmt.Printf("Appartment:       %s\n", tr.Apartment)
	fmt.Printf("Tenant:           %s\n", tr.Name)
	fmt.Printf("Rent owed:        %-6.2f\n", tr.RentOwed)
	fmt.Printf("Late fee owed:    %-6.2f\n", tr.LatePOwed)
	fmt.Printf("Bounce fee owed:  %-6.2f\n", tr.BounceOwed)
	fmt.Printf("Water owed:       %-6.2f\n", tr.WaterOwed)
	fmt.Printf("DepositOwed paid: %-6.2f\n", tr.DepositOwed)
	fmt.Printf("Next Payment due: %v\n", tr.NextPaymentDue)
	if withHistory {
		fmt.Printf("\t%d Rent Payments\n", len(tr.Payment))
		for k, v := range tr.Payment {
			fmt.Printf("\t\tPayment for %v:\n", k)
			fmt.Printf("\t\t\tAmount:  %-6.2f\n", v.Amount)
			fmt.Printf("\t\t\tPaid on: %v\n", v.Date)
		}
	}
}

func displayRental(i *jawaInfo, apt string, displayRenter bool) {
	if i == nil {
		return
	}
	rr, ok := i.Rental[apt]
	if ok {
		fmt.Printf("Apartment: %s\n", rr.Apartment)
		fmt.Printf("Tenant:    %s\n", rr.Tenant)
		fmt.Printf("TenantKey: %s\n", rr.TenantKey)
		fmt.Printf("Street:    %s\n", rr.Street)
		fmt.Printf("City:      %s\n", rr.City)
		fmt.Printf("State:     %s\n", rr.State)
		fmt.Printf("Zip:       %s\n", rr.Zip)
		fmt.Println()
		fmt.Printf("Rent:      %v\n", rr.Rent)
		fmt.Printf("Due by day %d\n", rr.DueDay)
		if displayRenter {
			tr, ok := i.Tenant[rr.TenantKey]
			if ok {
				fmt.Println()
				displayTenant(tr, true)
			} else {
				fmt.Printf("\tNo Current Tenant\n")
			}
		}
	} else {
		fmt.Printf("Record for %s does not yet exist.\n", apt)
	}
}

func printIfNonZero(amt float32, name string) {
	if amt > 0 {
		fmt.Printf("\t\t%-6.2f\t%s\n", amt, name)
	}
}

func createLetter(r rentalRecord, ji *jawaInfo) {
	fmt.Printf("\n")
	fmt.Printf("\n")
	fmt.Printf("%s\n", r.Tenant)
	fmt.Printf("%s\n", r.Street)
	fmt.Printf("%s, %s %s\n", r.City, r.State, r.Zip)

	tr := ji.Tenant[r.TenantKey]
	rm := tr.NextPaymentDue.Month()
	sum := r.Rent + tr.RentOwed + tr.LatePOwed + tr.BounceOwed + tr.WaterOwed
	t := timeNowRental()
	fmt.Printf("\n")
	fmt.Printf("%s %d, %d\n", t.Month(), t.Day(), t.Year())
	fmt.Printf("\n")
	fmt.Printf("Dear %s,\n", strings.Fields(r.Tenant)[0])
	fmt.Printf("\n")
	fmt.Printf("Your %s payment is based on the following details:\n", rm)
	fmt.Printf("\n")
	printIfNonZero(r.Rent, "Rent")
	printIfNonZero(tr.RentOwed, "Past due Rent")
	printIfNonZero(tr.WaterOwed, "Water payment")
	printIfNonZero(tr.LatePOwed, "Late fee")
	printIfNonZero(tr.BounceOwed, "Bounce fee")
	fmt.Printf("\t\t========\n")
	fmt.Printf("\t\t%-6.2f\t%s\n", sum, "Balance due")
	fmt.Printf("\n")
	fmt.Printf("Thank you for making your payment promtly,\n")
	fmt.Printf("\n")
	fmt.Printf("Yuli Auld\n")
	fmt.Printf("Apartment Manager\n")
	fmt.Printf("\n")
	fmt.Printf("\n")
}

func getUniqueDateKey(date time.Time) string {
	nano := timeNowRental().Nanosecond()
	return fmt.Sprintf("%s_%d_%d", date.Month().String(), date.Year(), nano)
}
func getUniqueKey(tr map[string]payment, date time.Time) (string, error) {
	if tr == nil {
		return "", fmt.Errorf("tenant.Payment map is nil")
	}
	i := 1
	for {
		key := fmt.Sprintf("%d_%d_%d_%d", date.Month() /*.String()*/, date.Day(), date.Year(), i)
		_, ok := tr[key]
		if !ok {
			return key, nil
		}
		i++
	}
}

func recordPayments(r rentalRecord, ji *jawaInfo, initialPayment bool) {
	//fmt.Printf("Current Month: %s\n", timeNowRental().Month())
	tr := ji.Tenant[r.TenantKey]
	rr := ji.Rental[tr.Apartment] // Why am i not just using "r"
	if tr.Payment == nil || initialPayment {
		tr.Payment = map[string]payment{}
		/*
			}
			if initialPayment {
		*/
		tr.RentOwed = 0
		tr.LatePOwed = 0
		tr.BounceOwed = 0
		tr.WaterOwed = 0
		tr.DepositOwed = 0
	}

	date := getInputDate()
	nextDue := tr.NextPaymentDue
	passedDue := nextDue.AddDate(0, 0, ji.CI.IncurrAfter)
	if len(tr.Payment) > 0 { // This is not the initial payment
		if date.After(passedDue) {
			tr.LatePOwed = tr.LatePOwed + ji.CI.LateFee
		}
	}
around:
	for {
		//printTenant(tr, rr.Rent, rr.DueDay) //Todo: dedub this function
		displayTenant(tr, true)
		fmt.Printf("Monthly Rent: %-6.2f, due by day %d\n", rr.Rent, rr.DueDay)
		pr := payment{}
		pr.Date = date
		//key := fmt.Sprintf("%s_%d", pr.Date.Month().String(), pr.Date.Year())
		//fmt.Printf("payment KEY is: %s\n", key)
		index := getInputInt("Apply payment towards?\n1)Rent\n2)Water\n3)Late fee\n4)Bounce fee\n5)Deposit\n0)Skip/done update\n\tSelect an index: ")
		switch index {
		case 0:
			ji.Tenant[r.TenantKey] = tr
			fmt.Printf("Finished with this tenant\n")
			return
		case 1:
			pr.Amount = getInputFloat("Enter Rent Payed: ")
			if len(tr.Payment) > 0 {
				tr.RentOwed = tr.RentOwed + rr.Rent - pr.Amount
			}
			pt := getInputPaidThrough(date)
			tr.NextPaymentDue = pt.AddDate(0, 0, 1)
			key, err := getUniqueKey(tr.Payment, date)
			if err != nil {
				fmt.Printf("getUniqueKey() failed: %s", err)
			}
			fmt.Printf("payment KEY is: %s\n", key)
			tr.Payment[key] = pr
			fmt.Printf("Rent[] has %d elements now\n", len(tr.Payment))
		case 2:
			pr.Amount = getInputFloat("Enter Water Payed: ")
			// TODO, adjust the WaterOwed???
			// TODO, All things water will wait for later
			key, err := getUniqueKey(tr.Payment, date)
			if err != nil {
				fmt.Printf("getUniqueKey() failed: %s", err)
			}
			tr.Payment[key] = pr
		case 3:
			pr.Amount = getInputFloat("Enter Late Fee Payed: ")
			tr.LatePOwed = tr.LatePOwed - pr.Amount
			key, err := getUniqueKey(tr.Payment, date)
			if err != nil {
				fmt.Printf("getUniqueKey() failed: %s", err)
			}
			tr.Payment[key] = pr
		case 4:
			pr.Amount = getInputFloat("Enter Bounce fee Payed: ")
			tr.BounceOwed = tr.BounceOwed - pr.Amount
			key, err := getUniqueKey(tr.Payment, date)
			if err != nil {
				fmt.Printf("getUniqueKey() failed: %s", err)
			}
			tr.Payment[key] = pr
		case 5:
			tr.DepositOwed = getInputFloat("Enter Deposit Payed: ")
		default:
			fmt.Printf("Incorrect entry, please try again\n")
			continue around
		}
	}
}

func process(i *jawaInfo, payments bool, firstPayment bool) {
	if firstPayment {
		payments = true
	}
	if len(i.Tenant) < 1 {
		fmt.Printf("No Tenants to process\n")
		return
	}
	if len(i.Rental) < 1 {
		fmt.Printf("No Rentals to process\n")
		return
	}
	for {
		fmt.Printf("Current Tenants:\n")
		for k, v := range i.Rental {
			//fmt.Printf("stuffi %v\n", v)
			//fmt.Printf("Apt: %s, Tenant: %s\n", k, i.Rental[v.Apartment].Tenant)
			fmt.Printf("\tApt: %s, Tenant: %s\n", k, v.Tenant)
		}
		apt := getInputStr("Enter the Apt to process: ")
		if apt == "" {
			fmt.Printf("No apartment chosen, exiting")
			return
		}
		rr, ok := i.Rental[apt]
		if !ok {
			fmt.Printf("Somethings wrong, things are not OK\n")
		}
		if payments {
			recordPayments(rr, i, firstPayment)
			// First payments should be uncommon and is distructive it
			// should not loop like other operations
			if firstPayment {
				return
			}
		} else {
			createLetter(rr, i)
		}
	}
}

func editInfo(i *jawaInfo) {
Pizza:
	for {
		displayCommon(i)
		index := getInputInt("Which field do you want to set?\n1)Late fee\n2)Incur after\n3)Bounce\n4)Apartment list\n5)Appartment record\n0)Skip/done edit\n\tSelect an index: ")
		switch index {
		case 0:
			break Pizza
		case 1:
			i.CI.LateFee = getInputFloat("Enter late fee: ")
		case 2:
			i.CI.IncurrAfter = getInputInt("Enter day after due to incur late fee: ")
		case 3:
			i.CI.BounceFee = getInputFloat("Enter check bounce fee: ")
		case 4:

		Pie:
			for {
				indx := getInputInt("Enter 1 to add, 2 to delete, 0 to skip: ")
				switch indx {
				case 0:
					break Pie
				case 1:
					i.CI.Apartments = inputAndAddString("Enter apartment name to add: ", i.CI.Apartments)
				case 2:
					i.CI.Apartments = inputAndDeleteString("Enter appartment name to remove: ", i.CI.Apartments)
				default:
					fmt.Printf("Incorrect selection, please try again\n")

				}
			}
		case 5:
		Pie2:
			for {
				displayCurrentAptList(i.CI.Apartments)
				apt := getInputStr("Enter the name of the apartment record to edit or return to skip: ")
				if apt == "" {
					break Pie2
				}
				if !isKnownApt(apt, i) {
					fmt.Printf("Incorrect selection, please try again\n")
					continue
				}
				editRental(i, apt)
			}
		default:
			fmt.Printf("Incorrect selection, please try again\n")

		}
	}

}

// Save Encodes via Gob to file
func Save(path string, object interface{}) error {
	// check if the source dir exist
	_, err := os.Stat(path)
	if err == nil {
		err := os.Rename(path, path+".bak")
		check(err)
	}
	//fmt.Printf("*******Stat.Error: %v\n", err)
	file, err := os.Create(path)
	if err == nil {
		defer file.Close()
		encoder := gob.NewEncoder(file)
		err = encoder.Encode(object)
	}
	check(err)
	return err
}

// Load Decodes Gob from file
func Load(path string, object interface{}) error {
	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	return err
}

func check(e error) {
	if e != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Println(line, "\t", file, "\n", e)
		os.Exit(1)
	}
}

func displayStatus(i *jawaInfo) {
	displayCommon(i)
	for _, apt := range *i.CI.Apartments {
		displayRental(i, apt, true)
	}

}

func updateRentsOwed(dataPtr *string, ji *jawaInfo) error {

	for _, v := range ji.Rental {
		ten := ji.Tenant[v.TenantKey]
		if timeNowRental().After(ten.RentChargedThru) {
			fmt.Printf("Updating %s RentChargedThru\n", v.TenantKey)
			if len(ten.Payment) > 0 {
				ten.RentOwed += v.Rent *
					float32(timeNowRental().Month()-ten.RentChargedThru.Month())
				fmt.Printf("Rent Owed is now %-6.2f\n", ten.RentOwed)
			}
			ten.RentChargedThru = now.EndOfMonth()
			if ji.Tenant == nil {
				ji.Tenant = map[string]renterRecord{}
			}
			ji.Tenant[v.TenantKey] = ten
		}
	}
	err := Save(*dataPtr, ji)
	return err
}

/*
func hardExit() {
	fmt.Printf("hardExit() called\n\n\n")
	fmt.Printf("hardExit() called\n\n\n")
	os.Exit(2) // this is to force all code to stop. mostly for net/http package where I don't have access to the server
}
*/

var dataTarget string

func main() {
	//defer hardExit() // to kill gui server after everything else
	dataPtr := pflag.String("dataloc", "", "Location for the appartment data file")
	timeNowPtr := pflag.String("timenow", "", "Set the current Now date as mm-dd-yyyy")
	versionPtr := pflag.Bool("version", false, "program version")
	editPtr := pflag.Bool("edit", false, "Edit Rental Information")
	recordPtr := pflag.Bool("record", false, "Record payment Information")
	initPayPtr := pflag.Bool("initPayment", false, "Record the initial payment for a Tenant")
	remindPtr := pflag.Bool("remind", false, "Create Reminder Letters")
	statusPtr := pflag.Bool("status", false, "Display the current status")
	WebGUIPtr := pflag.Bool("gui", false, "Use the Web Interface")
	pflag.Parse()

	if *versionPtr == true {
		fmt.Printf("\t Version %d.%d", version.major, version.minor)
		os.Exit(0)
	}
	if *dataPtr == "" {
		tempPath := os.TempDir()
		//fmt.Printf("TempDir: %s\n", tempPath)
		DirPath := filepath.Dir(tempPath)
		//fmt.Printf("DirPath: %s\n", DirPath)
		dataDir := filepath.Join(DirPath, "/RentalOps/")
		err := os.MkdirAll(dataDir, 666)
		if err != nil {
			fmt.Printf("MkdirAll(): %v\n", err)
		}
		*dataPtr = filepath.Join(dataDir, "Apartment.dat")
	}
	if *timeNowPtr != "" {
		var month, day, year int
		fmt.Sscanf(*timeNowPtr, "%d-%d-%d", &month, &day, &year)
		toDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
		initTimeRental(toDate)
	}

	ji := &jawaInfo{}
	err := Load(*dataPtr, ji)
	//check(err)
	fmt.Printf("Load err: %v\n", err)
	defer Save(*dataPtr, ji)
	if err == nil {
		updateRentsOwed(dataPtr, ji)
	}
	if ji.Rental == nil {
		// input file empty or non-existant, set up a basic default
		ji.Rental = make(map[string]rentalRecord)
		ji.Rental["undefined"] = rentalRecord{Apartment: "undefined"}
		ji.Tenant = make(map[string]renterRecord)
		ji.Tenant["undefined"] = renterRecord{Name: "undefined"}

	}
	if *editPtr == true {
		editInfo(ji)
		return
	}
	if *remindPtr == true || *recordPtr || *initPayPtr {
		process(ji, *recordPtr, *initPayPtr)
		return
	}
	if *statusPtr {
		displayStatus(ji)
		return
	}
	if *WebGUIPtr {
		dataTarget = *dataPtr
		establishWindow(ji)
	}
	pflag.Usage()
}

/*
var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	PrintDefaults()
}
*/
