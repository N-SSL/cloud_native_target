package MySQL

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "golang.org/x/crypto/bcrypt"
	"log"
)

var SqlDB *gorm.DB

func ConnectToServer() error {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		log.Println(fmt.Errorf("Fatal error config file: %w \n", err))
		return err
	}
	mysqlUser := viper.GetString("user")
	mysqlPassword := viper.GetString("password")
	mysqlHost := viper.GetString("host")
	mysqlPort := viper.GetInt("port")
	mysqlSchema := viper.GetString("schema")
	connecting := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",mysqlUser,mysqlPassword,mysqlHost,mysqlPort,mysqlSchema)

	db, err := gorm.Open(mysql.Open(connecting),&gorm.Config{})

	if err != nil{
		log.Println("cannot connect database")
		return err
	}

	SqlDB = db
	_ = SqlDB.AutoMigrate(&Users{},&EnvStruct{},&RunningStruct{})

	// bytepwd,_ := bcrypt.GenerateFromPassword([]byte("123456"),0)
	// var admin = Users{
	// 	Username: "admin",
	// 	Password: bytepwd,
	// 	Role: "admin",
	// }

	// _ = SqlDB.Create(&admin)

	return nil
}
