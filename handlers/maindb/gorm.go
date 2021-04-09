// Package maindb implements a adapter/wrapper of gorm pkg with different sql driver
package maindb

import (
	"errors"
	"simpleBackend/handlers/maindb/postgres"
	"time"

	"gorm.io/gorm"
)

// New returns gorm object
func New(dbType, dsn string, maxOpenConns, maxIdleConns int, connMaxLife time.Duration) (*gorm.DB, error) {
	switch dbType {
	case "postgresql":
		return postgres.New(dsn, maxIdleConns, maxOpenConns, connMaxLife)
	default:
		return nil, errors.New("unknown database type: " + dbType)
	}
}
