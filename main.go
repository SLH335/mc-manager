package main

import (
	"github.com/slh335/mc-manager/cmd"
	_ "github.com/slh335/mc-manager/cmd/mods"
	_ "github.com/slh335/mc-manager/cmd/server"
)

func main() {
	cmd.Execute()
}
