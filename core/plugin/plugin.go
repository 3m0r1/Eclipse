package plugin

import (
	"errors"

	lua "github.com/yuin/gopher-lua"
)

type PluginInfo struct {
	Name        string
	Description string
	Author      string
	Version     string
	Path        string
}

type Environment struct {
	Plugin *lua.LTable
}

type Plugin struct {
	Info        *PluginInfo
	Events      map[string]*Event
	Commands    map[string]*Command
	Imports     map[string]map[string]*Import
	Environment *Environment
	State       *lua.LState
}

func NewPlugin(path string) *Plugin {
	plugin := Plugin{
		Info:        &PluginInfo{Path: path},
		Events:      map[string]*Event{},
		Commands:    map[string]*Command{},
		Imports:     map[string]map[string]*Import{},
		Environment: &Environment{},
		State:       lua.NewState(),
	}

	return &plugin
}

func (plugin *Plugin) GetEvent(name string) (*Event, bool) {
	event, exists := plugin.Events[name]
	return event, exists
}

func (plugin *Plugin) AddEvent(name string, event *Event) error {
	_, exists := plugin.GetEvent(name)

	if exists {
		return errors.New("event already exists")
	}

	plugin.Events[name] = event
	return nil
}

func (plugin *Plugin) FireEvent(name string, args ...lua.LValue) error {
	event, ok := plugin.Events[name]

	if !ok {
		return errors.New("event doesn't exist")
	} else {
		event.Fire(
			plugin.State,
			args...,
		)
		return nil
	}
}

func (plugin *Plugin) GetCommand(name string) (*Command, bool) {
	command, exists := plugin.Commands[name]
	return command, exists
}

func (plugin *Plugin) AddCommand(name string, command *Command) error {
	_, exists := plugin.GetCommand(name)

	if exists {
		return errors.New("command already exists")
	}

	plugin.Commands[name] = command
	return nil
}

func (plugin *Plugin) InvokeCommandCtx(name string, ctx lua.LValue, args ...lua.LValue) (*lua.LTable, error) {
	cmd := plugin.Commands[name]
	return cmd.Invoke(plugin.State, ctx, args...)
}

func (plugin *Plugin) InvokeCommand(name string, args ...lua.LValue) (*lua.LTable, error) {
	cmd := plugin.Commands[name]
	ctx := plugin.Environment.Plugin
	return cmd.Invoke(plugin.State, ctx, args...)
}

func (plugin *Plugin) AddImport(name, proc string, imp *Import) {
	_, ok := plugin.Imports[name]

	if !ok {
		plugin.Imports[name] = map[string]*Import{}
	}

	plugin.Imports[name][proc] = imp
}

func (plugin *Plugin) SetPluginGlobal(pluginRet *lua.LTable) {
	state := plugin.State
	plugin.Environment.Plugin = pluginRet

	state.SetField(state.G.Global, "Plugin", pluginRet)
}

func (plugin *Plugin) Init() {
	state := plugin.State

	pluginReturn := state.ToTable(-1)
	state.Pop(1)

	plugin.SetPluginGlobal(pluginReturn)

	metadata := GetTableOrPanic(
		state,
		pluginReturn,
		"Metadata",
		"Couldn't find metadata table",
	)

	events := GetTableOrPanic(
		state,
		pluginReturn,
		"Events",
		"Couldn't find events table",
	)

	commands := GetTableOrPanic(
		state,
		pluginReturn,
		"Commands",
		"Couldn't find commands table",
	)

	imports := GetTableOrPanic(
		state,
		pluginReturn,
		"Imports",
		"Couldn't find imports table",
	)

	name := GetStringOrPanic(
		state,
		metadata,
		"Name",
		"Couldn't find plugin name",
	)

	description := GetStringOr(
		state,
		metadata,
		"Description",
		"N/A",
	)

	author := GetStringOr(
		state,
		metadata,
		"Author",
		"N/A",
	)

	version := GetStringOr(
		state,
		metadata,
		"Version",
		"N/A",
	)

	plugin.Info.Name = name
	plugin.Info.Description = description
	plugin.Info.Author = author
	plugin.Info.Version = version

	state.ForEach(events, func(key, value lua.LValue) {
		eventName := lua.LVAsString(key)
		eventFn := value.(*lua.LFunction)

		event := MakeEvent(eventName, eventFn)
		plugin.AddEvent(eventName, event)
	})

	state.ForEach(commands, func(key, value lua.LValue) {
		cmdName := lua.LVAsString(key)
		cmdEntry := value.(*lua.LTable)

		cmd := MakeCommand(state, cmdName, cmdEntry)
		plugin.AddCommand(cmdName, cmd)
	})

	state.ForEach(imports, func(key, value lua.LValue) {
		importEntry := value.(*lua.LTable)
		MakeCallableTable(
			state,
			importEntry,
			plugin.CallImport,
		)
	})
}

func (plugin *Plugin) CallImport(state *lua.LState) int {
	importTbl := state.ToTable(1)

	pluginName := GetStringOrPanic(
		state,
		importTbl,
		"Plugin",
		"Couldn't find plugin name",
	)

	procName := GetStringOrPanic(
		state,
		importTbl,
		"Procedure",
		"Couldn't find procedure name",
	)

	var args []lua.LValue

	top := state.GetTop()

	for i := 2; i <= top; i++ {
		args = append(args, state.Get(i))
	}

	for _, imp := range plugin.Imports[pluginName] {
		if imp.Command.Name == procName {
			ret, _ := imp.Command.Invoke(imp.Plugin.State, imp.Plugin.Environment.Plugin, args...)
			plugin.State.Push(ret)
		}
	}

	return 1
}

func (plugin *Plugin) Load() error {
	err := plugin.State.DoFile(plugin.Info.Path)

	if err != nil {
		return err
	}

	plugin.Init()
	plugin.FireEvent("OnLoad")

	return nil
}

func (plugin *Plugin) Unload() {
	plugin.FireEvent("OnUnload")
	plugin.State.Close()
}
