package model

import (
	"encoding/json"
	"fmt"
	"os"
	"pmon3/pmond/protos"
	"pmon3/pmond/utils/conv"
	"pmon3/pmond/utils/cpu"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ProcessStatus int64

const dateTimeFormat = "2006-01-02 15:04:05"

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

func StringToProcessStatus(s string) ProcessStatus {
	switch s {
	case "queued":
		return StatusQueued
	case "init":
		return StatusInit
	case "running":
		return StatusRunning
	case "stopped":
		return StatusStopped
	case "failed":
		return StatusFailed
	case "closed":
		return StatusClosed
	}
	return StatusFailed
}

type Process struct {
	ID           uint32        `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Pid          uint32        `gorm:"column:pid" json:"pid"`
	Log          string        `gorm:"column:log" json:"log"`
	Name         string        `gorm:"unique" json:"name"`
	ProcessFile  string        `json:"process_file"`
	Args         string        `json:"args"`
	Status       ProcessStatus `json:"status"`
	Pointer      *os.Process   `gorm:"-" json:"-"`
	AutoRestart  bool          `json:"auto_restart"`
	Uid          uint32        `gorm:"column:uid" json:"uid"`
	Username     string        `json:"username"`
	Gid          uint32        `gorm:"column:gid" json:"gid"`
	RestartCount uint32        `gorm:"-" json:"-"`
}

func (p Process) NoAutoRestartStr() string {
	return strconv.FormatBool(!p.AutoRestart)
}

func (Process) TableName() string {
	return "process"
}

func (p *Process) RenderTable() []string {
	cpuVal, memVal := "0%", "0.0 MB"
	if p.Status == StatusRunning {
		cpuVal, memVal = cpu.GetExtraInfo(int(p.Pid))
	}

	return []string{
		p.GetIdStr(),
		p.Name,
		p.GetPidStr(),
		p.GetRestartCountStr(),
		p.Status.String(),
		p.Username,
		cpuVal,
		memVal,
		p.UpdatedAt.Format(dateTimeFormat),
	}
}

func FindProcessByFileAndName(db *gorm.DB, process_file string, name string) (error, *Process) {
	var process Process
	err := db.First(&process, "process_file = ? AND name = ?", process_file, name).Error
	if err != nil {
		return err, nil
	}

	return nil, &process
}

func FindProcessByIdOrName(db *gorm.DB, idOrName string) (error, *Process) {
	var process Process
	err := db.First(&process, "id = ? or name = ?", idOrName, idOrName).Error
	if err != nil {
		return err, nil
	}

	return nil, &process
}

func (process *Process) Save(db *gorm.DB) (string, error) {
	err, originOne := FindProcessByFileAndName(db, process.ProcessFile, process.Name)
	if err == nil && originOne.ID > 0 { // process already exists
		process.ID = originOne.ID
	}

	err = db.Save(&process).Error
	if err != nil {
		return "", fmt.Errorf("pmon3 run err: %w", err)
	}

	output, err := json.Marshal(process.RenderTable())
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func (p Process) Stringify() string {
	return fmt.Sprintf("%s (%d)", p.Name, p.ID)
}

func (p *Process) GetIdStr() string {
	return conv.Uint32ToStr(p.ID)
}

func (p *Process) GetPidStr() string {
	return conv.Uint32ToStr(p.Pid)
}

var restartCount = make(map[uint32]uint32)

func (p *Process) GetRestartCount() uint32 {
	return restartCount[p.ID]
}

func (p *Process) GetRestartCountStr() string {
	return conv.Uint32ToStr(p.RestartCount)
}

func (p *Process) ResetRestartCount() {
	restartCount[p.ID] = 0
}

func (p *Process) IncrRestartCount() {
	restartCount[p.ID] += 1
}

func (p *Process) ToProtobuf() *protos.Process {
	newProcess := protos.Process{
		Id:           p.ID,
		CreatedAt:    p.CreatedAt.Format(dateTimeFormat),
		UpdatedAt:    p.UpdatedAt.Format(dateTimeFormat),
		Pid:          p.Pid,
		Log:          p.Log,
		Name:         p.Name,
		ProcessFile:  p.ProcessFile,
		Args:         p.Args,
		Status:       p.Status.String(),
		AutoRestart:  p.AutoRestart,
		Uid:          p.Uid,
		Username:     p.Username,
		Gid:          p.Gid,
		RestartCount: p.GetRestartCount(),
	}
	return &newProcess
}

func FromProtobuf(p *protos.Process) *Process {
	createdAt, error := time.Parse(dateTimeFormat, p.GetCreatedAt())
	if error != nil {
		fmt.Println(error)
	}
	updatedAt, error := time.Parse(dateTimeFormat, p.GetUpdatedAt())
	if error != nil {
		fmt.Println(error)
	}
	newProcess := Process{
		ID:           p.GetId(),
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		Pid:          p.GetPid(),
		Log:          p.GetLog(),
		Name:         p.GetName(),
		ProcessFile:  p.GetProcessFile(),
		Args:         p.GetArgs(),
		Status:       StringToProcessStatus(p.GetStatus()),
		AutoRestart:  p.GetAutoRestart(),
		Uid:          p.GetUid(),
		Username:     p.GetUsername(),
		Gid:          p.GetGid(),
		RestartCount: p.GetRestartCount(),
	}
	return &newProcess
}
