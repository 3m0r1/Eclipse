package plugin

import (
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

func (plugin *Plugin) AddEvent(name string, event *Event) {
	plugin.Events[name] = event
}

func (plugin *Plugin) FireEvent(name string, args ...lua.LValue) {
	event := plugin.Events[name]
	event.Fire(
		plugin.State,
		args...,
	)
}

func (plugin *Plugin) GetCommand(name string) *Command {
	for _, cmd := range plugin.Commands {
		if cmd.Name == name {
			return cmd
		}
	}

	return nil
}
func (plugin *Plugin) AddCommand(name string, command *Command) {
	plugin.Commands[name] = command
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

	metadata := state.GetField(
		pluginReturn,
		"Metadata",
	)

	events := state.GetField(
		pluginReturn,
		"Events",
	)

	commands := state.GetField(
		pluginReturn,
		"Commands",
	)

	name := lua.LVAsString(state.GetField(
		metadata,
		"Name",
	))

	description := lua.LVAsString(state.GetField(
		metadata,
		"Description",
	))

	author := lua.LVAsString(state.GetField(
		metadata,
		"Author",
	))

	version := lua.LVAsString(state.GetField(
		metadata,
		"Version",
	))

	plugin.Info.Name = name
	plugin.Info.Description = description
	plugin.Info.Author = author
	plugin.Info.Version = version

	state.ForEach(events.(*lua.LTable), func(key, value lua.LValue) {
		name := lua.LVAsString(key)
		fn := value.(*lua.LFunction)

		event := MakeEvent(name, fn)
		plugin.AddEvent(name, event)
	})

	state.ForEach(commands.(*lua.LTable), func(key, value lua.LValue) {
		name := lua.LVAsString(key)
		entry := value.(*lua.LTable)

		cmd := MakeCommand(state, name, entry)
		plugin.AddCommand(name, cmd)
	})
}

func (plugin *Plugin) CallImport(state *lua.LState) int {
	importTbl := state.ToTable(1)

	pluginName := lua.LVAsString(state.GetField(importTbl, "Plugin"))
	procName := lua.LVAsString(state.GetField(importTbl, "Procedure"))

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
func (plugin *Plugin) SetImportsMetatable() {
	state := plugin.State

	imports := state.GetField(
		plugin.Environment.Plugin,
		"Imports",
	)

	state.ForEach(imports.(*lua.LTable), func(key, value lua.LValue) {
		entry := value.(*lua.LTable)

		mt := state.NewTable()
		state.SetField(mt, "__call", state.NewFunction(plugin.CallImport))
		state.SetMetatable(entry, mt)
	})
}
func (plugin *Plugin) Load() {
	plugin.State.DoFile(plugin.Info.Path)
	plugin.Init()
	plugin.SetImportsMetatable()
	plugin.FireEvent("OnLoad")
}

func (plugin *Plugin) Unload() {
	plugin.FireEvent("OnUnload")
	plugin.State.Close()
}
