package util

import (
	"fmt"
	"log"
	"os"
	"os/user"
)

const dir = "/.goexpend"
const activeData = dir + "/active.json"
const logData = dir + "/log.json"

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

	err = file.Close()

	if err != nil {
		fmt.Printf("I/O error. It is possible operation completed. Check manually if " + GetActiveDataLoc() + " exists and rerun if needed.")
		os.Exit(1)
	}

	file, err = os.OpenFile(GetLogDataLoc(), os.O_CREATE, os.ModePerm)

	if err != nil {
		return err
	}

	err = file.Close()

	if err != nil {
		fmt.Printf("I/O error. It is possible operation completed. Check manually if " + GetLogDataLoc() + " exists and rerun if needed.")
		os.Exit(1)
	}

	fmt.Printf("init successful. Program is ready for use!\n")
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