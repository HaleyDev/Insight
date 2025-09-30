package data

import (
	"fmt"
	c "insight/config"
	log "insight/internal/pkg/logger"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var MysqlDB *gorm.DB

type Writer interface {
	Printf(string, ...interface{})
}

type WriterLog struct{}

func (w WriterLog) Printf(format string, args ...interface{}) {
	if c.GetConfig().MySQL.PrintSql {
		log.Logger.Sugar().Infof(format, args...)
	}
}

func GetLoggerLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info
	}
}

func initMysql() {
	logConfig := logger.New(
		WriterLog{},
		logger.Config{
			SlowThreshold:             0,
			LogLevel:                  GetLoggerLevel(c.GetConfig().MySQL.LogLevel), // Log level
			IgnoreRecordNotFoundError: false,                                        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                                         // Disable color
		},
	)

	configs := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: c.GetConfig().MySQL.TablePrefix, // 表名前缀
		},
		Logger: logConfig,
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.GetConfig().MySQL.Username,
		c.GetConfig().MySQL.Password,
		c.GetConfig().MySQL.Host,
		c.GetConfig().MySQL.Port,
		c.GetConfig().MySQL.Database,
	)

	var err error
	MysqlDB, err = gorm.Open(mysql.Open(dsn), configs)

	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}

	sqlDB, _ := MysqlDB.DB()
	sqlDB.SetMaxIdleConns(c.GetConfig().MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.GetConfig().MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.GetConfig().MySQL.MaxLifetime) * time.Second)
	log.Logger.Info("MySQL connected successfully")
}
