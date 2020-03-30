package main

import (
	"flag"
	"fmt"
	"github.com/adaminoue/goexpend/src/util"
	"os"
	"strconv"
)

func main() {
	// first deal with args
	args := os.Args

	if len(args) == 1 {
		// execute current default function
	}

	switch args[1] {
	case "init":
		err := util.Initialize()

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
				ShowErrorTextOfSomeKindAndExit() // no name or amount populated
			}

			add(*nameFlag, *amountFlag, *categoryFlag, *descriptionFlag, *mutableFlag, *recurrenceFlag)
		} else {
			ShowErrorTextOfSomeKindAndExit() // failed to parse args for add command
		}

	case "delete": // delete budget item
		if len(args) != 3 {
			ShowErrorTextOfSomeKindAndExit() // no flags allowed in deletion command
		}

		deleteId, err := strconv.ParseInt(args[2], 10, 0)

		if err != nil {
			ShowErrorTextOfSomeKindAndExit() // invalid id for deletion
		}

		del(int(deleteId))

	// change state of existing budget items
	case "modify": // edit accrued amount
		if len(args) < 3 {
			ShowErrorTextOfSomeKindAndExit() // no changes specified
		}

		// TODO continue building out modify logic
		// modifyCommand := flag.NewFlagSet("modify", flag.ExitOnError)

	case "accrue":
		// TODO build out accrue function logic
	case "realize": // add to actual amount (or subtract from with negative number)
		if len(args) != 4 {
			ShowErrorTextOfSomeKindAndExit() // no flags allowed in realize command
		}

		realizeId, err := strconv.ParseInt(args[2], 10, 0)

		if err != nil {
			ShowErrorTextOfSomeKindAndExit() // invalid id for deletion
		}

		realizeAmt, err := strconv.ParseFloat(args[3], 32)

		if err != nil {
			ShowErrorTextOfSomeKindAndExit() // invalid id for deletion
		}

		realize(int(realizeId), realizeAmt)

	// view and change month state
	case "month":
		// TODO build out month sub-function calls and handling
		// sub-switch related to month state
		// case "open":
			// open current month (to allow changes in realized amounts until closed)
		// case "close":
			// close current month (to not allow further changes of realized amounts)
		// case "reset":
			// reset current month (set all realized values in current month to zero)
		// case "change":
			// change current month
		// default:
			// return current month for clarity

	// reports and viewing
	case "info": // list info for one specific budget item
		if len(args) != 3 {
			ShowErrorTextOfSomeKindAndExit() // no flags allowed in info command
		}

		infoId, err := strconv.ParseInt(args[2], 10, 0)

		if err != nil {
			ShowErrorTextOfSomeKindAndExit() // invalid id
		}

		info(int(infoId))

	case "all": // list all budget items (name, amount, realized amount, category)
		if len(args) != 2 {
			ShowErrorTextOfSomeKindAndExit() // no args allowed in all command
		}
		all()
	case "report": // run report intended for viewing, on the current month
		if len(args) != 2 {
			ShowErrorTextOfSomeKindAndExit() // no args allowed in report command
		}
		report()
	default:
		ShowErrorTextOfSomeKindAndExit() // invalid command
	}
}

// TODO fill in this help text with stuff related to its locations in the code
// shows custom help text if input is invalid
func ShowErrorTextOfSomeKindAndExit() {
	fmt.Printf("PLACEHOLDER HELP TEXT" + "\n")
	os.Exit(1)
}

// TODO build out function
func add(name string, amount float64, category string, description string, mutable bool, recurrence string) {
	return
}

// TODO build out function
func del(itemId int)  {
	return
}

// TODO build out function
func realize(itemId int, amount float64) {
	return
}

// TODO build out function
func all() {
	return
}

// TODO build out function
func report() {
	return
}

// TODO build out function
func info(itemId int) {
	return
}