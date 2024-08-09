package conf

import (
	"fmt"
	"github.com/joe-at-startupmedia/depgraph"
	"github.com/sirupsen/logrus"
	"pmon3/pmond/model"
)

type AppsConfig struct {
	Apps []AppsConfigApp
}

type AppsConfigApp struct {
	File  string
	Flags model.ExecFlags
}

func ComputeDepGraph(appsPtr *[]AppsConfigApp) (*[]AppsConfigApp, *[]AppsConfigApp, error) {

	apps := *appsPtr

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

			return appsPtr, nil, nil
		}

	}

	return appsPtr, nil, nil
}

func AppNames(appMapPtr *[]AppsConfigApp) []string {

	if appMapPtr == nil {
		return []string{}
	}

	appMap := *appMapPtr

	if len(appMap) == 0 {
		return []string{}
	}

	keys := make([]string, len(appMap))

	i := 0
	for _, app := range appMap {
		keys[i] = app.Flags.Name
		i++
	}

	return keys
}

func GetAppByName(appName string, apps []AppsConfigApp) (AppsConfigApp, error) {
	for _, app := range apps {
		if app.Flags.Name == appName {
			return app, nil
		}
	}
	return AppsConfigApp{}, fmt.Errorf("could not find app in Apps Config with name %s", appName)
}
