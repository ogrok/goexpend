package main

import (
	"fmt"
	"github.com/adaminoue/goexpend/src/util"
	"os"
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

	case "delete":
		// delete budget item

	// change state of existing budget items
	case "modify":
		// edit accrued amount
	case "realize":
		// edit actual amount

	// view and change month state
	case "month":
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
	case "all":
		// list all budget items (name, amount, realized amount, category)
	case "report":
		// run report intended for viewing, on the current month
	default:
		// invalid argument
		showHelpTextAndExit()
	}
}

// TODO fill in this help text with basic syntax & args
// shows custom help text if input is invalid
func showHelpTextAndExit() {
	fmt.Printf("PLACEHOLDER HELP TEXT" + "\n")
	os.Exit(1)
}