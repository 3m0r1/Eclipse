package manager

import (
	"eclipse/core/plugin"
	"errors"
	"slices"

	lua "github.com/yuin/gopher-lua"
)

type PluginManager struct {
	Plugins []*plugin.Plugin
}

func NewPluginManager() *PluginManager {
	mgr := PluginManager{
		Plugins: []*plugin.Plugin{},
	}

	return &mgr
}

func (mgr *PluginManager) GetPlugin(name string) (*plugin.Plugin, bool) {
	for _, plugin := range mgr.Plugins {
		info := plugin.Info
		if info.Name == name {
			return plugin, true
		}
	}

	return nil, false
}

func (mgr *PluginManager) LoadPlugin(plugin *plugin.Plugin) error {
	if err := plugin.Load(); err != nil {
		return err
	}

	_, exists := mgr.GetPlugin(plugin.Info.Name)

	if exists {
		return errors.New("plugin already exists")
	}

	if plugin.Info.Name == "" {
		return errors.New("plugin doesn't have a name")
	}

	mgr.Plugins = append(mgr.Plugins, plugin)
	return nil
}

func (mgr *PluginManager) RemovePlugin(name string) error {
	for index, plugin := range mgr.Plugins {
		info := plugin.Info

		if info.Name == name {
			mgr.Plugins[index].Unload()
			mgr.Plugins = slices.Delete(mgr.Plugins, index, index+1)
			return nil
		}
	}

	return errors.New("plugin doesn't exist")
}

func (mgr *PluginManager) InitImports(targetPlugin *plugin.Plugin, imports *lua.LTable) {
	state := targetPlugin.State

	state.ForEach(imports, func(key, value lua.LValue) {
		entry := value.(*lua.LTable)

		pluginName := plugin.GetStringOrPanic(
			state,
			entry,
			"Plugin",
			"Couldn't find plugin name",
		)

		procName := plugin.GetStringOrPanic(
			state,
			entry,
			"Procedure",
			"Couldn't find proceedure name",
		)

		if foundPlugin, ok := mgr.GetPlugin(pluginName); ok {
			targetPlugin.InitImportTable(pluginName)

			if export, ok := foundPlugin.GetExport(procName); ok {
				imp := plugin.MakeImport(foundPlugin, (*plugin.Command)(export))
				targetPlugin.AddImport(pluginName, procName, imp)
			}
		}

	})
}

func (mgr *PluginManager) SetPluginsImports() {
	for _, plugin := range mgr.Plugins {
		mgr.InitImports(
			plugin,
			plugin.State.GetField(
				plugin.Environment.Plugin, "Imports",
			).(*lua.LTable),
		)
	}
}

func (mgr *PluginManager) LoadPlugins(plugins ...*plugin.Plugin) error {
	for _, plugin := range plugins {
		if err := mgr.LoadPlugin(plugin); err != nil {
			return err
		}
	}
	mgr.SetPluginsImports()

	for _, plugin := range plugins {
		plugin.FireEvent("OnReady", plugin.Environment.Plugin)
	}

	return nil
}
