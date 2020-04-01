package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/adaminoue/goexpend/src/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// parsing and validation of input occurs in main.
// no json is manipulated in main; this is abstracted out to the functions below it
// see `ioutil` folder for shared functions related to I/O etc.
func main() {
	args := os.Args

	// first check if config exists. force creation if not
	configExists := ioutil.ConfigExists()

	if !configExists {
		if args[1] != "init" {
			fmt.Printf("Config file does not exist at " + ioutil.GetConfigDataLoc() + ".\nDid you run `goex init` yet?\n")
			os.Exit(1)
		}
	}

	// then check config for whether we are past AskAgainAt
	config, _ := ioutil.GetConfig()

	if time.Now().Unix() > int64(config.AskAgainAfter) && args[1] != "purge" && args[1] != "month" {
		if askUserToClose(true) {
			closeMonth(true)
		} else {
			_ = ioutil.UpdateAskAgainAfter(1)
		}
	}

	// then begin processing input

	if len(args) == 1 {
		// execute current default function. TODO make this configurable
		report()
	}

	switch args[1] {
	case "init":
		err := ioutil.Initialize()

		if err != nil {
			fmt.Printf(err.Error() + "\n")
			os.Exit(1)
		}

	// adding and removing entire budget items
	case "add": // add new budget item
		addCommand := flag.NewFlagSet("add", flag.ExitOnError)
		nameFlag := addCommand.String("n", "", "Name of new budget item")
		amountFlag := addCommand.Int("a", 0.0, "Amount of new budget item")

		categoryFlag := addCommand.String("c", "", "Category of new budget item")
		descriptionFlag := addCommand.String("d", "", "Description of new budget item")
		mutableFlag := addCommand.Bool("m", true, "Mutability of new budget item")
		recurrenceFlag := addCommand.String("r", "monthly", "Recurrence behavior of new budget item")

		err := addCommand.Parse(args[2:])

		if err != nil {
			fmt.Printf(err.Error() + "\n")
			os.Exit(1)
		}

		if addCommand.Parsed() {
			if *nameFlag == "" || *amountFlag == 0 {
				cleanError("Both name and amount required for add command")
			}

			add(*nameFlag, *amountFlag, *categoryFlag, *descriptionFlag, *mutableFlag, *recurrenceFlag)
		} else {
			cleanError("Failed to parse arguments for add command")
		}

	case "delete": // delete budget item
		if len(args) != 3 {
			cleanError("Invalid input")
		}

		deleteId, err := strconv.ParseInt(args[2], 10, 0)

		if err != nil {
			cleanError("Invalid ID for deletion")
		}

		del(int(deleteId))

	// change state of existing budget items
	case "modify": // edit accrued amount
		if len(args) < 4 {
			cleanError("No modifications specified")
		}

		modifyId, err := strconv.ParseInt(args[2], 10, 0)

		if err != nil {
			cleanError("Invalid ID")
		}

		modifyCommand := flag.NewFlagSet("modify", flag.ExitOnError)
		amountFlag := modifyCommand.Int("a", 0, "Accrued amount of budget item (zero-value is ignored)")
		categoryFlag := modifyCommand.String("c", "", "Category of budget item")
		descriptionFlag := modifyCommand.String("d", "", "Description of budget item")
		nameFlag := modifyCommand.String("n", "", "Name of budget item")
		realizedFlag := modifyCommand.Int("r", 0, "Realized amount associated with budget item")

		err = modifyCommand.Parse(args[3:])

		// need to specifically check whether -r is used bc we want to allow setting realized value to 0;
		// default value is therefore invalid null case
		var realizedEdit = false

		for _, arg := range args {
			if arg == "-r" || arg == "--r" {
				realizedEdit = true
			}
		}

		if err != nil {
			fmt.Printf(err.Error() + "\n")
			os.Exit(1)
		}

		if modifyCommand.Parsed() {
			modify(int(modifyId), *amountFlag, *categoryFlag, *descriptionFlag, *nameFlag, realizedEdit, *realizedFlag)
		} else {
			cleanError("Failed to parse arguments for modify command")
		}

	case "accrue": // add to accrued amount (or subtract from with negative number)
		if len(args) != 4 {
			cleanError("Invalid input")
		}

		accrueId, err := strconv.ParseInt(args[2], 10, 0)

		if err != nil {
			cleanError("Invalid ID")
		}

		accrueAmt, err := strconv.Atoi(args[3])

		if err != nil {
			cleanError("Invalid amount")
		}

		accrue(int(accrueId), accrueAmt)
	case "realize": // add to actual amount (or subtract from with negative number)
		if len(args) != 4 {
			cleanError("Invalid input")
		}

		realizeId, err := strconv.ParseInt(args[2], 10, 0)

		if err != nil {
			cleanError("Invalid ID")
		}

		realizeAmt, err := strconv.Atoi(args[3])

		if err != nil {
			cleanError("Invalid amount")
		}

		realize(int(realizeId), realizeAmt)

	// view and change month state
	case "month":
		if len(args) == 2 {
			_ = showCurrentMonth()
			os.Exit(0)
		}

		// sub-switch on
		switch args[2] {
		case "close": // close current month, permanently logging it
			if len(args) > 3 && args[3] == "-f" {
				closeMonth(true)
			} else {
				closeMonth(false)
			}
		case "reset": // reset current month (set all realized values in current month to zero)
			if len(args) > 3 && args[3] == "-f" {
				reset(true)
			} else {
				reset(false)
			}
		default:
			cleanError("Invalid command")
		}

	// reports and viewing
	case "info": // list info for one specific budget item
		if len(args) != 3 {
			cleanError("No flags allowed in info command")
		}

		infoId, err := strconv.ParseInt(args[2], 10, 0)

		if err != nil {
			cleanError("Invalid ID")
		}

		info(int(infoId))

	case "all": // list all budget items (name, amount, realized amount, category)
		if len(args) != 2 {
			cleanError("No arguments allowed in all command")
		}
		all()
	case "report": // run report intended for viewing, on the current month
		if len(args) != 2 {
			cleanError("no arguments allowed in report command")
		}
		report()

	case "purge": // remove all data then run init
		if userConfirms("remove all data and reset goexpend") {
			err := Purge()

			if err != nil {
				fmt.Printf(err.Error())
				os.Exit(1)
			}
		} else {
			fmt.Printf("Purge aborted")
		}

	default:
		cleanError("Invalid command")
	}
}

// TODO fill in this help text with stuff related to its locations in the code
// shows custom help text if input is invalid
func cleanError(input string) {
	fmt.Printf(input + "\n")
	os.Exit(1)
}

func add(name string, amount int, category string, description string, mutable bool, recurrence string) {
	newItem := ioutil.ItemTemplate{
		ID:              0,
		Name:            name,
		Category:        category,
		Description:     description,
		Amount:          amount,
		Recurrence:      recurrence,
		RecurrenceMonth: showCurrentMonth(),
		Mutable:         mutable,
	}

	id, err := ioutil.WriteNewTemplate(&newItem, true)

	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	if recurrence == "none" {
		err := ioutil.DeleteTemplateItem(id)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
	}
}

func del(itemId int)  {
	err := ioutil.DeleteItem(itemId, true)

	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
}

func accrue(itemId int, amount int) {
	currentItem, err := ioutil.GetSpecificActiveItem(itemId)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var mod = ioutil.ModTemplate{
		ID:          itemId,
		Amount:      currentItem.Accrued + amount,
	}

	err = ioutil.ModifyItem(mod, false, false)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	println("Accrued amount of " + strconv.Itoa(amount))
}

func realize(itemId int, amount int) {
	currentItem, err := ioutil.GetSpecificActiveItem(itemId)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var mod = ioutil.ModTemplate{
		ID:          itemId,
		Realized:    currentItem.Realized + amount,
	}

	err = ioutil.ModifyItem(mod, true, false)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}



	println("Realized amount of " + strconv.Itoa(amount))
}

// TODO build out function
func all() {
	if ioutil.ConfigExists() {
		fmt.Printf("CONFIG EXISTS")
	} else {
		fmt.Printf("Config does NOT exist!")
	}
}

// TODO build out function
func report() {
	return
}

func info(itemId int) {
	templateExists := true

	template, err := ioutil.GetSpecificTemplate(itemId)

	if err != nil {
		templateExists = false
	}

	activeItem, err := ioutil.GetSpecificActiveItem(itemId)

	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	// account for when template doesn't exist, in case of non-recurrence
	var amountToShow int
	var recurToShow string
	var monthToShow int

	if templateExists {
		amountToShow = template.Amount
		recurToShow = template.Recurrence
		monthToShow = template.RecurrenceMonth
	} else {
		amountToShow = activeItem.Accrued
		recurToShow = "none"
		monthToShow = -1
	}

	viewmodel := ioutil.ViewmodelInfo{
		ID:              activeItem.ID,
		Name:            activeItem.Name,
		Category:        activeItem.Category,
		Description:     activeItem.Description,
		CurrentAccrued:  activeItem.Accrued,
		Realized:        activeItem.Realized,
		Mutable:         activeItem.Mutable,
		Amount:          amountToShow,
		Recurrence:      recurToShow,
		RecurrenceMonth: monthToShow,
	}

	var recurrenceDesc string

	if viewmodel.Recurrence == "yearly" {
		recurrenceDesc = "yearly in " + strings.ToLower(time.Month(viewmodel.RecurrenceMonth).String())
	} else if viewmodel.Recurrence == "monthly" {
		recurrenceDesc = "monthly"
	} else if viewmodel.Recurrence == "none" {
		recurrenceDesc = "none"
	}

	fmt.Printf(""+
		"\nName:                " + viewmodel.Name +
		"\nID:                  " + strconv.Itoa(viewmodel.ID) +
		"\nCategory:            " + viewmodel.Category +
		"\nDescription:         " + viewmodel.Description +
		"\nRegular Amount:      " + strconv.Itoa(viewmodel.Amount) +
		"\nCurrent Accrual:     " + strconv.Itoa(viewmodel.CurrentAccrued) +
		"\nRealized / Remains:  " + strconv.Itoa(viewmodel.Realized) + " / " + strconv.Itoa(viewmodel.Remains()) +
		"\nMutable:             " + strconv.FormatBool(viewmodel.Mutable) +
		"\nRecurs:              " + recurrenceDesc +
		"\n")
}

func modify(itemId int, amount int, category string, description string, name string, realizedEdit bool, realizedAmount int) {
	var modTemplate = ioutil.ModTemplate{
		ID:          itemId,
		Amount:      amount,
		Category:    category,
		Description: description,
		Name:        name,
		Realized:    realizedAmount,
	}

	err := ioutil.ModifyItem(modTemplate, realizedEdit, true)

	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func closeMonth(force bool) {
	config, err := ioutil.GetConfig()

	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	if int64(config.MonthEnd) < time.Now().Unix() {
		println("Error: You cannot close the month until it is over")
		os.Exit(1)
	}

	if !force {
		if !userConfirms("close the current month and open the month of " + time.Now().Month().String()) {
			os.Exit(1)
		}
	}

	if err := ioutil.CloseMonth(); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

// TODO build out this function
func reset(force bool) {
	if !force {
		userConfirms("reset the active month")
	}

	// do the function
}

func showCurrentMonth() int {
	// return current month for clarity, both printed and in return value. if non-current, ask to close it. ask daily, not every time
	// json required for this somewhere has current month as string and timestamp as "ask again after"
	return int(time.Now().Month())
}

func userConfirms(operation string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Are you sure you would like to " + operation + "? [Y/n] ")
	text, _ := reader.ReadString('\n')
	if strings.ToUpper(strings.TrimSpace(text)) == "Y" {
		return true
	}

	return false
}

func askUserToClose(summary bool) bool {
	reader := bufio.NewReader(os.Stdin)
	if summary {
		config, err := ioutil.GetConfig()

		if err == nil {
			fmt.Println("It is now " + time.Now().Month().String() + " " +
				strconv.Itoa(time.Now().Year()) + " and the current active month is " +
				time.Month(config.CurrentMonth).String() + " " +
				strconv.Itoa(config.CurrentYear) + ".")
		}
	}
	fmt.Print("Would you like to close the old month and open the budget for the month of "+ time.Now().Month().String() +"? [Y/n] ")
	text, _ := reader.ReadString('\n')
	if strings.ToUpper(strings.TrimSpace(text)) == "Y" {
		return true
	}

	return false
}

func Purge() error {
	dir := ioutil.GetDir()

	err := os.RemoveAll(dir)

	if err != nil {
		return err
	}

	fmt.Println("Purge complete. Rerunning init...")
	err = ioutil.Initialize()

	if err != nil {
		return err
	}

	return nil
}