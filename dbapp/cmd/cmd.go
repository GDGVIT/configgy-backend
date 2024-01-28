package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var DbAppsCmd = &cobra.Command{
	Use:   "dbapp",
	Short: "run different db operations",
	Run: func(cmd *cobra.Command, args []string) {
		// Current environment
		env := os.Getenv("ENVIRONMENT")
		fmt.Println("Current Environment: ", env)
	},
}

func init() {
	DbAppsCmd.AddCommand(DropTables())
	DbAppsCmd.AddCommand(Migrate())
	DbAppsCmd.AddCommand(Seed())
	DbAppsCmd.AddCommand(Backup())
}
