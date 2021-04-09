// Package postgres is a wrapper of gorm works with postgresql.
// Avoid to meet side-effect, so create this pkg
package postgres

import (
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// New returns gorm object with connection pool
// GORM using database/sql to maintain connection pool (from: https://gorm.io/docs/connecting_to_the_database.html#Connection-Pool)
func New(dsn string, maxIdleConns, maxOpenConns int, connMaxLifeTime time.Duration) (*gorm.DB, error) {

	// conn pool by database/sql
	if maxIdleConns <= 0 || maxOpenConns <= 0 || connMaxLifeTime == 0 {
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to new main db with conn pool")
	}

	sqldb, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "failed to set up connection pool by database/sql")
	}

	sqldb.SetMaxIdleConns(maxIdleConns)
	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetConnMaxIdleTime(connMaxLifeTime * time.Second)

	return db, nil

}
