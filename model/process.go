package model

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-set/v2"
	"github.com/joe-at-startupmedia/depgraph"
	"github.com/sirupsen/logrus"
	"os/user"
	"pmon3/protos"
	"pmon3/utils/conv"
	"strings"
	"time"
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
	case StatusBackoff:
		return "backoff"
	case StatusRestarting:
		return "restarting"
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
	case "backoff":
		return StatusBackoff
	case "restarting":
		return StatusRestarting
	}
	return StatusFailed
}

func (p *Process) RenderTable() []string {
	cpuVal, memVal := "0.0%", "0.0 MB"
	if p.CpuUsage != "" {
		cpuVal = p.CpuUsage
	}
	if p.MemoryUsage != "" {
		memVal = p.MemoryUsage
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

func (p *Process) Stringify() string {
	return fmt.Sprintf("%s (%d)", p.Name, p.ID)
}

func (p *Process) Json() (string, error) {
	output, err := json.Marshal(p)
	return string(output), err
}

func (p *Process) GetIdStr() string {
	return conv.Uint32ToStr(p.ID)
}

func (p *Process) GetPidStr() string {
	return conv.Uint32ToStr(p.Pid)
}

func (p *Process) GetRestartCount() uint32 {
	return processRestartCounter[p.ID]
}

func (p *Process) GetRestartCountStr() string {
	return conv.Uint32ToStr(p.RestartCount)
}

func (p *Process) ResetRestartCount() {
	processRestartCounter[p.ID] = 0
}

func (p *Process) IncrRestartCount() {
	processRestartCounter[p.ID] += 1
}

func (p *Process) GetGroupHashSet() *set.HashSet[*Group, string] {
	return set.HashSetFrom[*Group, string](p.Groups)
}

func (p *Process) GetGroupNames() []string {
	groupNames := make([]string, len(p.Groups))
	for i := range p.Groups {
		groupNames[i] = p.Groups[i].Name
	}
	return groupNames
}

func (p *Process) SetUsageStats() {
	if p.Status == StatusRunning {
		p.MemoryUsage, p.CpuUsage = ProcessUsageStatsAccessor.GetUsageStats(int(p.Pid))
	}
}

func (p *Process) ToExecFlags() *ExecFlags {

	flags := ExecFlags{
		File:          p.ProcessFile,
		User:          p.Username,
		Log:           p.Log,
		Args:          p.Args,
		EnvVars:       p.EnvVars,
		Name:          p.Name,
		NoAutoRestart: !p.AutoRestart,
	}

	if len(p.Dependencies) > 0 {
		flags.Dependencies = strings.Split(p.Dependencies, " ")
	}

	if len(p.Groups) > 0 {
		flags.Groups = p.GetGroupNames()
	}

	return &flags
}

//non-receiver methods begin

func FromExecFlags(flags *ExecFlags, logPath string, user *user.User, groups []*Group) *Process {

	var processParams = []string{flags.Name}
	if len(flags.Args) > 0 {
		processParams = append(processParams, strings.Split(flags.Args, " ")...)
	}

	p := Process{
		Pid:          0,
		Log:          logPath,
		Name:         flags.Name,
		ProcessFile:  flags.File,
		Args:         strings.Join(processParams[1:], " "),
		EnvVars:      flags.EnvVars,
		Pointer:      nil,
		Status:       StatusQueued,
		AutoRestart:  !flags.NoAutoRestart,
		Dependencies: strings.Join(flags.Dependencies, " "),
		Groups:       groups,
	}

	if user != nil {
		p.Uid = conv.StrToUint32(user.Uid)
		p.Gid = conv.StrToUint32(user.Gid)
		p.Username = user.Username
	}

	return &p
}

func ComputeDepGraph(processesPtr *[]Process) (*[]Process, *[]Process, error) {

	processes := *processesPtr

	if len(processes) > 1 {
		g := depgraph.New()
		depAppNames := make(map[string]Process)
		nonDepAppNames := make(map[string]Process)
		for _, p := range processes {
			if len(p.Dependencies) > 0 {
				pDependencies := strings.Split(p.Dependencies, " ")
				depAppNames[p.Name] = p
				for _, dep := range pDependencies {
					err := g.DependOn(p.Name, dep)
					if err != nil {
						logrus.Errorf("encountered error building process dependency tree: %s", err)
						return nil, nil, err
					}
				}
			} else {
				nonDepAppNames[p.Name] = p
			}
		}

		if len(g.Leaves()) > 0 {

			dependentProcesses := make([]Process, 0)

			topoSorted := g.TopoSorted()
			for _, processName := range topoSorted {
				if depAppNames[processName].ProcessFile != "" {
					dependentProcesses = append(dependentProcesses, depAppNames[processName])
				} else if nonDepAppNames[processName].ProcessFile != "" {
					dependentProcesses = append(dependentProcesses, nonDepAppNames[processName])
					delete(nonDepAppNames, processName)
				} else if nonDepAppNames[processName].ProcessFile == "" {
					logrus.Warnf("dependencies: %s is not a valid process name", processName)
				}
			}

			nonDependentProcesses := make([]Process, len(nonDepAppNames))
			i := 0
			for pName := range nonDepAppNames {
				nonDependentProcesses[i] = nonDepAppNames[pName]
				i++
			}

			return &nonDependentProcesses, &dependentProcesses, nil
		} else {

			return processesPtr, nil, nil
		}

	}

	return processesPtr, nil, nil
}

func ProcessNames(processesPtr *[]Process) []string {

	if processesPtr == nil {
		return []string{}
	}

	processes := *processesPtr

	if len(processes) == 0 {
		return []string{}
	}

	names := make([]string, len(processes))

	i := 0
	for _, p := range processes {
		names[i] = p.Name
		i++
	}

	return names
}

//protobuf methods begin

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
		EnvVars:      p.EnvVars,
		Status:       p.Status.String(),
		AutoRestart:  p.AutoRestart,
		Uid:          p.Uid,
		Username:     p.Username,
		Gid:          p.Gid,
		RestartCount: p.GetRestartCount(),
		MemoryUsage:  p.MemoryUsage,
		CpuUsage:     p.CpuUsage,
		Dependencies: p.Dependencies,
		Groups:       GroupsArrayToProtobuf(p.Groups),
	}
	return &newProcess
}

func ProcessFromProtobuf(p *protos.Process) *Process {
	createdAt, err := time.Parse(dateTimeFormat, p.GetCreatedAt())
	if err != nil {
		fmt.Println(err)
	}
	updatedAt, err := time.Parse(dateTimeFormat, p.GetUpdatedAt())
	if err != nil {
		fmt.Println(err)
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
		EnvVars:      p.GetEnvVars(),
		Status:       StringToProcessStatus(p.GetStatus()),
		AutoRestart:  p.GetAutoRestart(),
		Uid:          p.GetUid(),
		Username:     p.GetUsername(),
		Gid:          p.GetGid(),
		RestartCount: p.GetRestartCount(),
		MemoryUsage:  p.GetMemoryUsage(),
		CpuUsage:     p.GetCpuUsage(),
		Dependencies: p.GetDependencies(),
		Groups:       GroupsArrayFromProtobuf(p.GetGroups()),
	}
	return &newProcess
}

func GetGroupString(p *protos.Process) string {
	var processNamesStr string
	groupLength := len(p.Groups)
	if groupLength > 0 {
		processNameArray := make([]string, groupLength)
		for i := range p.Groups {
			processNameArray[i] = p.Groups[i].Name
		}
		processNamesStr = strings.Join(processNameArray, ", ")
	}
	return processNamesStr
}
