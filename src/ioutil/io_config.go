package ioutil

import (
	"encoding/json"
	"github.com/adaminoue/goexpend/src/models"
	"io/ioutil"
	"os"
	"time"
)

func ConfigExists() bool {
	_, err := os.OpenFile(GetConfigDataLoc(), os.O_RDONLY, os.ModePerm)

	if err != nil {
		return false
	}

	return true
}

// returns first moment of next month in Epoch time
func endOfCurrentMonth() int {
	now := time.Now()
	year := now.Year()
	month := now.Month()

	if month == 12 {
		year += 1
		month = 1
	} else {
		month += 1
	}

	return int(time.Date(year, month, 1, 0, 0, 0, 0, time.Local).Unix())
}

// only called upon initialization to save initial config
func saveInitialConfig() error {
	initialConfig := models.Config{
		CurrentMonth:  int(time.Now().Month()),
		CurrentYear:   time.Now().Year(),
		AskAgainAfter: endOfCurrentMonth(),
	}

	jsonConfig, err := json.Marshal(initialConfig)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(GetConfigDataLoc(), jsonConfig, os.ModePerm)

	if err != nil {
		return err
	}

	return nil
}