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

func ComputeDepGraph(apps []AppsConfigApp) (*[]AppsConfigApp, *[]AppsConfigApp, error) {

	if len(apps) > 0 {
		g := depgraph.New()
		nonDependentApps := make([]AppsConfigApp, 0)
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

			topoSortedLayers := g.TopoSortedLayers()
			for _, appNames := range topoSortedLayers {
				for _, appName := range appNames {
					if depAppNames[appName].File != "" {
						dependentApps = append(dependentApps, depAppNames[appName])
					} else if nonDepAppNames[appName].File != "" {
						dependentApps = append(dependentApps, nonDepAppNames[appName])
						delete(nonDepAppNames, appName)
					} else if nonDepAppNames[appName].File == "" {
						logrus.Warnf("dependencies: %s is not a valid app name", appName)
					}
				}
			}

			for appName := range nonDepAppNames {
				nonDependentApps = append(nonDependentApps, nonDepAppNames[appName])
			}

			return &nonDependentApps, &dependentApps, nil
		} else {

			return &apps, nil, nil
		}

	}

	return nil, nil, nil
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
