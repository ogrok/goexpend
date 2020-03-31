package main

import (
	"flag"
	"fmt"
	"github.com/adaminoue/goexpend/src/models"
	"github.com/adaminoue/goexpend/src/ioutil"
	"os"
	"strconv"
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
		amountFlag := addCommand.Float64("a", 0.0, "Amount of new budget item")

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
			cleanError("No flags allowed in deletion command")
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
		amountFlag := modifyCommand.Float64("a", 0, "Accrued amount of budget item (zero-value is ignored)")
		categoryFlag := modifyCommand.String("c", "", "Category of budget item")
		descriptionFlag := modifyCommand.String("d", "", "Description of budget item")
		nameFlag := modifyCommand.String("n", "", "Name of budget item")
		realizedFlag := modifyCommand.Float64("r", 0, "Realized amount associated with budget item")

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
			cleanError("No flags allowed in accrue command")
		}

		accrueId, err := strconv.ParseInt(args[2], 10, 0)

		if err != nil {
			cleanError("Invalid ID")
		}

		accrueAmt, err := strconv.ParseFloat(args[3], 32)

		if err != nil {
			cleanError("Invalid amount")
		}

		accrue(int(accrueId), accrueAmt)
	case "realize": // add to actual amount (or subtract from with negative number)
		if len(args) != 4 {
			cleanError("No flags allowed in realize command")
		}

		realizeId, err := strconv.ParseInt(args[2], 10, 0)

		if err != nil {
			cleanError("Invalid ID")
		}

		realizeAmt, err := strconv.ParseFloat(args[3], 32)

		if err != nil {
			cleanError("Invalid amount")
		}

		realize(int(realizeId), realizeAmt)

	// view and change month state
	case "month":
		if len(args) == 2 {
			_ = showCurrentMonth()
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

func add(name string, amount float64, category string, description string, mutable bool, recurrence string) {
	newItem := models.ItemTemplate{
		ID:              0,
		Name:            name,
		Category:        category,
		Amount:          amount,
		Recurrence:      recurrence,
		RecurrenceMonth: showCurrentMonth(),
		Mutable:         mutable,
	}

	err := ioutil.WriteNewTemplate(&newItem)

	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
}

// TODO build out function
func del(itemId int)  {
	err := ioutil.DeleteItem(itemId)

	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
}

// TODO build out function
func accrue(itemId int, amount float64) {
	fmt.Printf("accrue function call successful")
}

// TODO build out function
func realize(itemId int, amount float64) {
	return
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

// TODO build out function
func info(itemId int) {
	return
}

// TODO build out function
func modify(itemId int, amount float64, category string, description string, name string, realizedEdit bool, realizedAmount float64) {
	return
}

// TODO build out this function
func closeMonth(force bool) {
	if !force {
		// ask for confirmation and exit if not received
	}

	// rest of function
}

// TODO build out this function
func reset(force bool) {
	if !force {
		// ask for confirmation and exit if not received
	}

	// rest of function
}

// TODO build out this function
func showCurrentMonth() int {
	// return current month for clarity, both printed and in return value. if non-current, ask to close it. ask daily, not every time
	// json required for this somewhere has current month as string and timestamp as "ask again after"
	return 0
}