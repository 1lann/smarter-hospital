package core

import (
	"fmt"
	"log"
	"reflect"
	"runtime/debug"
	"strconv"
)

var registeredLogic = make(map[string]logicInfo)
var moduleToLogic = make(map[string][]string)
var connectHandler = func(moduleID string) {}
var disconnectHandler = func(moduleID string) {}

type logicInfo struct {
	moduleIDs    []string
	logicManager reflect.Value
}

func (l logicInfo) trigger(s *Server) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("core: logic handler for \"" +
				l.logicManager.Type().String() + "\" panic: " + fmt.Sprint(r) +
				"\n" + string(debug.Stack()))
		}
	}()

	arguments := []reflect.Value{reflect.ValueOf(s)}
	for _, moduleID := range l.moduleIDs {
		arguments = append(arguments, setupModules[moduleID].module.Elem())
	}

	l.logicManager.MethodByName("Handle").Call(arguments)
}

// RegisterConnect registers a callback handler whenever a module connects
// (or reconnects) to the system. Used only for servers.
func RegisterConnect(handler func(moduleID string)) {
	connectHandler = handler
}

// RegisterDisconnect registers a callback handler whenever a
// module disconnects from the system. Used only for servers.
func RegisterDisconnect(handler func(moduleID string)) {
	disconnectHandler = handler
}

// RegisterEventLogic registers a logic handler which performs automated
// actions based on events. Only used for servers.
func RegisterEventLogic(moduleIDs []string, logicManager interface{}) {
	managerType := reflect.TypeOf(logicManager)
	if managerType.Kind() != reflect.Ptr ||
		managerType.Elem().Kind() != reflect.Struct {
		panic("core: register event logic: logicManager must be a pointer to " +
			"a struct")
	}

	id := managerType.String()
	panicPrefix := "core: register event logic: " + id + ": "

	handle, exists := managerType.MethodByName("Handle")
	if !exists {
		panic(panicPrefix + "must have a Handle method")
	}

	expectedIn := len(moduleIDs) + 1
	if handle.Type.NumIn()-1 != expectedIn {
		panic(panicPrefix + "expected len(moduleIDs) + 1 = " +
			strconv.Itoa(expectedIn) + " arguments, instead got " +
			strconv.Itoa(handle.Type.NumIn()-1))
	}

	if !reflect.TypeOf(&Server{}).AssignableTo(handle.Type.In(1)) {
		panic(panicPrefix + "first argument must be of of type *core.Server")
	}

	for i := 2; i < len(moduleIDs)+2; i++ {
		module, found := setupModules[moduleIDs[i-2]]
		if !found {
			panic(panicPrefix + "module ID not found: " + moduleIDs[i-2])
		}
		if !handle.Type.In(i).AssignableTo(module.module.Elem().Type()) {
			panic(panicPrefix + "argument " + strconv.Itoa(i) + " must be " +
				"assignable to matching module type: " +
				module.module.Elem().Type().String())
		}
	}

	registeredLogic[id] = logicInfo{
		moduleIDs:    moduleIDs,
		logicManager: reflect.ValueOf(logicManager),
	}

	for _, moduleID := range moduleIDs {
		moduleToLogic[moduleID] = append(moduleToLogic[moduleID], id)
	}
}
