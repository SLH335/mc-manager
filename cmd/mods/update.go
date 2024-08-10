package mods

import (
	"log"

	"github.com/slh335/mc-manager/services"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update [ids or slugs]",
	Short: "Update mods",
	Long:  `Update mods`,
	Run: func(cmd *cobra.Command, args []string) {
		loader, mcVersion, err := services.GetVersion()
		if err != nil {
			log.Fatal("Error: Version information could not be loaded")
			return
		}

		modIds := []string{}
		for _, arg := range args {
			if arg != "" {
				modIds = append(modIds, arg)
			}
		}

		services.UpdateMods(modIds, loader, mcVersion)
	},
}

func init() {
	modsCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
