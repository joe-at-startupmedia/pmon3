package model

type AppsConfigApp struct {
	File  string    `json:"file"`
	Flags ExecFlags `json:"flags"`
}

func AppsConfigAppNames(appMapPtr *[]AppsConfigApp) []string {

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
