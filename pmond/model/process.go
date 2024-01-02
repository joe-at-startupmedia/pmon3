package model

import (
	"encoding/json"
	"fmt"
	"os"
	"pmon3/pmond/utils/cpu"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

type ProcessStatus int64

const (
	StatusQueued ProcessStatus = iota
	StatusInit
	StatusRunning
	StatusStopped
	StatusFailed
	StatusClosed
)

func (s ProcessStatus) String() string {
	switch s {
	case StatusQueued:
		return "queued"
	case StatusInit:
		return "init"
	case StatusRunning:
		return "running"
	case StatusStopped:
		return "stopped"
	case StatusFailed:
		return "failed"
	case StatusClosed:
		return "closed"
	}
	return "unknown"
}

type Process struct {
	ID          uint          `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Pid         int           `gorm:"column:pid" json:"pid"`
	Log         string        `gorm:"column:log" json:"log"`
	Name        string        `gorm:"unique" json:"name"`
	ProcessFile string        `json:"process_file"`
	Args        string        `json:"args"`
	Status      ProcessStatus `json:"status"`
	Pointer     *os.Process   `gorm:"-" json:"-"`
	AutoRestart bool          `json:"auto_restart"`
	Uid         string
	Username    string
	Gid         string
}

func (p Process) NoAutoRestartStr() string {
	if !p.AutoRestart {
		return "true"
	} else {
		return "false"
	}
}

func (Process) TableName() string {
	return "process"
}

func (p Process) MustJson() string {
	data, _ := json.Marshal(&p)

	return string(data)
}

func (p Process) RenderTable() []string {
	cpuVal, memVal := "0", "0"
	if p.Status == StatusRunning {
		cpuVal, memVal = cpu.GetExtraInfo(p.Pid)
	}

	return []string{
		strconv.Itoa(int(p.ID)),
		p.Name,
		strconv.Itoa(p.Pid),
		p.Status.String(),
		p.Username,
		cpuVal,
		memVal,
		p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func FindByProcessFileAndName(db *gorm.DB, process_file string, name string) (error, *Process) {
	var process Process
	err := db.First(&process, "process_file = ? AND name = ?", process_file, name).Error
	if err != nil {
		return err, nil
	}

	return nil, &process
}

func (p Process) Stringify() string {
	return fmt.Sprintf("%s (%d)", p.Name, p.ID)
}
