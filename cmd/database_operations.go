package cmd

import (
	"simpleBackend/config/load"
	"simpleBackend/handlers/maindb"
	"simpleBackend/handlers/maindb/models/nasa/migrations"
	"simpleBackend/log"

	"github.com/go-gormigrate/gormigrate/v2"

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

	// TODO: main migrate {model name} {up/down}
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

			// use the default settings about connection pool
			mdb, err := maindb.New(conf.DB.Main.Type, conf.DB.Main.DSN, 0, 0, 0)
			if err != nil {
				log.Logger.Error("failed to New gorm DB", zap.Error(err))
				return
			}

			sqlDB, err := mdb.DB()
			if err != nil {
				log.Logger.Error("failed to load db from database/sql", zap.Error(err))
				return
			}
			defer sqlDB.Close()

			migrate := gormigrate.New(mdb, gormigrate.DefaultOptions, migrations.Migrations())

			if err := migrate.Migrate(); err != nil {
				log.Logger.Error("failed to make migration", zap.Error(err))
			}

		},
	}

	mainDBSchemeMigration.Flags().StringVarP(&configPath, "config", "c", "config/example.yaml", "path to the database config file")

	database.AddCommand(
		mainDBSchemeMigration,
	)

	return database
}
