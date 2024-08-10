package mods

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/slh335/mc-manager/services"
	"github.com/slh335/mc-manager/ui"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for mods",
	Long:  `Search for mods`,
	Run: func(cmd *cobra.Command, args []string) {
		loader, mcVersion, err := services.GetVersion()
		if err != nil {
			log.Fatal("Error: Version information could not be loaded")
			return
		}

		query := strings.Join(args, "+")

		p := tea.NewProgram(ui.InitialSearchModel(query, loader, mcVersion), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
			return
		}
	},
}

func init() {
	modsCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
