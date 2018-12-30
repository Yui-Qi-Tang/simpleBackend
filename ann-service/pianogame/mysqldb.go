package pianogame

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// MysqlDB DB client exported
var MysqlDB *gorm.DB

func init() {

	log.Println("Init mysql db start....")
	var err error
	// TO-DO: DB name from; connect schema : user:pwd@DB name?option1=x&option2=y...
	MysqlDB, err = gorm.Open("mysql", "root:@/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("Mysql DB open failed!")
	}
	// Migration test
	MysqlDB.AutoMigrate(&User{})

	MysqlDB.AutoMigrate(&Email{})
	MysqlDB.AutoMigrate(&Address{})
	MysqlDB.AutoMigrate(&CreditCard{})

	// Add table suffix when create tables
	// db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})

	if x := MysqlDB.HasTable("users555"); x == false {
		log.Println("No table")
	}
	defer MysqlDB.Close()
}
