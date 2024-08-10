package mods

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/slh335/mc-manager/services"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [mod loader] [game version]",
	Short: "Initialize mod index",
	Long:  `Initialize mod index`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(2)(cmd, args); err != nil {
			return err
		}
		if strings.ToLower(args[0]) != "fabric" {
			return fmt.Errorf("invalid mod loader specified: %s", args[0])
		}
		pattern := regexp.MustCompile("^1(\\.[0-9]{1,2}){1,2}$")
		if !pattern.MatchString(args[1]) {
			return fmt.Errorf("invalid Minecraft version specified: %s", args[1])
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		loader := args[0]
		mcVersion := args[1]

		err := services.InitIndex(loader, mcVersion)
		if err != nil {
			fmt.Println("Error: Failed to initialize mod index")
		} else {
			fmt.Println("Initialized mod index for", loader, mcVersion)
		}
	},
}

func init() {
	modsCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
