package cmd

import (
	"eclipse/core/manager"
	"eclipse/core/plugin"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load [plugins]",
	Short: "Loads the plugins in the specified directory",
	Args:  cobra.ExactArgs(1),
	Run:   LoadPlugins,
}

func LoadPlugins(cmd *cobra.Command, args []string) {
	dir := args[0]

	mgr := manager.NewPluginManager()
	var plugins []*plugin.Plugin

	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".lua") {
			plugins = append(plugins, plugin.NewPlugin(path))
		}

		return nil
	})

	if err := mgr.LoadPlugins(plugins...); err != nil {
		fmt.Println(err)
	}
}

func Execute() {
	if err := loadCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
