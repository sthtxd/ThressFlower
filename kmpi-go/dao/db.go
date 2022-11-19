package dao

import (
	"fmt"
	"kmpi-go/config"
	"kmpi-go/log"
	"time"

	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	userName  string = ""
	password  string = ""
	ipAddrees string = ""
	port      int    = -1
	dbName    string = "furnace_monitor"
	charset   string = "utf8"
)
var CronDb *gorm.DB

func init() {
	initConfig()
	CronDb = InitDb()
}

// func initSqlServerDb() (*gorm.DB, error) {
// 	dsn := fmt.Sprintf("server=%s;port%d;database=%s;user id=%s;password=%s", ipAddrees, port, dbName, userName, password)
// 	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

// 	err = db.AutoMigrate(&FurnaceTestDao{})
// 	if err != nil {
// 		log.Error("AutoMigrate error", err)
// 	}
// 	return db, err
// }
func initMysqlDb() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&timeout=1s&readTimeout=360s&parseTime=true&loc=Local", userName, password, ipAddrees, port, dbName, charset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Error("mysql connect failed, detail is [%v]", err.Error())
		panic(err)
	}
	realDb, err := db.DB()
	if err != nil {
		log.Error("mysql db() error,", err.Error())
		panic(err)
	}
	realDb.SetConnMaxLifetime(3600 * time.Second)
	realDb.SetMaxIdleConns(20)
	realDb.SetMaxOpenConns(50)
	err = db.AutoMigrate(&AllDeviceDataDao{})
	if err != nil {
		log.Error("AllDeviceDataDao error", err)
	}
	err = db.AutoMigrate(&WeightRfidDao{})
	if err != nil {
		log.Error("WeightRfidDao error", err)
	}
	err = db.AutoMigrate(&WeightRfidLogDao{})
	if err != nil {
		log.Error("WeightRfidLogDao error", err)
	}
	err = db.AutoMigrate(&DeviceLogDao{})
	if err != nil {
		log.Error("DeviceLogDao error", err)
	}
	err = db.AutoMigrate(&UserDao{})
	if err != nil {
		log.Error("UserDao error", err)
	}
	err = db.AutoMigrate(&DownloadDao{})
	if err != nil {
		log.Error("DownloadDao error", err)
	}

	return db, err

}
func InitDb() *gorm.DB {
	//db, err := initMysqlDb()
	db, err := initMysqlDb()

	if err != nil {
		panic(err)
	}

	return db
}
func initConfig() {
	userNameRes, err := config.GetValue("mysql", "username")
	if err != nil {
		log.Error("GetValue error", err.Error())
	}
	userName = userNameRes
	passwordRes, err := config.GetValue("mysql", "password")
	if err != nil {
		log.Error("GetValue error", err.Error())
	}
	password = passwordRes
	ipAddreesRes, err := config.GetValue("mysql", "ipAddrees")
	if err != nil {
		log.Error("GetValue error", err.Error())
	}
	ipAddrees = ipAddreesRes

	portString, err := config.GetValue("mysql", "port")
	if err != nil {
		log.Error("GetValue error", err.Error())
	}

	portNew, err := strconv.Atoi(portString)
	if err != nil {
		log.Error("atoi error:%s", err.Error())
		port = -1
	} else {
		port = portNew
	}
}
