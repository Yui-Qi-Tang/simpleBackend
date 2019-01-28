package pianogame

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Email comment
type Email struct {
	ID         int
	UserID     int    `gorm:"index"`                          // Foreign key (belongs to), tag `index` will create index for this column
	Email      string `gorm:"type:varchar(100);unique_index"` // `type` set sql type, `unique_index` will create unique index for this column
	Subscribed bool
}

// Address comment
type Address struct {
	ID       int
	Address1 string         `gorm:"not null;unique"` // Set field as not nullable and unique
	Address2 string         `gorm:"not null;unique"` // Set field as not nullable and unique
	Address3 string         `gorm:"not null;unique"` // Set field as not nullable and unique
	Post     sql.NullString `gorm:"not null"`
}

// Language comment
type Language struct {
	ID   int
	Name string `gorm:"index:idx_name_code"` // Create index with name, and will create combined index if find other fields defined same name
	Code string `gorm:"index:idx_name_code"` // `unique_index` also works
}

// CreditCard comment
type CreditCard struct {
	gorm.Model
	UserID uint
	Number string
}

// User comment
type User struct {
	gorm.Model
	Account  sql.NullString `gorm:"type:varchar(256);unique;not null"`
	Password sql.NullString `gorm:"type:varchar(256);not null"`
	Profile  UserProfile
}

// UserProfile detail info of user
type UserProfile struct {
	ID       int
	Birthday time.Time
	Name     string  `gorm:"size:255"` // Default size for string is 255, reset it with this tag
	Emails   []Email // One-To-Many relationship (has many - use Email's UserID as foreign key)
}
