package cmd

import (
	"github.com/GabrielDCelery/collab-todo-tui/internals/model"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "short",
	Long:  `long`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		model.Run()
	},
}
