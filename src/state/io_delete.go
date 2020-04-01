package state

import (
	"encoding/json"
	"errors"
	"github.com/adaminoue/goexpend/src/models"
	"io/ioutil"
	"os"
	"strconv"
)

func DeleteItem(id int, deleteTemplate bool) error {
	_ = DeleteActiveItem(id)

	if deleteTemplate {
		err := DeleteTemplateItem(id)

		if err != nil {
			return err
		}
	}


	return nil
}

func DeleteTemplateItem(id int) error {
	// delete item. we load all existing items and rewrite smaller set
	// rather inefficient, but simplistic and nice

	fileLoc := GetTemplateDataLoc()

	file, err := ioutil.ReadFile(fileLoc)

	if err != nil {
		return err
	}

	var newFileContents []byte
	var templates []models.Template

	if len(file) == 0 {
		return errors.New("Nothing to delete")
	}

	err = json.Unmarshal(file, &templates)

	// maybe json only has one object in it
	if err != nil {
		var singleTemplate models.Template

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

func DeleteActiveItem(id int) error {
	// delete item. we load all existing items and rewrite smaller set
	// rather inefficient, but simplistic and nice

	fileLoc := GetActiveDataLoc()

	file, err := ioutil.ReadFile(fileLoc)

	if err != nil {
		return err
	}

	var newFileContents []byte
	var monthItems []models.ActiveItem

	if len(file) == 0 {
		return errors.New("Nothing to delete")
	}

	err = json.Unmarshal(file, &monthItems)

	// maybe json only has one object in it
	if err != nil {
		var singleItem models.ActiveItem

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

func removeTemplate(slice []models.Template, s int) []models.Template {
	return append(slice[:s], slice[s+1:]...)
}

func removeActiveItem(slice []models.ActiveItem, s int) []models.ActiveItem {
	return append(slice[:s], slice[s+1:]...)
}