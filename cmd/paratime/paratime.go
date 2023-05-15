package paratime

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "paratime",
	Short:   "ParaTime layer operations",
	Aliases: []string{"p", "pt"},
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(addCmd)
	Cmd.AddCommand(registerCmd)
	Cmd.AddCommand(removeCmd)
	Cmd.AddCommand(setDefaultCmd)
	Cmd.AddCommand(showCmd)
	Cmd.AddCommand(statsCmd)
}
