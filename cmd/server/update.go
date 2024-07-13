package server

import (
	"fmt"
	"log"

	"github.com/slh335/mc-modpack-manager/services"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update [software] [version] [name]",
	Short: "update a new Minecraft server and initialize the mod manager",
	Long:  `update a new Minecraft server and initialize the mod manager`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		loader, mcVersion, err := services.DetectVersion()
		fmt.Println("Detected", loader, mcVersion)
		if err != nil {
			log.Fatal(err)
			return
		}
		if loader == "fabric" {
			err = services.UpdateFabricServer(mcVersion)
			if err != nil {
				log.Fatal("Error: Failed to update server")
				return
			}
		}
	},
}

func init() {
	serverCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
