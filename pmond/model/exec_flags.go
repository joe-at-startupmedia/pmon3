package model

import (
	"encoding/json"
)

type ExecFlags struct {
	User          string   `json:"user"`
	Log           string   `json:"log,omitempty" yaml:"log,omitempty" toml:"Log,omitempty"`
	LogDir        string   `json:"log_dir,omitempty" yaml:"log_dir,omitempty" toml:"LogDir,omitempty"`
	Args          string   `json:"args"`
	EnvVars       string   `json:"env_vars,omitempty" yaml:"env_vars,omitempty" toml:"EnvVars,omitempty"`
	Name          string   `json:"name"`
	Dependencies  []string `json:"dependencies,omitempty" yaml:"dependencies,omitempty" toml:"dependencies,omitempty"`
	Groups        []string `json:"groups,omitempty" yaml:"groups,omitempty" toml:"groups,omitempty" `
	NoAutoRestart bool     `json:"no_auto_restart" yaml:"no_auto_restart,omitempty" toml:"NoAutoRestart,omitempty"`
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
