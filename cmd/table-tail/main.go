package main

import (
	"github.com/gitu/table-tail/cmd/table-tail/commands"
	_ "github.com/gitu/table-tail/pkg/utils/oracle"
	_ "github.com/gitu/table-tail/pkg/utils/postgres"
	_ "github.com/godror/godror"
)

func main() {
	commands.Execute()
}
