package state

import (
	"github.com/adaminoue/goexpend/src/models"
	"os"
	"strconv"
	"time"
)

func ResetMonth() error {
	items, err := GetAllActiveItems()

	if err != nil {
		return err
	}

	var mod = models.Modification{ Realized: 0 }

	errCount := 0

	for _, item := range items {
		mod.ID = item.ID
		err = ModifyItem(&mod, true, false)

		if err != nil {
			errCount += 1
		}
	}

	if errCount > 0 {
		println(strconv.Itoa(errCount) + " errors occurred during month reset. Please check your filesystem and try again.")
		os.Exit(1)
	}

	err = WriteConfig(true, -1)

	if err != nil {
		return err
	}

	println("Reset successful. Welcome to " + time.Now().Month().String())
	return nil
}