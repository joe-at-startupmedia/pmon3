package controller

import (
	"fmt"
	"github.com/joe-at-startupmedia/depgraph"
	"pmon3/conf"
	"pmon3/pmond"
	"pmon3/pmond/db"
	"pmon3/pmond/model"
	"pmon3/pmond/protos"
	"strings"
)

func Dgraph(cmd *protos.Cmd) *protos.CmdResp {
	queueOrder, dGraph, err := computeGraph()
	if err != nil {
		return ErroredCmdResp(cmd, fmt.Errorf("command error: could not get graph: %w", err))
	}
	output1 := strings.Join(queueOrder, "\n")
	output2 := strings.Join(dGraph, "\n")
	newCmdResp := protos.CmdResp{
		Id:       cmd.GetId(),
		Name:     cmd.GetName(),
		ValueStr: fmt.Sprintf("%s||%s", output1, output2),
	}
	return &newCmdResp
}

func computeGraph() ([]string, []string, error) {

	var all []model.Process
	err := db.Db().Find(&all).Error
	apps := pmond.Config.AppsConfig.Apps

	enqueueOrder := make([]string, 0)

	if len(apps) > 0 {
		g := depgraph.New()
		depAppNames := make(map[string]conf.AppsConfigApp)
		nonDepAppNames := make(map[string]conf.AppsConfigApp)
		for _, app := range apps {
			if len(app.Flags.Dependencies) > 0 {
				depAppNames[app.Flags.Name] = app
				for _, dep := range app.Flags.Dependencies {
					err = g.DependOn(app.Flags.Name, dep)
					if err != nil {
						pmond.Log.Errorf("encountered error building app dependency tree: %s", err)
					}
				}
			} else {
				nonDepAppNames[app.Flags.Name] = app
			}
		}

		if len(g.Leaves()) > 0 {

			sortedLayers := make([]string, 0)

			topoSortedLayers := g.TopoSortedLayers()
			for _, appNames := range topoSortedLayers {
				for _, appName := range appNames {
					if depAppNames[appName].File != "" {
						sortedLayers = append(sortedLayers, appName)
					} else if nonDepAppNames[appName].File != "" {
						sortedLayers = append(sortedLayers, appName)
						nonDepAppNames[appName] = conf.AppsConfigApp{}
					} else if nonDepAppNames[appName].File == "" {
						pmond.Log.Warnf("dependencies: %s is not a valid app name", appName)
					}
				}
			}

			enqueueOrder = append(enqueueOrder, sortedLayers...)
			for appName := range nonDepAppNames {
				if nonDepAppNames[appName].File != "" {
					enqueueOrder = append(enqueueOrder, appName)
				}
			}

			return enqueueOrder, sortedLayers, nil
		} else {
			for _, app := range apps {
				enqueueOrder = append(enqueueOrder, app.Flags.Name)
			}

			return enqueueOrder, nil, nil
		}

	}

	return nil, nil, nil
}
