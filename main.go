package main

import (
	"fmt"
	"os"
)

func main() {
	// first deal with args
	args := os.Args

	if len(args) <= 1 {
		showHelpTextAndExit()
	}

	switch args[1] {
	case "init":
		err := initialize()

		if err != nil {
			fmt.Printf(err.Error() + "\n")
			os.Exit(1)
		}

	// adding and removing entire budget items
	case "add":
		// add new budget item
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
			// open current month (to allow changes until closed)
		// case "close":
			// close current month (to not allow further changes until opened)
		// case "reset":
			// reset existing month (set all actual values in current month to zero)
		// default:
			// return current month for clarity

	// reports and viewing
	case "all":
		// list all budget items (name, amount, realized amount)
	case "report":
		// run report intended for viewing
	default:
		// run quick-default behavior (which ought to be configurable)
	}
}

// TODO fill in this help text with basic syntax & args
// shows custom help text if input is invalid
func showHelpTextAndExit() {
	fmt.Printf("PLACEHOLDER HELP TEXT" + "\n")
	os.Exit(1)
}