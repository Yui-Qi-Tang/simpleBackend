package cmd

import (
	"simpleBackend/config/load"
	"simpleBackend/handlers/maindb"
	nasaMigrations "simpleBackend/handlers/maindb/models/nasa/migrations"
	"simpleBackend/log"

	"github.com/go-gormigrate/gormigrate/v2"

	"github.com/pkg/errors"
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

	mainDatabase := &cobra.Command{
		Use:   "main",
		Short: "main database",
		Run: func(c *cobra.Command, args []string) {
			c.Usage()
		},
	}

	mainDBSchemeMigration := &cobra.Command{
		Use:   "migrate",
		Args:  cobra.MaximumNArgs(3), // args[0]: model name, args[1]: op{up/down/list}, args[2]: migration id
		Short: "run scheme migration of main database",
		Run: func(c *cobra.Command, args []string) {
			log.Logger.Info("start database scheme migrating...")
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

			// migrations
			switch args[0] {
			case "nasa":
				migrate := gormigrate.New(mdb, gormigrate.DefaultOptions, nasaMigrations.Migrations())
				if len(args) >= 3 {
					if err := mainDBSchemeMigration(migrate, args[1], args[2]); err != nil {
						log.Logger.Error("failed to migrate", zap.Error(err))
						return
					}
				} else {
					if err := mainDBSchemeMigration(migrate, args[1], ""); err != nil {
						log.Logger.Error("failed to migrate", zap.Error(err))
						return
					}

				}
			default:
				log.Logger.Error("unknown model")
			}

			log.Logger.Info("Database scheme migration is done...")

		},
	}

	mainDBSchemeMigration.Flags().StringVarP(&configPath, "config", "c", "config/example.yaml", "path to the database config file")

	mainDatabase.AddCommand(
		mainDBSchemeMigration,
	)

	database.AddCommand(
		mainDatabase,
	)

	return database
}

func mainDBSchemeMigration(m *gormigrate.Gormigrate, op string, mID string) error {

	switch op {
	case "up":
		if len(mID) > 0 {
			if err := m.MigrateTo(mID); err != nil {
				return errors.Wrap(err, "failed to migrate up")
			}
		} else {
			if err := m.Migrate(); err != nil {
				return errors.Wrap(err, "failed to migrate up")
			}
		}
	case "down":
		if len(mID) > 0 {
			if err := m.RollbackTo(mID); err != nil {
				return errors.Wrap(err, "failed to migrate down")
			}
		} else {
			if err := m.RollbackLast(); err != nil {
				return errors.Wrap(err, "failed to migrate down")
			}
		}
	default:
		return errors.New("unknown opera")
	}
	return nil
}
