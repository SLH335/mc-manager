package mods

import (
	"fmt"

	"github.com/slh335/mc-manager/cmd"
	"github.com/spf13/cobra"
)

// modsCmd represents the mods command
var modsCmd = &cobra.Command{
	Use:   "mods",
	Short: "Manage mods",
	Long:  `Manage mods`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mods called")
	},
}

func init() {
	cmd.RootCmd.AddCommand(modsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// modsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// modsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
