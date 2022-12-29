package datasource

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

//func InstanceDb(conf map[string]interface{}, logMode bool) *gorm.DB {
//	//conf := model.GetConfig("database").(map[string]interface{})[handle].(map[string]interface{})
//	driveSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
//		conf["user"], conf["password"], conf["host"], conf["port"], conf["name"])
//	db, err := gorm.Open("mysql", driveSource)
//	if err != nil {
//		util.Error("dbhelper.InstanceDb:", err.Error())
//		//log.Println("dbhelper.InstanceDb:", err)
//		return nil
//	}
//	// 开启打印sql
//	db.LogMode(logMode)
//	//db.SetLogger(log.New(os.Stdout, "\r\n", 0))
//	db.DB().SetMaxOpenConns(256)
//	db.DB().SetMaxIdleConns(5)
//	db.DB().SetConnMaxLifetime(8 * time.Second) //hulk mysql默认超时时间是10秒
//	// 全局禁用表名复数
//	db.SingularTable(true)
//	return db
//}

func InstanceDb(conf map[string]interface{}, debug bool) *gorm.DB {
	driveSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%ds&readTimeout=%ds&charset=%s&parseTime=True&loc=Local",
		conf["user"],
		conf["password"],
		conf["host"],
		conf["port"],
		conf["db_name"],
		conf["timeout"],
		conf["read_timeout"],
		conf["charset"])
	db, err := gorm.Open(mysql.Open(driveSource),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
	if err != nil {
		//log.Println("dbhelper.InstanceDb:", err)
		panic(err)
		return nil
	}
	if debug {
		db = db.Debug()
	}
	rawDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	//db.SetLogger(log.New(os.Stdout, "\r\n", 0))
	rawDB.SetMaxOpenConns(int(conf["max_idle"].(int64)))
	rawDB.SetMaxIdleConns(int(conf["max_connection"].(int64)))
	rawDB.SetConnMaxLifetime(8 * time.Second) //hulk mysql默认超时时间是10秒
	return db
}
