package main

import (
	"eclipse/core/manager"
	"eclipse/core/plugin"

	lua "github.com/yuin/gopher-lua"
)

func main() {
	mgr := manager.NewPluginManager()

	plugin1 := plugin.NewPlugin("./plugins/basic.lua")
	plugin2 := plugin.NewPlugin("./plugins/terminal.lua")

	mgr.LoadPlugins(plugin1, plugin2)

	plugin1.InvokeCommand("Greet", lua.LString("World"))
}
