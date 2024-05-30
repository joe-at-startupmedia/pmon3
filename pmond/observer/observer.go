package observer

import (
	"fmt"
	"github.com/goinbox/shell"
	"pmon3/pmond"
	"pmon3/pmond/model"
	"strings"
)

type EventType int

const (
	RestartEvent EventType = iota + 1
	FailedEvent
)

func (w EventType) String() string {
	return [...]string{"restarted", "failed"}[w-1]
}

// Event with some payload
type Event struct {
	Process *model.Process
	Type    EventType
}

func HandleEvent(evt *Event) {
	pmond.Log.Infof("Received event: %v\n", evt)
	switch evt.Type {
	case FailedEvent:
		onFailedEvent(evt)
	case RestartEvent:
		onRestartEvent(evt)
	}
}

func onFailedEvent(evt *Event) {
	pmond.Log.Errorf("process failed and not restarting: %s", evt.Process.Stringify())
	oPE := pmond.Config.OnProcessFailureExec
	if len(oPE) > 0 {
		onEventExec(evt, oPE)
	}
}

func onRestartEvent(evt *Event) {
	pmond.Log.Warnf("restarting process: %s", evt.Process.Stringify())
	oPE := pmond.Config.OnProcessRestartExec
	if len(oPE) > 0 {
		onEventExec(evt, oPE)
	}
}

func jsonEscape(rawJson string) string {
	//remove any existing backslashes and then escape double quotes with backslashes
	replacer := strings.NewReplacer("\\", "", "\"", "\\\"")
	return replacer.Replace(rawJson)
}

func onEventExec(evt *Event, cmd string) {
	processJson, err := evt.Process.Json()
	if err != nil {
		pmond.Log.Errorf("Could not json-serialize process for event executor(%s): %s", evt.Type.String(), err)
		return
	}
	execCmd := fmt.Sprintf("%s \"%s\"", cmd, jsonEscape(processJson))
	pmond.Log.Debugf("Attempting event executor(%s): %s", evt.Type.String(), execCmd)
	rel := shell.RunCmd(execCmd)
	if !rel.Ok {
		pmond.Log.Errorf("event executor encountered an err: %s", rel.Output)
	}
}
