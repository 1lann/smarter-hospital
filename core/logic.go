package core

import (
	"reflect"
	"strconv"
)

var registeredLogic = make(map[string]logicInfo)
var moduleToLogic = make(map[string][]string)
var registeredDisconnect = make(map[string]func())
var registeredConnect = make(map[string]func())

type logicInfo struct {
	moduleIDs    []string
	logicManager reflect.Value
}

func (l logicInfo) trigger() {
	var arguments []reflect.Value
	for _, moduleID := range l.moduleIDs {
		arguments = append(arguments, setupModules[moduleID].module)
	}

	l.logicManager.MethodByName("Handle").Call(arguments)
}

// RegisterConnect registers a callback handler whenever the specified
// module connects (or reconnects) to the system. Only used for servers.
func RegisterConnect(moduleID string, handler func()) {
	if _, found := registeredConnect[moduleID]; found {
		panic("core: register connect: handler already exists for module: " +
			moduleID)
	}

	registeredConnect[moduleID] = handler
}

// RegisterDisconnect registers a callback handler whenever the specified
// module disconnects from the system. Only used for servers.
func RegisterDisconnect(moduleID string, handler func()) {
	if _, found := registeredDisconnect[moduleID]; found {
		panic("core: register disconnect: handler already exists for module: " +
			moduleID)
	}

	registeredConnect[moduleID] = handler
}

// RegisterEventLogic registers a logic handler which performs automated
// actions based on events. Only used for servers.
func RegisterEventLogic(moduleIDs []string, logicManager interface{}) {
	managerType := reflect.TypeOf(logicManager)
	if managerType.Kind() != reflect.Struct {
		panic("core: register event logic: logicManager must be a struct")
	}

	id := managerType.String()
	panicPrefix := "core: register event logic: " + id + ": "

	handle, exists := managerType.MethodByName("Handle")
	if !exists {
		panic(panicPrefix + "must have a Handle method")
	}

	expectedIn := len(moduleIDs) + 1
	if handle.Type.NumIn() != expectedIn {
		panic(panicPrefix + "expected len(moduleIDs) + 1 = " +
			strconv.Itoa(expectedIn) + " arguments, " + "instead got " +
			strconv.Itoa(handle.Type.NumIn()))
	}

	if !reflect.TypeOf(&Server{}).AssignableTo(handle.Type.In(0)) {
		panic(panicPrefix + "first argument must be of of type *core.Server")
	}

	for i := 1; i <= len(moduleIDs); i++ {
		module := setupModules[moduleIDs[i-1]]
		if !handle.Type.In(i).AssignableTo(module.module.Type()) {
			panic(panicPrefix + "argument " + strconv.Itoa(i) + " must be " +
				"assignable to matching module type: " +
				module.module.Elem().Type().String())
		}
	}

	registeredLogic[id] = logicInfo{
		moduleIDs:    moduleIDs,
		logicManager: reflect.ValueOf(&logicManager),
	}

	for _, moduleID := range moduleIDs {
		moduleToLogic[moduleID] = append(moduleToLogic[moduleID], id)
	}
}
