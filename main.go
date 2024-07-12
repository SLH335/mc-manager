package main

import (
	"github.com/slh335/mc-modpack-manager/cmd"
	_ "github.com/slh335/mc-modpack-manager/cmd/mods"
	_ "github.com/slh335/mc-modpack-manager/cmd/server"
)

func main() {
	cmd.Execute()
}
