package mods

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/slh335/mc-modpack-manager/services"
	"github.com/spf13/cobra"
)

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
		loader, mcVersion, err := services.GetVersion()
		if err != nil {
			log.Fatal("Error: Version information could not be loaded from index")
			return
		}

		modIds := []string{}
		for _, arg := range args {
			if arg != "" {
				modIds = append(modIds, arg)
			}
		}

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

		mods, err := services.GetModrinthModInformation(modIds, loader, mcVersion, true)
		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Println("Loading dependencies")
		for _, mod := range mods {
			newDependencies := []string{}
			for _, dependencyId := range mod.LatestVersion.Dependencies {
				for _, curr := range mods {
					if curr.Id == dependencyId {
						newDependencies = append(newDependencies, dependencyId)
					}
				}

			}
			modDependencies, err := services.GetModrinthModInformation(newDependencies, loader, mcVersion, false)
			if err == nil {
				mods = append(mods, modDependencies...)
			}
		}

		slices.SortFunc(mods, func(a, b services.ModrinthMod) int {
			if a.Id < b.Id {
				return -1
			} else if b.Id < a.Id {
				return 1
			} else {
				return 0
			}
		})
		mods = slices.CompactFunc(mods, func(a, b services.ModrinthMod) bool {
			return a.Id == b.Id
		})

		if len(mods) == 0 {
			fmt.Println("No valid mod ids or slugs were supplied")
			return
		}

		for i, mod := range mods {
			fmt.Printf("\033[2K\rDownloading mod [%d/%d] %s", i+1, len(mods), mod.Title)
			err = mod.Install()
			if err != nil {
				log.Fatal(err)
				return
			}
		}
		fmt.Print("\n")
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
