package app

import (
	"io"
	"io/ioutil"
	"os"
)

func (a *App) readLog(path string) (*string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	result := string(data)
	return &result, nil
}

func (a *App) writeToLog(path string, lastIssueID *string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, *lastIssueID)
	if err != nil {
		return err
	}
	return file.Sync()
}
