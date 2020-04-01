package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adaminoue/goexpend/src/models"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

func WriteNewTemplate(item *models.Template, alsoMonthItem bool) (int, error) {
	// first validate recurrence input
	if item.Recurrence != "yearly" {
		item.RecurrenceMonth = 0
	}

	if item.Recurrence != "yearly" && item.Recurrence != "monthly" && item.Recurrence != "none" {
		return -1, errors.New("Invalid recurrence parameter. Valid parameters are: `none`, `monthly`, `yearly`.\n")
	}

	// other parameters are passed cleanly. just need to deal with ID
	actualID, err := GetNextSequentialId()

	if err != nil || actualID < 1 {
		return -1, err
	}

	item.ID = actualID

	// then write item. we load all existing items to double-check for ID conflicts
	// rather inefficient, but data is so small and simple in this program so it's fine

	file, err := ioutil.ReadFile(GetTemplateDataLoc())

	if err != nil {
		return -1, err
	}

	var newFileContents []byte
	var templates []models.Template

	if len(file) != 0 {
		err = json.Unmarshal(file, &templates)

		// single-objects need to be unmarshaled into single-obj var then appended to array
		if err != nil {
			var singleTemplate models.Template
			err = json.Unmarshal(file, &singleTemplate)

			if err != nil {
				return -1, err
			}

			templates = append(templates, singleTemplate)
		}

		templates = append(templates, *item)

		for k, v := range templates {
			for a, b := range templates {
				if a != k && v.ID == b.ID {
					return -1, errors.New("ID conflict ("+string(v.ID)+"). No new item created")
				}
			}
		}

		sort.Slice(templates, func(i, j int) bool { return templates[i].ID < templates[j].ID })

		newFileContents, err = json.Marshal(templates)

		if err != nil {
			return -1, err
		}
	} else {
		templates = append(templates, *item)

		newFileContents, err = json.Marshal(templates)
	}

	if alsoMonthItem {
		err = WriteNewMonthItem(item, 0)

		if err != nil {
			return -1, err
		}
	}

	err = ioutil.WriteFile(GetTemplateDataLoc(), newFileContents, os.ModePerm)

	if err != nil {
		// try to roll back creation of new month item from previous step
		_ = DeleteActiveItem(item.ID)

		return -1, err
	}

	fmt.Println("Budget item " + strconv.Itoa(actualID) + " created successfully")
	return actualID, nil
}

// create new item in active month concurrently with new template
func WriteNewMonthItem(input *models.Template, realizedAmount int) error {

	oneTime := true

	if input.Recurrence != "none" {
		oneTime = false
	}

	monthItem := models.ActiveItem{
		ID:       input.ID,
		Name:     input.Name,
		Category: input.Category,
		Accrued:  input.Amount,
		Realized: realizedAmount,
		Mutable:  input.Mutable,
		Amount:   input.Amount,
		OneTime:  oneTime,
	}

	file, err := ioutil.ReadFile(GetActiveDataLoc())

	if err != nil {
		return err
	}

	var newFileContents []byte
	var activeItems []models.ActiveItem

	if len(file) != 0 {
		err = json.Unmarshal(file, &activeItems)

		// single-objects need to be unmarshaled into single-obj var then appended to array
		if err != nil {
			var singleTemplate models.ActiveItem
			err = json.Unmarshal(file, &singleTemplate)

			if err != nil {
				return err
			}

			activeItems = append(activeItems, singleTemplate)
		}

		activeItems = append(activeItems, monthItem)

		for k, v := range activeItems {
			for a, b := range activeItems {
				if a != k && v.ID == b.ID {
					return errors.New("ID conflict ("+string(v.ID)+"). No new item created")
				}
			}
		}

		sort.Slice(activeItems, func(i, j int) bool { return activeItems[i].ID < activeItems[j].ID })

		newFileContents, err = json.Marshal(activeItems)

		if err != nil {
			return err
		}
	} else {
		activeItems = append(activeItems, monthItem)

		newFileContents, err = json.Marshal(activeItems)
	}

	err = ioutil.WriteFile(GetActiveDataLoc(), newFileContents, os.ModePerm)

	if err != nil {
		return err
	}

	return nil
}