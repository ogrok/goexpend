package ioutil

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

func WriteNewTemplate(item *models.ItemTemplate) error {
	// first validate recurrence input
	if item.Recurrence != "yearly" {
		item.RecurrenceMonth = 0
	}

	if item.Recurrence != "yearly" && item.Recurrence != "monthly" && item.Recurrence != "none" {
		return errors.New("Invalid recurrence parameter. Valid parameters are: `none`, `monthly`, `yearly`.\n")
	}

	// other parameters are passed cleanly. just need to deal with ID
	actualID, err := GetNextSequentialId()

	if err != nil || actualID < 1 {
		return err
	}

	item.ID = actualID

	// then write item. we load all existing items to double-check for ID conflicts
	// rather inefficient, but data is so small and simple in this program so it's fine

	file, err := ioutil.ReadFile(GetTemplateDataLoc())

	if err != nil {
		return err
	}

	var newFileContents []byte
	var templates []models.ItemTemplate

	if len(file) != 0 {
		err = json.Unmarshal(file, &templates)

		// single-objects need to be unmarshaled into single-obj var then appended to array
		if err != nil {
			var singleTemplate models.ItemTemplate
			err = json.Unmarshal(file, &singleTemplate)

			if err != nil {
				return err
			}

			templates = append(templates, singleTemplate)
		}

		templates = append(templates, *item)

		for k, v := range templates {
			for a, b := range templates {
				if a != k && v.ID == b.ID {
					return errors.New("ID conflict ("+string(v.ID)+"). No new item created")
				}
			}
		}

		sort.Slice(templates, func(i, j int) bool { return templates[i].ID < templates[j].ID })

		newFileContents, err = json.Marshal(templates)

		if err != nil {
			return err
		}
	} else {
		templates = append(templates, *item)

		newFileContents, err = json.Marshal(templates)
	}


	err = ioutil.WriteFile(GetTemplateDataLoc(), newFileContents, os.ModePerm)

	if err != nil {
		return err
	}

	fmt.Println("Budget item " + strconv.Itoa(actualID) + " created successfully")
	return nil
}