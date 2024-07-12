package server

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/slh335/mc-modpack-manager/services"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [software] [version] [name]",
	Short: "Create a new Minecraft server and initialize the mod manager",
	Long:  `Create a new Minecraft server and initialize the mod manager`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(3)(cmd, args); err != nil {
			return err
		}
		if strings.ToLower(args[0]) != "fabric" {
			return fmt.Errorf("invalid server software specified: %s", args[0])
		}
		pattern := regexp.MustCompile("^1(\\.[0-9]{1,2}){1,2}$")
		if !pattern.MatchString(args[1]) {
			return fmt.Errorf("invalid Minecraft version specified: %s", args[1])
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		software := strings.ToLower(args[0])
		mcVersion := args[1]
		name := strings.Join(args[2:], " ")
		fmt.Printf("Creating %s server on %s named %s...\n", software, mcVersion, name)

		dir := filepath.Base(strings.ReplaceAll(strings.ToLower(name), " ", "-"))
		os.Mkdir(dir, os.ModePerm)

		fileName, err := services.DownloadLatestFabricServer(mcVersion, dir)
		if err != nil {
			log.Fatal("Error: Failed to download server")
			return
		}
		fmt.Println("Downloaded server " + fileName)

		command := exec.Command("sh", "start.sh")
		command.Dir = dir
		err = command.Run()
		if err != nil {
			log.Fatal("Error: Failed to install server", err)
			return
		}
		fmt.Println("Installed server")

		os.WriteFile(filepath.Join(dir, "eula.txt"), []byte("eula=true"), os.ModePerm)
	},
}

func init() {
	serverCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
