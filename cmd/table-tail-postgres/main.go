package main

import (
	"github.com/gitu/table-tail/cmd/table-tail/commands"
	_ "github.com/gitu/table-tail/pkg/utils/postgres"
	_ "github.com/lib/pq"
)

func main() {
	commands.SetDefaultDriver("postgres")
	commands.Execute()
}
