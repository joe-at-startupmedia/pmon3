package conf

import (
	"github.com/joe-at-startupmedia/depgraph"
	"github.com/sirupsen/logrus"
)

func ComputeDepGraph(apps []AppsConfigApp) ([]AppsConfigApp, []AppsConfigApp, error) {

	if len(apps) > 0 {
		g := depgraph.New()
		enqueueOrder := make([]AppsConfigApp, 0)
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

			sortedLayers := make([]AppsConfigApp, 0)

			topoSortedLayers := g.TopoSortedLayers()
			for _, appNames := range topoSortedLayers {
				for _, appName := range appNames {
					if depAppNames[appName].File != "" {
						sortedLayers = append(sortedLayers, depAppNames[appName])
						enqueueOrder = append(enqueueOrder, depAppNames[appName])
					} else if nonDepAppNames[appName].File != "" {
						sortedLayers = append(sortedLayers, nonDepAppNames[appName])
						enqueueOrder = append(enqueueOrder, nonDepAppNames[appName])
						delete(nonDepAppNames, appName)
					} else if nonDepAppNames[appName].File == "" {
						logrus.Warnf("dependencies: %s is not a valid app name", appName)
					}
				}
			}

			for appName := range nonDepAppNames {
				enqueueOrder = append(enqueueOrder, nonDepAppNames[appName])
			}

			return enqueueOrder, sortedLayers, nil
		} else {

			return apps, nil, nil
		}

	}

	return nil, nil, nil
}

func MapKeys(appMap []AppsConfigApp) []string {

	keys := make([]string, len(appMap))

	i := 0
	for _, app := range appMap {
		keys[i] = app.Flags.Name
		i++
	}

	return keys
}
