package ioutil

import (
	"encoding/json"
	"errors"
	"fmt"
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
		err = WriteConfig()

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
	var templates []ItemTemplate

	file, err := ioutil.ReadFile(GetTemplateDataLoc())

	if err != nil {
		return -1, err
	}

	if len(file) == 0 {
		return 1, nil
	}

	err = json.Unmarshal(file, &templates)

	if err != nil {
		var singleTemplate ItemTemplate
		err = json.Unmarshal(file, &singleTemplate)

		if err != nil {
			return -1, err
		}

		templates = append(templates, singleTemplate)
	}

	var activeItems []MonthItem

	file, err = ioutil.ReadFile(GetActiveDataLoc())

	if err != nil {
		return -1, err
	}

	if len(file) == 0 {
		return 1, nil
	}

	err = json.Unmarshal(file, &activeItems)

	if err != nil {
		var singleMonthItem MonthItem
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

func GetAllTemplates() ([]ItemTemplate, error) {
	var result []ItemTemplate

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

func GetAllActiveItems() ([]MonthItem, error) {
	var result []MonthItem

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

func GetSpecificTemplate(id int) (ItemTemplate, error) {
	all, err := GetAllTemplates()

	if err != nil {
		return ItemTemplate{}, err
	}

	for _, v := range all {
		if v.ID == id {
			return v, nil
		}
	}

	return ItemTemplate{}, errors.New("Item with ID "+strconv.Itoa(id)+" not found")
}

func GetSpecificActiveItem(id int) (MonthItem, error) {
	all, err := GetAllActiveItems()

	if err != nil {
		return MonthItem{}, err
	}

	for _, v := range all {
		if v.ID == id {
			return v, nil
		}
	}

	return MonthItem{}, errors.New("Item with ID "+strconv.Itoa(id)+" not found")
}