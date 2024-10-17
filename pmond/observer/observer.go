package observer

import (
	"fmt"
	"pmon3/model"
	"pmon3/pmond"
	"pmon3/pmond/shell"
	"strings"
)

type EventType int

var OnRestartEventFunc func(evt *Event)
var OnFailedEventFunc func(evt *Event)
var OnBackOffEventFunc func(evt *Event)

//*runtime.Func

const (
	RestartEvent EventType = iota + 1
	FailedEvent
	BackoffEvent
)

func (w EventType) String() string {
	return [...]string{"restarted", "failed", "backoff"}[w-1]
}

type Event struct {
	Process *model.Process
	Type    EventType
}

func HandleEvent(evt *Event) {
	pmond.Log.Infof("Observer: Received event: %v\n", evt)
	switch evt.Type {
	case FailedEvent:
		onFailedEvent(evt)
	case RestartEvent:
		onRestartEvent(evt)
	case BackoffEvent:
		onBackoffEvent(evt)
	}
}

func onFailedEvent(evt *Event) {
	pmond.Log.Warnf("Observer: process failed and not restarting: %s", evt.Process.Stringify())
	oPE := pmond.Config.EventHandler.ProcessFailure
	if len(oPE) > 0 {
		onEventExec(evt, oPE)
	}
	if OnFailedEventFunc != nil {
		OnFailedEventFunc(evt)
	}
}

func onRestartEvent(evt *Event) {
	pmond.Log.Warnf("Observer: restarting process: %s", evt.Process.Stringify())
	oPE := pmond.Config.EventHandler.ProcessRestart
	if len(oPE) > 0 {
		onEventExec(evt, oPE)
	}
	if OnRestartEventFunc != nil {
		OnRestartEventFunc(evt)
	}
}

func onBackoffEvent(evt *Event) {
	pmond.Log.Warnf("Observer: process backed off and not restarting: %s", evt.Process.Stringify())
	oPE := pmond.Config.EventHandler.ProcessBackoff
	if len(oPE) > 0 {
		onEventExec(evt, oPE)
	}
	if OnBackOffEventFunc != nil {
		OnBackOffEventFunc(evt)
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
	pmond.Log.Debugf("Observer: Attempting event executor(%s): %s", evt.Type.String(), execCmd)
	shell.HandleOnEventExec(execCmd)
}
