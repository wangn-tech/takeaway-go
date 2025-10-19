package database

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"takeaway-go/internal/app/config"
	"takeaway-go/pkg/logger"
)

var DB *gorm.DB

func InitDB() {
	// 基于 zap 实现 gorm 日志接口 logger.Interface
	ormLogger := logger.NewGormZapLogger(logger.Log)
	// 根据环境设置 GORM 的日志级别
	if gin.Mode() == "debug" {
		ormLogger.LogLevel = gormlogger.Info
	} else {
		ormLogger.LogLevel = gormlogger.Warn
	}

	// 构建 DSN 字符串
	conf := config.AppConf.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)

	// 连接数据库
	var err error
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	logger.Log.Info("MySQL database connected successfully")
}
