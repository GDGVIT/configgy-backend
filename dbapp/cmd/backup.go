package cmd

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/GDGVIT/configgy-backend/dbapp/pkg"
	"github.com/spf13/cobra"
)

func Backup() *cobra.Command {
	return &cobra.Command{
		Use: "backup",
		RunE: func(cmd *cobra.Command, args []string) error {

			dbConnection, sqlConnection := pkg.Connection()
			defer sqlConnection.Close()

			var tableNames []string
			if err := dbConnection.Table("information_schema.tables").
				Where("table_schema = ?", "public").Pluck("table_name", &tableNames).Error; err != nil {
				panic(err)
			}

			directory := "backup"
			// Check if directory exists
			if _, err := os.Stat(directory); os.IsNotExist(err) {
				// Create directory if it doesn't exist
				if err := os.Mkdir(directory, 0755); err != nil {
					panic(err)
				}
			}
			// Save all data in respective csv files
			for _, tableName := range tableNames {
				file, err := os.Create(directory + "/" + tableName + ".csv")
				if err != nil {
					panic(err)
				}
				defer file.Close()

				rows, err := dbConnection.Table(tableName).Rows()
				if err != nil {
					panic(err)
				}

				defer rows.Close()

				csvWriter := csv.NewWriter(file)
				defer csvWriter.Flush()

				columns, err := rows.Columns()
				if err != nil {
					panic(err)
				}

				csvWriter.Write(columns)

				for rows.Next() {
					columns, err := rows.Columns()
					if err != nil {
						panic(err)
					}

					values := make([]interface{}, len(columns))
					valuePtrs := make([]interface{}, len(columns))

					for i := range columns {
						valuePtrs[i] = &values[i]
					}

					rows.Scan(valuePtrs...)
					// Properly format the values
					for i := range columns {
						if values[i] == nil {
							values[i] = "\\N"
						}
					}

					var row []string
					for i := range columns {
						// add to csv
						data := fmt.Sprintf("%v", values[i])
						row = append(row, data)
					}
					csvWriter.Write(row)
				}
			}

			return nil
		},
	}
}
