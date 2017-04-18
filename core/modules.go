package core

import (
	"encoding/gob"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

var registeredModules = make(map[string]registeredModule)
var setupModules = make(map[string]*setupModule)

type registeredModule struct {
	eventType        reflect.Type
	actionType       reflect.Type
	settingsType     reflect.Type
	moduleType       reflect.Type
	hasEventHandler  bool
	hasInfoProvider  bool
	hasActionHandler bool
}

type setupModule struct {
	registration registeredModule
	module       reflect.Value
}

// RegisterModule registers a module
func RegisterModule(module interface{}) {
	moduleType := reflect.TypeOf(module)
	moduleName := strings.Split(moduleType.String(), ".")[0]
	panicPrefix := "comm: register module: " + moduleName + ": "

	if moduleType.Kind() != reflect.Struct {
		panic(panicPrefix + "module must be a struct")
	}

	register := registeredModule{moduleType: moduleType}

	eventField, exists := moduleType.FieldByName("Event")
	if exists {
		register.eventType = eventField.Type
		gob.RegisterName("evt_"+moduleName, reflect.New(eventField.Type))
	}

	actionField, exists := moduleType.FieldByName("Action")
	if exists {
		register.actionType = actionField.Type
		gob.RegisterName("act_"+moduleName, reflect.New(actionField.Type))
	}

	settingsField, exists := moduleType.FieldByName("Settings")
	if exists {
		register.settingsType = settingsField.Type
	}

	eventHandler, exists := moduleType.MethodByName("HandleEvent")
	if exists {
		if eventHandler.Type.NumIn() != 1 {
			panic(panicPrefix + "HandleEvent: expected 1 argument, " +
				"instead got " + strconv.Itoa(eventHandler.Type.NumIn()))
		}

		if !register.eventType.AssignableTo(eventHandler.Type.In(0)) {
			panic(panicPrefix + "HandleEvent: first argument must be the " +
				"same type as Module.Event")
		}

		if eventHandler.Type.NumOut() != 1 {
			panic(panicPrefix + "HandleEvent: expected 1 return value, " +
				"instead got " + strconv.Itoa(eventHandler.Type.NumOut()))
		}

		if !eventHandler.Type.Out(0).Implements(reflect.TypeOf((*error)(nil)).
			Elem()) {
			panic(panicPrefix + "HandleEvent: return type must be error")
		}

		register.hasEventHandler = true
	}

	infoProvider, exists := moduleType.MethodByName("Info")
	if exists {
		if infoProvider.Type.NumIn() != 0 {
			panic(panicPrefix + "Info: expected 0 arguments, instead " +
				"got " + strconv.Itoa(infoProvider.Type.NumIn()))
		}

		if infoProvider.Type.NumOut() != 1 {
			panic(panicPrefix + "Info: expected 1 return argument, " +
				"instead got " + strconv.Itoa(infoProvider.Type.NumOut()))
		}

		register.hasInfoProvider = true
	}

	actionHandler, exists := moduleType.MethodByName("HandleAction")
	if exists {
		if actionHandler.Type.NumIn() != 1 {
			panic(panicPrefix + "HandleAction: expected 1 argument, " +
				"instead got " + strconv.Itoa(actionHandler.Type.NumIn()))
		}

		if !register.actionType.AssignableTo(actionHandler.Type.In(0)) {
			panic(panicPrefix + "HandleAction: first argument must be the " +
				"same type as Module.Action")
		}

		register.hasActionHandler = true
	}

	if _, found := registeredModules[moduleName]; found {
		panic(panicPrefix + "module already registered")
	}
	registeredModules[moduleName] = register
}

// SetupModule sets up a module for use with the given settings.
func SetupModule(moduleName string, id string, settings ...interface{}) {
	if _, found := setupModules[id]; found {
		panic("comm: setup module: module ID already exists: " + id)
	}

	panicPrefix := "comm: setup module: " + id + ": "

	if len(settings) >= 1 {
		panic(panicPrefix + "at most 1 settings value can be provided")
	}

	module, found := registeredModules[moduleName]
	if !found {
		panic(panicPrefix + "setup module: module not found: " + moduleName)
	}

	newModule := reflect.New(module.moduleType)

	if module.settingsType != nil && len(settings) > 0 {
		if !reflect.TypeOf(settings[0]).AssignableTo(module.settingsType) {
			panic(panicPrefix + "provided settings must be assignable to " +
				"Module.Settings")
		}

		newModule.FieldByName("Settings").Set(reflect.ValueOf(settings[0]))
	}

	setupModules[id] = &setupModule{
		registration: module,
		module:       newModule,
	}
}

// ErrNoSuchAction is returned if no such action value is registered when
// calling ActionValue.
var ErrNoSuchAction = errors.New("comm: no such action")

// ActionValue returns the registered sample value given the ID of a module.
func ActionValue(id string) (interface{}, error) {
	value, found := setupModules[id]
	if !found {
		return nil, ErrNoSuchAction
	}

	if value.registration.actionType == nil {
		return nil, ErrNoSuchAction
	}

	return reflect.New(value.registration.actionType).Interface(), nil
}
