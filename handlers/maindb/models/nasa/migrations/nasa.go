package nasa

import (
	"simpleBackend/handlers/maindb/models/nasa"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Migrations returns gormigrate.Migration instances
func Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "add_nasa_apod_1617960618",
			Migrate: func(tx *gorm.DB) error {
				tx.AutoMigrate()
				return tx.AutoMigrate(&nasa.Apod{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&nasa.Apod{})
			},
		},
	}
}
