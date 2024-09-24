package model

import (
	"encoding/json"
	"fmt"
	"github.com/joe-at-startupmedia/depgraph"
	"github.com/sirupsen/logrus"
)

type AppsConfig struct {
	Apps []AppsConfigApp `json:"apps"`
}

func (ac *AppsConfig) Json() string {
	content, _ := json.Marshal(ac)
	return string(content)
}

func (ac *AppsConfig) ComputeDepGraph() (*[]AppsConfigApp, *[]AppsConfigApp, error) {

	apps := ac.Apps

	if len(apps) > 1 {
		g := depgraph.New()
		depAppNames := make(map[string]AppsConfigApp)
		nonDepAppNames := make(map[string]AppsConfigApp)
		for _, app := range apps {
			if len(app.Flags.Dependencies) > 0 {
				depAppNames[app.Flags.Name] = app
				for _, dep := range app.Flags.Dependencies {
					err := g.DependOn(app.Flags.Name, dep)
					if err != nil {
						logrus.Errorf("encountered error building app dependency tree: %s", err)
						return nil, nil, err
					}
				}
			} else {
				nonDepAppNames[app.Flags.Name] = app
			}
		}

		if len(g.Leaves()) > 0 {

			dependentApps := make([]AppsConfigApp, 0)

			topoSorted := g.TopoSorted()
			for _, appName := range topoSorted {
				if depAppNames[appName].File != "" {
					dependentApps = append(dependentApps, depAppNames[appName])
				} else if nonDepAppNames[appName].File != "" {
					dependentApps = append(dependentApps, nonDepAppNames[appName])
					delete(nonDepAppNames, appName)
				} else if nonDepAppNames[appName].File == "" {
					logrus.Warnf("dependencies: %s is not a valid app name", appName)
				}
			}

			nonDependentApps := make([]AppsConfigApp, len(nonDepAppNames))
			i := 0
			for appName := range nonDepAppNames {
				nonDependentApps[i] = nonDepAppNames[appName]
				i++
			}

			return &nonDependentApps, &dependentApps, nil
		} else {

			return &ac.Apps, nil, nil
		}

	}

	return &ac.Apps, nil, nil
}

func (ac *AppsConfig) GetAppByName(appName string) (AppsConfigApp, error) {
	for _, app := range ac.Apps {
		if app.Flags.Name == appName {
			return app, nil
		}
	}
	return AppsConfigApp{}, fmt.Errorf("could not find app in Apps Config with name %s", appName)
}
