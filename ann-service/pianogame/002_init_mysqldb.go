package pianogame

import (
	"io/ioutil"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	yaml "gopkg.in/yaml.v2"
)

// MysqlDB DB client exported
var MysqlDB *gorm.DB

func init() {

	log.Println("Init mysql db start....")
	/* Load config by yaml format */
	var config MysqlConfig
	var DBOpenSetting string
	configFile, configErr := ioutil.ReadFile("config/database/mysql/config.yaml") // open file and read
	if configErr != nil {
		log.Panicf("read mysql database config file error: %v", configErr)
	} // fi
	configUnmarshalError := yaml.Unmarshal([]byte(configFile), &config)
	if configUnmarshalError != nil {
		log.Panicf("error while unmarshal from db config: %v", configUnmarshalError)
	} // fi
	log.Println("Load mysql config file finished")
	//x = strConcate(config.User, config.Password, config.DBName)
	//log.Println(x)
	DBOpenSetting = strConcateF(
		"%s:%s@/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.DBName,
	)

	var err error
	// TO-DO: DB name from; connect schema : user:pwd@DB name?option1=x&option2=y...
	MysqlDB, err = gorm.Open("mysql", DBOpenSetting)
	if err != nil {
		panic("Mysql DB open failed!")
	}
	// Migration test
	/*
		WARNING: AutoMigrate will ONLY create tables, missing columns and missing indexes,
		         and WON'T change existing column's type or delete unused columns to protect your data.
	*/
	// reset tables in test!
	//MysqlDB.DropTableIfExists(&User{}, &Email{}, &Address{}, &CreditCard{}, &UserProfile{}) // comment out if in production!!
	MysqlDB = MysqlDB.Set("gorm:table_options", "ENGINE=InnoDB")

	MysqlDB.AutoMigrate(&User{})

	MysqlDB.AutoMigrate(&Email{})
	// MysqlDB.AutoMigrate(&Address{})
	// MysqlDB.AutoMigrate(&CreditCard{})
	MysqlDB.AutoMigrate(&UserProfile{})

	// Add table suffix when create tables
	//MysqlDB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})

	// defer MysqlDB.Close()
}
