package goex

import (
	"encoding/json"
	"errors"
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

func GetConfig() (Config, error) {
	if !ConfigExists() {
		return Config{}, errors.New("config does not exist")
	}

	file, err := ioutil.ReadFile(GetConfigDataLoc())
	var config Config

	if err != nil {
		return config, err
	}

	err = json.Unmarshal(file, &config)

	if err != nil {
		return Config{}, err
	}

	return config, nil
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

func WriteConfig() error {
	eom := endOfCurrentMonth()

	initialConfig := Config{
		CurrentMonth:  int(time.Now().Month()),
		CurrentYear:   time.Now().Year(),
		MonthEnd:      eom,
		AskAgainAfter: eom,
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

func UpdateAskAgainAfter(days int) error {
	if days == 0 {
		return nil
	}

	config, err := GetConfig()

	if err != nil {
		return err
	}

	config.AskAgainAfter += (86400 * days) // 1 day * number of days (passed as input)

	jsonConfig, err := json.Marshal(config)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(GetConfigDataLoc(), jsonConfig, os.ModePerm)

	if err != nil {
		return err
	}

	return nil
}