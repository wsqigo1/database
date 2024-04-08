package gorm

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

//func init() {
//	// 定义 dsn
//	const dsn = "root:secret@tcp(127.0.0.1:13306)/gormExample?charset=utf8mb4&parseTime=True&loc=Local"
//
//	// 连接服务器（池）
//	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
//		Logger: logger.Default.LogMode(logger.Info),
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	DB = db
//}

var logWriter io.Writer

func init() {
	// 定义DSN
	const dsn = "root:secret@tcp(127.0.0.1:13306)/gormExample?charset=utf8mb4&parseTime=True&loc=Local"
	// 初始化logWrite
	var err error
	logWriter, err = os.OpenFile("./sql.log", os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	customLogger := logger.New(log.New(logWriter, "\n", log.LstdFlags),
		logger.Config{
			// 慢查询阈值 200ms
			SlowThreshold: 200 * time.Millisecond,
			// 日志级别
			LogLevel: logger.Info,
			// 是否忽略记录不存在的错误
			IgnoreRecordNotFoundError: false,
			// 不彩色化
			Colorful: false,
		})
	// 连接服务器（池）
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 设置为自定义的日志
		Logger: customLogger,
		// 设置默认的命名策略的选项
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "wsq_",
			SingularTable: true,
			NoLowerCase:   false,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	DB = db
	// 注册序列化器
}
