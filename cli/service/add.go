package service

import (
	"encoding/json"
	"fmt"
	"pmon3/pmond"
	"pmon3/pmond/model"
)

func AddData(process *model.Process) (string, error) {
	// save to db
	var originOne model.Process
	err := pmond.Db().First(&originOne, "process_file = ? AND name = ?", process.ProcessFile, process.Name).Error
	if err == nil && originOne.ID > 0 { // process already exist
		process.ID = originOne.ID
	}

	err = pmond.Db().Save(&process).Error
	if err != nil {
		return "", fmt.Errorf("pmon3 run err: %w", err)
	}

	output, err := json.Marshal(process.RenderTable())
	if err != nil {
		return "", err
	}

	return string(output), nil
}
