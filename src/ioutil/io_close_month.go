package ioutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

func CloseMonth() error {
	// first check that config exists and can be used
	config, err := GetConfig()

	if err != nil {
		return err
	}

	// log all active items in their current state
	items, err := GetAllActiveItems()

	if len(items) == 0 {
		fmt.Println("No items to log. Advancing to current month...")

		var blank []ItemTemplate
		if err := generateNewMonth(&blank, true); err != nil {
			return err
		} else {
			return nil
		}
	}

	if err != nil {
		return err
	}

	var logs []MonthLogItem

	for _, item := range items {
		logs = append(logs, MonthLogItem{
			Name:        item.Name,
			Category:    item.Category,
			Description: item.Description,
			Accrued:     item.Accrued,
			Excess:      item.Excess(),
			Realized:    item.Realized,
			Remaining:   item.Remaining(),
			Mutable:     item.Mutable,
		})
	}

	newLogId, err := getNextLogId()

	if err != nil {
		return err
	}

	var completeMonthLog = MonthLog{
		ID:       newLogId,
		Month:    config.CurrentMonth,
		Year:     config.CurrentYear,
		LogItems: logs,
	}

	existingLogs, err := getExistingLogs()

	if err != nil {
		return err
	}

	existingLogs = append(existingLogs, completeMonthLog)

	sort.Slice(existingLogs, func(i, j int) bool { return existingLogs[i].ID < existingLogs[j].ID })

	err = saveOverLogFile(&existingLogs)

	if err != nil {
		return err
	}

	// now build next month: config and activeitems

	templates, err := GetAllTemplates()

	if err != nil {
		return err
	}

	return generateNewMonth(&templates, true)
}

func getNextLogId() (int, error) {
	existingLogs, err := getExistingLogs()

	if err != nil {
		return -1, err
	}

	// then find lowest candidate ID not in use and return it
	candidateId := 1

	for {
		goodCandidate := true

		for _, i := range existingLogs {
			if candidateId == i.ID {
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

func getExistingLogs() ([]MonthLog, error) {
	var result []MonthLog

	fileLoc := GetLogDataLoc()

	file, err := ioutil.ReadFile(fileLoc)

	if err != nil {
		return result, err
	}

	err = json.Unmarshal(file, &result)

	return result, err
}

func saveOverLogFile(logs *[]MonthLog) error {
	fileLoc := GetLogDataLoc()

	logsJson, err := json.Marshal(logs)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileLoc, logsJson, os.ModePerm)

	return err
}

// deeply inefficient because one-each R/W operation occurs per item. should refactor, but works
func generateNewMonth(templates *[]ItemTemplate, newConfig bool) error {
	fileLoc := GetActiveDataLoc()

	err := os.RemoveAll(fileLoc)

	err = ioutil.WriteFile(fileLoc, []byte("[]"), os.ModePerm)

	if err != nil {
		return err
	}

	for _, item := range *templates {
		err := WriteNewMonthItem(&item, 0)

		if err != nil {
			return err
		}
	}

	if newConfig {
		println("Welcome to " + time.Now().Month().String())
		return WriteConfig()
	}

	return nil
}