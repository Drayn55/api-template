package config

import (
	"log"
	"modify-framework-mod/utils"
	"os"
)

func HandleCLI() {
	if len(os.Args) > 1 {
		args := os.Args[1:]
		cmd := args[0]

		switch cmd {
		case "migrate", "migrate:refresh", "make:migration":
			utils.HandleMigrationCommand(args)

		case "db:seed", "create:seeder":
			utils.HandleSeederCommand(args)

		case "fix:dirty":
			utils.FixDirtyMigration()

		default:
			log.Println("⚠️ Command tidak dikenali:", cmd)
			log.Println("Usage:")
			log.Println("  go run main.go migrate")
			log.Println("  go run main.go migrate:refresh")
			log.Println("  go run main.go make:migration <name>")
			log.Println("  go run main.go db:seed")
			log.Println("  go run main.go create:seeder <name>")
			log.Println("  go run main.go fix:dirty")
		}
		os.Exit(0)
	}
}
