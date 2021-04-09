package cmd

import (
	"simpleBackend/config/load"
	"simpleBackend/handlers/maindb"
	"simpleBackend/log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func databaseOperations() *cobra.Command {
	var configPath string

	database := &cobra.Command{
		Use:   "database",
		Short: "database operations",
		Run: func(c *cobra.Command, args []string) {
			c.Usage()
		},
	}

	mainDBSchemeMigration := &cobra.Command{
		Use:   "mainMigrate",
		Short: "run scheme migration of main database",
		Run: func(c *cobra.Command, args []string) {
			// config
			conf, err := load.LoadFromFile(configPath)
			if err != nil {
				log.Logger.Error("failed to load config from file", zap.Error(err))
				return
			}

			// because
			maindb.New(conf.DB.Main.Type, conf.DB.Main.DSN, 0, 0, 0)
		},
	}

	mainDBSchemeMigration.Flags().StringVarP(&configPath, "config", "-c", "config/example.yaml", "path to the database config file")

	database.AddCommand(
		mainDBSchemeMigration,
	)

	return database
}
