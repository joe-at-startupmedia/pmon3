package exec

import (
	"encoding/json"
	"fmt"
	"pmon3/cli/proxy"
	"pmon3/pmond/model"
)

func restart(m *model.Process, flags string) ([]string, error) {
	// only stopped process or failed process allow run start
	if m.Status == model.StatusStopped || m.Status == model.StatusFailed {
		newData, err := reloadProcess(m, flags)
		if err != nil {
			return nil, err
		}

		return newData, nil
	}

	return nil, fmt.Errorf("process already running with the name provided: %s", m.Name)
}

func reloadProcess(m *model.Process, flags string) ([]string, error) {
	data, err := proxy.RunProcess([]string{"restart", m.ProcessFile, flags})

	if err != nil {
		return nil, err
	}

	var rel []string
	err = json.Unmarshal(data, &rel)
	if err != nil {
		return nil, err
	}

	return rel, nil
}
