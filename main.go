package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/adaminoue/goexpend/src/models"
	"github.com/adaminoue/goexpend/src/state"
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
	configExists := state.ConfigExists()

	if !configExists {
		if args[1] != "init" {
			fmt.Printf("Config file does not exist at " + state.GetConfigDataLoc() + ".\nDid you run `goex init` yet?\n")
			os.Exit(1)
		}
	}

	// then check config for whether we are past AskAgainAt and prompt user to close if so
	config, _ := state.GetConfig()

	if time.Now().Unix() > int64(config.AskAgainAfter) && args[1] != "purge" && args[1] != "month" {
		suggestMonthClose()
	}

	// then begin processing input

	if len(args) == 1 {
		// execute current default function. TODO make this configurable
		err := state.ShowFullReport()

		if err != nil {
			println(err.Error())
			os.Exit(1)
		}

		os.Exit(0)
	}

	switch args[1] {
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
	case "add": // add new budget item

		// first grab all the flags
		addCommand := flag.NewFlagSet("add", flag.ExitOnError)

		nameFlag := addCommand.String("n", "", "Name of new budget item")
		amountFlag := addCommand.Int("a", 0.0, "Amount of new budget item")
		categoryFlag := addCommand.String("c", "", "Category of new budget item")
		descriptionFlag := addCommand.String("d", "", "Description of new budget item")
		immutableFlag := addCommand.Bool("i", false, "Mutability of new budget item")
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

			newItem := models.Template{
				ID:              0,
				Name:            *nameFlag,
				Category:        *categoryFlag,
				Description:     *descriptionFlag,
				Amount:          *amountFlag,
				Recurrence:      *recurrenceFlag,
				RecurrenceMonth: currentMonth(false),
				Immutable:       *immutableFlag,
			}

			add(&newItem)
		} else {
			cleanError("Failed to parse arguments for add command")
		}
	case "all": // list all budget items (name, amount, realized amount, category)
		if len(args) != 2 {
			cleanError("No arguments allowed in all command")
		}
		all()
	case "delete": // delete budget item
		if len(args) != 3 {
			cleanError("Invalid input")
		}

		deleteId, err := strconv.ParseInt(args[2], 10, 0)

		if err != nil {
			cleanError("Invalid ID for deletion")
		}

		del(int(deleteId))
	case "income":
		if len(args) != 3 {
			println("Invalid input; please provide number only")
			os.Exit(1)
		}

		intIncome, _ := strconv.Atoi(args[2])

		income(intIncome)
	case "info": // list info for one specific budget item
		if len(args) != 3 {
			cleanError("No flags allowed in info command")
		}

		infoId, err := strconv.ParseInt(args[2], 10, 0)

		if err != nil {
			cleanError("Invalid ID")
		}

		info(int(infoId))
	case "init":
		err := state.Initialize()

		if err != nil {
			fmt.Printf(err.Error() + "\n")
			os.Exit(1)
		}
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
			var mods = models.Modification{
				ID:          int(modifyId),
				Amount:      *amountFlag,
				Category:    *categoryFlag,
				Description: *descriptionFlag,
				Name:        *nameFlag,
				Realized:    *realizedFlag,
			}

			modify(&mods, realizedEdit)
		} else {
			cleanError("Failed to parse arguments for modify command")
		}
	case "month":
		if len(args) == 2 {
			_ = currentMonth(true)
			os.Exit(0)
		}

		// sub-switch
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
	case "purge": // remove all data then run init
		if userConfirms("remove all data and reset goexpend") {
			err := purge()

			if err != nil {
				fmt.Printf(err.Error())
				os.Exit(1)
			}
		} else {
			fmt.Printf("purge aborted")
		}
	case "realize": // add to actual amount (or subtract from with negative number)

		if len(args) < 3 || len(args) > 4 {
			cleanError("Invalid input")
		}

		realizeId, err := strconv.ParseInt(args[2], 10, 0)

		if err != nil {
			cleanError("Invalid ID")
		}

		switch len(args) {
		case 3:
			item, err := state.GetSpecificActiveItem(int(realizeId))

			if err != nil {
				cleanError("Could not fetch full amount to realize")
			}

			realize(int(realizeId), item.Remaining())
		case 4:
			realizeAmt, err := strconv.Atoi(args[3])

			if err != nil {
				cleanError("Invalid amount")
			}

			realize(int(realizeId), realizeAmt)
		}
	case "report": // run report intended for viewing, on the current month
		if len(args) != 2 {
			cleanError("no arguments allowed in report command")
		}
		err := state.ShowFullReport()

		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
	default:
		cleanError("Invalid command")
	}

	os.Exit(0)
}

// shows custom help text if input is invalid
func cleanError(input string) {
	fmt.Printf(input + "\n")
	os.Exit(0)
}

func add(newItem *models.Template) {

	id, err := state.WriteNewTemplat(newItem, true)

	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	if newItem.Recurrence == "none" {
		err := state.DeleteTemplateItem(id, false)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
	}
}

func del(itemId int) {
	err := state.DeleteItem(itemId)

	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
}

func accrue(itemId int, amount int) {
	currentItem, err := state.GetSpecificActiveItem(itemId)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var mod = models.Modification{
		ID:     itemId,
		Amount: currentItem.Accrued + amount,
	}

	err = state.ModifyItem(&mod, false, false)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	println("Accrued amount of " + strconv.Itoa(amount))
}

func realize(itemId int, amount int) {
	currentItem, err := state.GetSpecificActiveItem(itemId)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var mod = models.Modification{
		ID:       currentItem.ID,
		Realized: currentItem.Realized + amount,
	}

	err = state.ModifyItem(&mod, true, false)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	println("Realized amount of " + strconv.Itoa(amount))
}

func all() {
	items, err := state.GetAllActiveItems()

	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	if len(items) == 0 {
		println("No items to list! Try adding some items with `goex add -n {NAME} -a {AMOUNT}`.\nCheck README for more info.")
	}

	for _, item := range items {
		info(item.ID)
	}
}

func info(itemId int) {
	templateExists := true

	template, err := state.GetSpecificTemplate(itemId)

	if err != nil {
		templateExists = false
	}

	activeItem, err := state.GetSpecificActiveItem(itemId)

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

	viewmodel := models.ActiveItemView{
		ID:              activeItem.ID,
		Name:            activeItem.Name,
		Category:        activeItem.Category,
		Description:     activeItem.Description,
		CurrentAccrued:  activeItem.Accrued,
		Realized:        activeItem.Realized,
		Immutable:       activeItem.Immutable,
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

	fmt.Printf("" +
		"\nName:                " + viewmodel.Name +
		"\nID:                  " + strconv.Itoa(viewmodel.ID) +
		"\nCategory:            " + viewmodel.Category +
		"\nDescription:         " + viewmodel.Description +
		"\nRegular Amount:      " + strconv.Itoa(viewmodel.Amount) +
		"\nCurrent Accrual:     " + strconv.Itoa(viewmodel.CurrentAccrued) +
		"\nRealized / Remains:  " + strconv.Itoa(viewmodel.Realized) + " / " + strconv.Itoa(viewmodel.Remains()) +
		"\nImmutable:           " + strconv.FormatBool(viewmodel.Immutable) +
		"\nRecurs:              " + recurrenceDesc +
		"\n")
}

func modify(mods *models.Modification, realizedEdit bool) {
	err := state.ModifyItem(mods, realizedEdit, true)

	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func closeMonth(force bool) {
	config, err := state.GetConfig()

	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	if int64(config.MonthEnd) > time.Now().Unix() {
		println("Error: You cannot close the month until it is over")
		os.Exit(1)
	}

	if !force {
		if !userConfirms("close the current month and open the month of " + time.Now().Month().String()) {
			os.Exit(1)
		}
	}

	if err := state.CloseMonth(); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func reset(force bool) {
	if !force {
		if !userConfirms("reset the active month") {
			println("Aborted")
			os.Exit(0)
		}
	}

	err := state.ResetMonth()

	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func income(income int) {
	err := state.WriteConfig(false, income)

	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	println("Income updated to " + strconv.Itoa(income))
}

func currentMonth(showToUser bool) int {
	now := time.Now()

	if showToUser {
		println("Current active month is " + now.Month().String() + " " + strconv.Itoa(now.Year()))
	}

	return int(time.Now().Month())
}

func purge() error {
	dir := state.GetDir()

	err := os.RemoveAll(dir)

	if err != nil {
		return err
	}

	fmt.Println("purge complete. Rerunning init...")
	err = state.Initialize()

	if err != nil {
		return err
	}

	return nil
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
		config, err := state.GetConfig()

		if err == nil {
			fmt.Println("It is now " + time.Now().Month().String() + " " +
				strconv.Itoa(time.Now().Year()) + " and the current active month is " +
				time.Month(config.CurrentMonth).String() + " " +
				strconv.Itoa(config.CurrentYear) + ".")
		}
	}
	fmt.Print("Would you like to close the old month and open the budget for the month of " + time.Now().Month().String() + "? [Y/n] ")
	text, _ := reader.ReadString('\n')
	if strings.ToUpper(strings.TrimSpace(text)) == "Y" {
		return true
	}

	return false
}

func suggestMonthClose() {
	if askUserToClose(true) {
		closeMonth(true)
	} else {
		_ = state.UpdateAskAgainAfter(1)
	}
}
