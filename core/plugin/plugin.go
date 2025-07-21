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
	Exports     map[string]*Export
	Environment *Environment
	State       *lua.LState
}

func NewPlugin(path string) *Plugin {
	plugin := Plugin{
		Info:        &PluginInfo{Path: path},
		Events:      map[string]*Event{},
		Commands:    map[string]*Command{},
		Imports:     map[string]map[string]*Import{},
		Exports:     map[string]*Export{},
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
	if _, exists := plugin.GetEvent(name); exists {
		return errors.New("event already exists")
	}
	plugin.Events[name] = event
	return nil
}

func (plugin *Plugin) FireEvent(name string, args ...lua.LValue) error {
	event, ok := plugin.Events[name]
	if !ok {
		return errors.New("event doesn't exist")
	}

	return event.Fire(plugin.State, args...)
}

func (plugin *Plugin) GetCommand(name string) (*Command, bool) {
	command, exists := plugin.Commands[name]
	return command, exists
}

func (plugin *Plugin) GetExport(name string) (*Export, bool) {
	export, exists := plugin.Exports[name]
	return export, exists
}

func (plugin *Plugin) AddCommand(name string, command *Command) error {
	if _, exists := plugin.GetCommand(name); exists {
		return errors.New("command already exists")
	}
	plugin.Commands[name] = command

	if command.Export {
		if _, exists := plugin.GetExport(name); exists {
			return errors.New("export already exists")
		}
		plugin.Exports[name] = (*Export)(command)
	}

	return nil
}

func (plugin *Plugin) InvokeCommandCtx(name string, ctx lua.LValue, args ...lua.LValue) (*lua.LTable, error) {
	cmd, exists := plugin.GetCommand(name)
	if !exists {
		return nil, errors.New("command doesn't exist")
	}

	return cmd.Invoke(plugin.State, ctx, args...)
}

func (plugin *Plugin) InvokeCommand(name string, args ...lua.LValue) (*lua.LTable, error) {
	ctx := plugin.Environment.Plugin
	return plugin.InvokeCommandCtx(name, ctx, args...)
}

func (plugin *Plugin) GetImportTable(name string) (map[string]*Import, bool) {
	table, exists := plugin.Imports[name]
	return table, exists
}

func (plugin *Plugin) InitImportTable(name string) error {
	if _, exists := plugin.GetImportTable(name); exists {
		return errors.New("import table already exists")
	}
	plugin.Imports[name] = make(map[string]*Import)
	return nil
}

func (plugin *Plugin) GetImport(name, proc string) (*Import, bool) {
	table, exists := plugin.GetImportTable(name)
	if !exists {
		return nil, false
	}
	imp, exists := table[proc]
	if !exists {
		return nil, false
	}
	return imp, true
}

func (plugin *Plugin) AddImport(name, proc string, imp *Import) error {
	table, exists := plugin.GetImportTable(name)
	if !exists {
		return errors.New("import table doesn't exist")
	}
	table[proc] = imp
	return nil
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
		// TODO: handle error from AddEvent
		plugin.AddEvent(eventName, event)
	})

	state.ForEach(commands, func(key, value lua.LValue) {
		cmdName := lua.LVAsString(key)
		cmdEntry := value.(*lua.LTable)

		cmd := MakeCommand(state, cmdName, cmdEntry)
		// TODO: handle error from AddCommand
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

	if imp, exists := plugin.GetImport(pluginName, procName); exists {
		ret, _ := imp.Command.Invoke(
			imp.Plugin.State,
			imp.Plugin.Environment.Plugin,
			args...,
		)
		plugin.State.Push(ret)
		return 1
	} else {
		return 0
	}
}

func (plugin *Plugin) Load() error {
	if err := plugin.State.DoFile(plugin.Info.Path); err != nil {
		return err
	}

	plugin.Init()

	return plugin.FireEvent("OnLoad")
}

func (plugin *Plugin) Unload() error {
	err := plugin.FireEvent("OnUnload")
	plugin.State.Close()
	return err
}
