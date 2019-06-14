package main

import (
	"github.com/gitu/table-tail/cmd/table-tail/commands"
	_ "github.com/gitu/table-tail/pkg/utils/oracle"
	_ "gopkg.in/goracle.v2"
)

func main() {
	commands.SetDefaultDriver("goracle")
	commands.Execute()
}
