package model

import (
	"encoding/json"
	"os/user"
)

type ExecFlags struct {
	User          string   `json:"user"`
	Log           string   `json:"log"`
	LogDir        string   `json:"log_dir"`
	Args          string   `json:"args"`
	EnvVars       string   `json:"env_vars"`
	Name          string   `json:"name"`
	Dependencies  []string `json:"dependencies"`
	Groups        []string `json:"groups"`
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

func (e *ExecFlags) SetCurrentUser() {
	if e.User == "" {
		user, err := user.Current()
		if err == nil {
			e.User = user.Username
		}
	}
}

func (e *ExecFlags) Json() string {
	content, _ := json.Marshal(e)

	return string(content)
}
