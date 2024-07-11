package mods

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/slh335/mc-modpack-manager/services"
	"github.com/spf13/cobra"
)

const mcVersion string = "1.21"

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [id or slug]",
	Short: "Add a mod",
	Long:  `Add a mod`,
	Args: func(cmd *cobra.Command, args []string) error {
		err := cobra.MinimumNArgs(1)(cmd, args)
		if cmd.Flag("file").Value.String() == "" && err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		modIds := args

		file := cmd.Flag("file").Value.String()
		if file != "" {
			content, err := os.ReadFile(file)
			if err == nil {
				lines := strings.Split(string(content), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line != "" {
						modIds = append(modIds, line)
					}
				}
			}
		}

		mods, err := services.GetModrinthModInformation(modIds, mcVersion)
		if err != nil {
			log.Fatal(err)
			return
		}

		for _, mod := range mods {
			fmt.Println("Downloading", mod.Title, "version", mod.LatestVersion.Filename)
			err = mod.Install()
			if err != nil {
				log.Fatal(err)
				return
			}
		}
	},
}

func init() {
	modsCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	addCmd.Flags().StringP("file", "f", "", "Add and install mods from file")
}
