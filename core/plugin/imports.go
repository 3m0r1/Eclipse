package plugin

type Import struct {
	Plugin  *Plugin
	Command *Command
}

func MakeImport(plugin *Plugin, cmd *Command) *Import {
	imp := Import{
		Plugin:  plugin,
		Command: cmd,
	}

	return &imp
}
