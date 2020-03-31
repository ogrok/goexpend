package ioutil

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
)

func DeleteItem(id int) error {
	_ = deleteActiveItem(id)

	err := deleteTemplateItem(id)

	if err != nil {
		return err
	}

	return nil
}

func deleteTemplateItem(id int) error {
	// delete item. we load all existing items and rewrite smaller set
	// rather inefficient, but simplistic and nice

	fileLoc := GetTemplateDataLoc()

	file, err := ioutil.ReadFile(fileLoc)

	if err != nil {
		return err
	}

	var newFileContents []byte
	var templates []ItemTemplate

	if len(file) == 0 {
		return errors.New("Nothing to delete")
	}

	err = json.Unmarshal(file, &templates)

	// maybe json only has one object in it
	if err != nil {
		var singleTemplate ItemTemplate

		err := json.Unmarshal(file, &singleTemplate)

		if err != nil {
			return err
		}

		templates = append(templates, singleTemplate)
	}

	deletionOccurred := false

	for k, v := range templates {
		if v.ID == id {
			templates = removeTemplate(templates, k)
			deletionOccurred = true
		}
	}

	newFileContents, err = json.Marshal(templates)

	err = ioutil.WriteFile(fileLoc, newFileContents, os.ModePerm)

	if err != nil {
		return err
	}

	if deletionOccurred {
		println("Item " + strconv.Itoa(id) + " deleted successfully.")
	} else {
		println("Item " + strconv.Itoa(id) + " not found.")
	}

	return nil
}

func deleteActiveItem(id int) error {
	// delete item. we load all existing items and rewrite smaller set
	// rather inefficient, but simplistic and nice

	fileLoc := GetActiveDataLoc()

	file, err := ioutil.ReadFile(fileLoc)

	if err != nil {
		return err
	}

	var newFileContents []byte
	var monthItems []MonthItem

	if len(file) == 0 {
		return errors.New("Nothing to delete")
	}

	err = json.Unmarshal(file, &monthItems)

	// maybe json only has one object in it
	if err != nil {
		var singleItem MonthItem

		err := json.Unmarshal(file, &singleItem)

		if err != nil {
			return err
		}

		monthItems = append(monthItems, singleItem)
	}

	for k, v := range monthItems {
		if v.ID == id {
			monthItems = removeActiveItem(monthItems, k)
		}
	}

	newFileContents, err = json.Marshal(monthItems)

	err = ioutil.WriteFile(fileLoc, newFileContents, os.ModePerm)

	if err != nil {
		return err
	}

	return nil
}

func removeTemplate(slice []ItemTemplate, s int) []ItemTemplate {
	return append(slice[:s], slice[s+1:]...)
}

func removeActiveItem(slice []MonthItem, s int) []MonthItem {
	return append(slice[:s], slice[s+1:]...)
}