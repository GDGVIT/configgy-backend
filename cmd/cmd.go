package cmd

import (
	"fmt"
	"os"

	api "github.com/GDGVIT/configgy-backend/api/cmd"
	dbapp "github.com/GDGVIT/configgy-backend/dbapp/cmd"
	mailer "github.com/GDGVIT/configgy-backend/mailer/cmd"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello World")
	},
}

// Execute - starts the CLI
func init() {
	cmd.AddCommand(mailer.RootCmd)
	cmd.AddCommand(api.RootCmd())
	cmd.AddCommand(dbapp.DbAppsCmd)
}

func Execute() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
