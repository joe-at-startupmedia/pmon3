package model

import (
	"encoding/json"
)

func (e *ExecFlags) Parse(jsonStr string) (*ExecFlags, error) {
	var m ExecFlags
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (e *ExecFlags) Json() string {
	content, _ := json.Marshal(e)

	return string(content)
}

func ExecFlagsNames(execFlagsPtr *[]ExecFlags) []string {

	if execFlagsPtr == nil {
		return []string{}
	}

	execFlags := *execFlagsPtr

	if len(execFlags) == 0 {
		return []string{}
	}

	keys := make([]string, len(execFlags))

	i := 0
	for _, execFlag := range execFlags {
		keys[i] = execFlag.Name
		i++
	}

	return keys
}
