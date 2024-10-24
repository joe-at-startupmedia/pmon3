package model

import (
	"gorm.io/gorm"
	"os"
	"time"
)

const dateTimeFormat = "2006-01-02 15:04:05"

var processRestartCounter = make(map[uint32]uint32)

type ProcessStatus int64

const (
	StatusQueued ProcessStatus = iota
	StatusInit
	StatusRunning
	StatusStopped
	StatusFailed
	StatusClosed
	StatusBackoff
	StatusRestarting
)

type Process struct {
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Pointer      *os.Process   `gorm:"-" json:"-"`
	Log          string        `gorm:"column:log" json:"log"`
	Name         string        `gorm:"unique" json:"name"`
	ProcessFile  string        `json:"process_file"`
	Args         string        `json:"args"`
	EnvVars      string        `json:"env_vars"`
	Username     string        `json:"username"`
	MemoryUsage  string        `json:"memory_usage"`
	CpuUsage     string        `json:"cpu_usage"`
	Dependencies string        `json:"dependencies"`
	Groups       []*Group      `gorm:"many2many:process_groups;"`
	Status       ProcessStatus `json:"status"`
	ID           uint32        `gorm:"primary_key" json:"id"`
	Pid          uint32        `gorm:"column:pid" json:"pid"`
	Uid          uint32        `gorm:"column:uid" json:"uid"`
	Gid          uint32        `gorm:"column:gid" json:"gid"`
	RestartCount uint32        `gorm:"-" json:"-"`
	AutoRestart  bool          `json:"auto_restart"`
}

type ExecFlags struct {
	File          string   `json:"file"`
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

type Group struct {
	gorm.Model
	Name      string     `gorm:"unique" json:"name"`
	Processes []*Process `gorm:"many2many:process_groups;"`
	ID        uint32     `gorm:"primary_key" json:"id"`
}

type ProcessConfig struct {
	Processes []ExecFlags `json:"processes"`
}

var ProcessUsageStatsAccessor ProcessUsageStatsAccessorInterface

type ProcessUsageStatsAccessorInterface interface {
	GetUsageStats(int) (string, string)
}
