package cmd

import (
	"fmt"
	"os"

	"github.com/GDGVIT/configgy-backend/dbapp/pkg"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func Seed() *cobra.Command {
	var err error
	return &cobra.Command{
		Use: "seed",
		RunE: func(cmd *cobra.Command, args []string) error {
			// if config.App.Env != "development" {
			// 	fmt.Println("Warning: Environment is not development. Tables wont be seeded")
			// 	return nil
			// }
			godotenv.Load()
			dbConnection, sqlConnection := pkg.Connection()
			defer sqlConnection.Close()
			begin := dbConnection.Begin()

			if os.Getenv("ENVIRONMENT") != "development" {
				// for i, seed := range pkg.ProdSeeder(begin) {
				// 	if err = seed.Run(begin); err != nil {
				// 		begin.Rollback()
				// 		fmt.Println("[Seeder] Running seed failed")
				// 		panic(err)
				// 	}
				// 	fmt.Println("[", i, "]: ", "Seed table: ", seed.TableName)
				// }
			} else {
				fmt.Println("App env is development")
				for i, seed := range pkg.Seeder(begin) {
					if err = seed.Run(begin); err != nil {
						begin.Rollback()
						fmt.Println("[Seeder] Running seed failed")
						panic(err)
					}
					fmt.Println("[", i, "]: ", "Seed table: ", seed.TableName)
				}
			}
			begin.Commit()
			fmt.Println("Seeding Completed with Environment: ", os.Getenv("ENVIRONMENT"))
			return nil
		},
	}
}
