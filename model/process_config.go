package model

import (
	"encoding/json"
	"fmt"
	"github.com/joe-at-startupmedia/depgraph"
	"github.com/sirupsen/logrus"
)

type ProcessConfig struct {
	Processes []ExecFlags `json:"processes"`
}

func (ac *ProcessConfig) Json() string {
	content, _ := json.Marshal(ac)
	return string(content)
}

func (ac *ProcessConfig) ComputeDepGraph() (*[]ExecFlags, *[]ExecFlags, error) {

	processes := ac.Processes

	if len(processes) > 1 {
		g := depgraph.New()
		depProcessNames := make(map[string]ExecFlags)
		nonDepProcessNames := make(map[string]ExecFlags)
		for _, p := range processes {
			if len(p.Dependencies) > 0 {
				depProcessNames[p.Name] = p
				for _, dep := range p.Dependencies {
					err := g.DependOn(p.Name, dep)
					if err != nil {
						logrus.Errorf("encountered error building process dependency tree: %s", err)
						return nil, nil, err
					}
				}
			} else {
				nonDepProcessNames[p.Name] = p
			}
		}

		if len(g.Leaves()) > 0 {

			dependentProcesses := make([]ExecFlags, 0)

			topoSorted := g.TopoSorted()
			for _, processName := range topoSorted {
				if depProcessNames[processName].File != "" {
					dependentProcesses = append(dependentProcesses, depProcessNames[processName])
				} else if nonDepProcessNames[processName].File != "" {
					dependentProcesses = append(dependentProcesses, nonDepProcessNames[processName])
					delete(nonDepProcessNames, processName)
				} else if nonDepProcessNames[processName].File == "" {
					logrus.Warnf("dependencies: %s is not a valid process name", processName)
				}
			}

			nonDependentProcesses := make([]ExecFlags, len(nonDepProcessNames))
			i := 0
			for pName := range nonDepProcessNames {
				nonDependentProcesses[i] = nonDepProcessNames[pName]
				i++
			}

			return &nonDependentProcesses, &dependentProcesses, nil
		} else {

			return &ac.Processes, nil, nil
		}

	}

	return &ac.Processes, nil, nil
}

func (ac *ProcessConfig) GetExecFlagsByName(name string) (ExecFlags, error) {
	for _, execFlags := range ac.Processes {
		if execFlags.Name == name {
			return execFlags, nil
		}
	}
	return ExecFlags{}, fmt.Errorf("could not find process in Process Config with name %s", name)
}
