package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adaminoue/goexpend/src/models"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"
)

const dir = "/.goexpend"
const activeData = dir + "/active.json"
const logData = dir + "/log.json"
const configData = dir + "/config.json"
const templateData = dir + "/template.json"

var userHomeDir string

// creates some empty json files the rest of the app expects to exist; does not overwrite them if they do exist.
func Initialize() error {

	if _, err := os.Stat(GetHomeDir() + dir); os.IsNotExist(err) {
		err = os.Mkdir(GetHomeDir() + dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(GetActiveDataLoc(), os.O_CREATE, os.ModePerm)

	if err != nil {
		return err
	}

	blankFile := []byte("[]")

	_ = ioutil.WriteFile(GetActiveDataLoc(), blankFile, os.ModePerm)

	err = file.Close()

	if err != nil {
		fmt.Printf("I/O error. It is possible operation completed. Check manually if " + GetActiveDataLoc() + " exists and rerun if needed.")
		os.Exit(1)
	}

	file, err = os.OpenFile(GetLogDataLoc(), os.O_CREATE, os.ModePerm)

	_ = ioutil.WriteFile(GetLogDataLoc(), blankFile, os.ModePerm)

	if err != nil {
		return err
	}

	err = file.Close()

	if err != nil {
		fmt.Printf("I/O error. It is possible operation completed. Check manually if " + GetLogDataLoc() + " exists and rerun if needed.")
		os.Exit(1)
	}

	file, err = os.OpenFile(GetTemplateDataLoc(), os.O_CREATE, os.ModePerm)

	_ = ioutil.WriteFile(GetTemplateDataLoc(), blankFile, os.ModePerm)

	if err != nil {
		return err
	}

	err = file.Close()

	if err != nil {
		fmt.Printf("I/O error. It is possible operation completed. Check manually if " + GetLogDataLoc() + " exists and rerun if needed.")
		os.Exit(1)
	}

	if !ConfigExists() {
		err = WriteConfig(true, -1)

		if err != nil {
			fmt.Printf(err.Error()+"\nIt is possible operation completed. Check manually if " + GetLogDataLoc() + " exists and rerun if needed.")
			os.Exit(1)
		}

		fmt.Printf("init successful. Program is ready for use!\n")
	} else {
		fmt.Printf("Config file already existed at " + GetConfigDataLoc() + ".\nManually remove and run init again to recreate.\n")
	}

	return nil
}

func GetHomeDir() string {
	if userHomeDir == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		userHomeDir = usr.HomeDir
	}

	return userHomeDir
}

func GetDir() string {
	if userHomeDir == "" {
		userHomeDir = GetHomeDir()
	}
	return userHomeDir + dir
}

func GetActiveDataLoc() string {
	if userHomeDir == "" {
		userHomeDir = GetHomeDir()
	}
	return userHomeDir + activeData
}

func GetLogDataLoc() string {
	if userHomeDir == "" {
		userHomeDir = GetHomeDir()
	}
	return userHomeDir + logData
}

func GetTemplateDataLoc() string {
	if userHomeDir == "" {
		userHomeDir = GetHomeDir()
	}
	return userHomeDir + templateData
}

func GetConfigDataLoc() string {
	if userHomeDir == "" {
		userHomeDir = GetHomeDir()
	}
	return userHomeDir + configData
}

func GetNextSequentialId() (int, error) {
	var templates []models.Template

	file, err := ioutil.ReadFile(GetTemplateDataLoc())

	if err != nil {
		return -1, err
	}

	if len(file) == 0 {
		return 1, nil
	}

	err = json.Unmarshal(file, &templates)

	if err != nil {
		var singleTemplate models.Template
		err = json.Unmarshal(file, &singleTemplate)

		if err != nil {
			return -1, err
		}

		templates = append(templates, singleTemplate)
	}

	var activeItems []models.ActiveItem

	file, err = ioutil.ReadFile(GetActiveDataLoc())

	if err != nil {
		return -1, err
	}

	if len(file) == 0 {
		return 1, nil
	}

	err = json.Unmarshal(file, &activeItems)

	if err != nil {
		var singleMonthItem models.ActiveItem
		err = json.Unmarshal(file, &singleMonthItem)

		if err != nil {
			return -1, err
		}

		activeItems = append(activeItems, singleMonthItem)
	}

	// then find lowest candidate ID not in use and return it
	candidateId := 1

	for {
		goodCandidate := true

		for _, i := range templates {
			if candidateId == i.ID {
				goodCandidate = false
				break
			}
		}

		for _, j := range activeItems {
			if candidateId == j.ID {
				goodCandidate = false
				break
			}
		}

		if goodCandidate {
			return candidateId, nil
		} else {
			candidateId += 1
		}
	}
}

func GetAllTemplates() ([]models.Template, error) {
	var result []models.Template

	file, err := ioutil.ReadFile(GetTemplateDataLoc())

	if err != nil {
		return result, err
	}

	err = json.Unmarshal(file, &result)

	if err != nil {
		return result, err
	}

	return result, nil
}

func GetAllActiveItems() ([]models.ActiveItem, error) {
	var result []models.ActiveItem

	file, err := ioutil.ReadFile(GetActiveDataLoc())

	if err != nil {
		return result, err
	}

	err = json.Unmarshal(file, &result)

	if err != nil {
		return result, err
	}

	return result, nil
}

func GetSpecificTemplate(id int) (models.Template, error) {
	all, err := GetAllTemplates()

	if err != nil {
		return models.Template{}, err
	}

	for _, v := range all {
		if v.ID == id {
			return v, nil
		}
	}

	return models.Template{}, errors.New("Item with ID "+strconv.Itoa(id)+" not found")
}

func GetSpecificActiveItem(id int) (models.ActiveItem, error) {
	all, err := GetAllActiveItems()

	if err != nil {
		return models.ActiveItem{}, err
	}

	for _, v := range all {
		if v.ID == id {
			return v, nil
		}
	}

	return models.ActiveItem{}, errors.New("Item with ID "+strconv.Itoa(id)+" not found")
}

// returns amount of excess from template accrued amount, e.g. amount of manual accrual at the moment
func Excess(i *models.ActiveItem) int {
	template, err := GetSpecificTemplate(i.ID)

	if err != nil {
		return 0
	}

	return i.Accrued - template.Amount
}