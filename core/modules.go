package core

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
)

// ErrNotFound is returned if a value could not be found.
var ErrNotFound = errors.New("core: not found")

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
	hasPollEvents    bool
	hasInitHandler   bool
}

type setupModule struct {
	registration registeredModule
	module       reflect.Value
}

// RegisterModule registers a module. Used for clients and servers.
func RegisterModule(module interface{}) {
	moduleType := reflect.TypeOf(module)
	moduleValue := reflect.New(moduleType)
	moduleName := strings.Split(moduleType.String(), ".")[0]
	panicPrefix := "core: register module: " + moduleName + ": "

	if moduleType.Kind() != reflect.Struct {
		panic(panicPrefix + "module must be a struct")
	}

	idField, exists := moduleType.FieldByName("ID")
	if !exists {
		panic(panicPrefix + "module must have an ID field")
	}

	if idField.Type.Kind() != reflect.String {
		panic(panicPrefix + "ID field must be of type string")
	}

	register := registeredModule{moduleType: moduleType}

	settingsField, exists := moduleType.FieldByName("Settings")
	if exists {
		register.settingsType = settingsField.Type
		gob.RegisterName("set_"+moduleName,
			reflect.Zero(settingsField.Type).Interface())
	}

	eventHandler := moduleValue.MethodByName("HandleEvent")
	if eventHandler.IsValid() {
		handlerType := eventHandler.Type()
		if handlerType.NumIn() != 1 {
			panic(panicPrefix + "HandleEvent: expected 1 input argument, " +
				"instead got " + strconv.Itoa(handlerType.NumIn()))
		}

		register.eventType = handlerType.In(0)
		gob.RegisterName("evt_"+moduleName, reflect.Zero(register.eventType).
			Interface())

		if handlerType.NumOut() != 0 {
			panic(panicPrefix + "HandleEvent: expected no return values, " +
				"instead got " + strconv.Itoa(handlerType.NumOut()))
		}

		register.hasEventHandler = true
	}

	infoProvider := moduleValue.MethodByName("Info")
	if infoProvider.IsValid() {
		infoType := infoProvider.Type()
		if infoType.NumIn() != 0 {
			panic(panicPrefix + "Info: expected 0 arguments, instead " +
				"got " + strconv.Itoa(infoType.NumIn()))
		}

		if infoType.NumOut() != 1 {
			panic(panicPrefix + "Info: expected 1 return argument, " +
				"instead got " + strconv.Itoa(infoType.NumOut()))
		}

		register.hasInfoProvider = true
	}

	actionHandler := moduleValue.MethodByName("HandleAction")
	if actionHandler.IsValid() {
		handlerType := actionHandler.Type()
		if handlerType.NumIn() != 2 {
			panic(panicPrefix + "HandleAction: expected 2 arguments, " +
				"instead got " + strconv.Itoa(handlerType.NumIn()))
		}

		if !reflect.TypeOf(&Client{}).AssignableTo(handlerType.In(0)) {
			panic(panicPrefix + "HandleAction: first argument must be of " +
				"type *core.Client")
		}

		register.actionType = handlerType.In(1)
		gob.RegisterName("act_"+moduleName, reflect.Zero(register.actionType).
			Interface())

		if handlerType.NumOut() != 1 {
			panic(panicPrefix + "HandleAction: expected 1 return argument, " +
				"instead got " + strconv.Itoa(handlerType.NumOut()))
		}

		if !handlerType.Out(0).Implements(reflect.TypeOf((*error)(nil)).
			Elem()) {
			panic(panicPrefix + "HandleAction: return type must be error")
		}

		register.hasActionHandler = true
	}

	pollEvents := moduleValue.MethodByName("PollEvents")
	if pollEvents.IsValid() {
		handlerType := pollEvents.Type()
		if handlerType.NumIn() != 1 {
			panic(panicPrefix + "PollEvents: expected 1 input argument, " +
				"instead got " + strconv.Itoa(handlerType.NumIn()))
		}

		if !reflect.TypeOf(&Client{}).AssignableTo(handlerType.In(0)) {
			panic(panicPrefix + "PollEvents: first argument must be of " +
				"type *core.Client")
		}

		if handlerType.NumOut() != 0 {
			panic(panicPrefix + "PollEvents: expected no return arguments, " +
				"instead got " + strconv.Itoa(handlerType.NumOut()))
		}

		register.hasPollEvents = true
	}

	initHandler := moduleValue.MethodByName("Init")
	if initHandler.IsValid() {
		handlerType := initHandler.Type()
		if handlerType.NumIn() != 0 {
			panic(panicPrefix + "Init: expected no input arguments, " +
				"instead got " + strconv.Itoa(handlerType.NumIn()))
		}

		if handlerType.NumOut() != 0 {
			panic(panicPrefix + "Init: expected no return arguments, " +
				"instead got " + strconv.Itoa(handlerType.NumOut()))
		}

		register.hasInitHandler = true
	}

	if _, found := registeredModules[moduleName]; found {
		panic(panicPrefix + "module already registered")
	}
	registeredModules[moduleName] = register
}

// SetupModule sets up a module for use with the given settings.
// Used for clients and servers.
func SetupModule(moduleName string, id string, settings ...interface{}) {
	if _, found := setupModules[id]; found {
		panic("core: setup module: module ID already exists: " + id)
	}

	panicPrefix := "core: setup module: " + id + ": "

	if len(settings) > 1 {
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

		newModule.Elem().FieldByName("Settings").
			Set(reflect.ValueOf(settings[0]))
	}

	if module.hasInitHandler {
		newModule.MethodByName("Init").Call([]reflect.Value{})
	}

	newModule.Elem().FieldByName("ID").SetString(id)

	setupModules[id] = &setupModule{
		registration: module,
		module:       newModule,
	}
}

// ActionType returns the Action type given the ID of a module.
func ActionType(id string) (reflect.Type, error) {
	value, found := setupModules[id]
	if !found {
		return nil, ErrNotFound
	}

	if value.registration.actionType == nil {
		return nil, ErrNotFound
	}

	return value.registration.actionType, nil
}

// Info returns the information value given the ID of a module.
func Info(id string) (result interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("core: info panic: " +
				fmt.Sprint(r) + "\n" + string(debug.Stack()))
			err = errors.New(fmt.Sprint(r))
		}

	}()

	value, found := setupModules[id]
	if !found {
		return nil, ErrNotFound
	}

	if !value.registration.hasInfoProvider {
		return nil, ErrNotFound
	}

	callResult := value.module.MethodByName("Info").Call([]reflect.Value{})
	return callResult[0].Interface(), nil
}
