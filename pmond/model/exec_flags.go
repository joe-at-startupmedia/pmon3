package model

import (
	"encoding/json"
)

type ExecFlags struct {
	User          string   `json:"user"`
	Log           string   `json:"log,omitempty"`
	LogDir        string   `json:"log_dir,omitempty"`
	Args          string   `json:"args"`
	EnvVars       string   `json:"env_vars,omitempty"`
	Name          string   `json:"name"`
	Dependencies  []string `json:"dependencies,omitempty"`
	Groups        []string `json:"groups,omitempty"`
	NoAutoRestart bool     `json:"no_auto_restart"`
}

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
