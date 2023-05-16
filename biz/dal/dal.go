package dal

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/three-body/hertz-scaffold/config"
	"gorm.io/driver/mysql"
	"gorm.io/gen/examples/dal/query"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

var (
	DB *gorm.DB
)

func InitMySQL() error {
	level := logger.Warn
	if config.GetEnv() == config.EnvDev {
		level = logger.Info
	}
	var err error
	DB, err = gorm.Open(
		mysql.New(mysql.Config{
			DSN:                       config.GetConf().Mysql.Master,
			DefaultStringSize:         uint(config.GetConf().Mysql.DefaultStringSize),   // default size for string fields
			DisableDatetimePrecision:  config.GetConf().Mysql.DisableDatetimePrecision,  // disable datetime precision, which not supported before MySQL 5.6
			SkipInitializeWithVersion: config.GetConf().Mysql.SkipInitializeWithVersion, // autoConfigure based on currently MySQL version
		}),
		&gorm.Config{
			PrepareStmt:            config.GetConf().Mysql.Gorm.PrepareStmt,
			SkipDefaultTransaction: config.GetConf().Mysql.Gorm.SkipDefaultTx,
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					SlowThreshold:             time.Duration(config.GetConf().Mysql.Gorm.SlowThreshold) * time.Millisecond,
					LogLevel:                  level,
					IgnoreRecordNotFoundError: config.GetConf().Mysql.Gorm.IgnoreRecordNotFoundError,
					Colorful:                  true,
				},
			),
		},
	)
	if err != nil {
		return fmt.Errorf("open sqlite %q fail: %w", config.GetConf().Mysql.Master, err)
	}

	if !config.GetConf().Mysql.Separation {
		sqlDB, err := DB.DB()
		if err != nil {
			return fmt.Errorf("init mysql failed: %w", err)
		}
		sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(config.GetConf().Mysql.ConnMaxIdleTime))
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.GetConf().Mysql.ConnMaxLifeTime))
		sqlDB.SetMaxIdleConns(config.GetConf().Mysql.ConnMaxIdle)
		sqlDB.SetMaxOpenConns(config.GetConf().Mysql.ConnMaxOpen)
	} else {
		sources := []gorm.Dialector{mysql.Open(config.GetConf().Mysql.Master)}
		var replicas []gorm.Dialector
		for _, dsn := range config.GetConf().Mysql.Slave {
			replicas = append(replicas, mysql.Open(dsn))
		}
		err = DB.Use(
			dbresolver.Register(dbresolver.Config{
				Sources:  sources,
				Replicas: replicas,
				// sources/replicas load balancing policy
				Policy: dbresolver.RandomPolicy{},
				// print sources/replicas mode in logger
				TraceResolverMode: true,
			}).
				SetConnMaxIdleTime(time.Second * time.Duration(config.GetConf().Mysql.ConnMaxIdleTime)).
				SetConnMaxLifetime(time.Second * time.Duration(config.GetConf().Mysql.ConnMaxLifeTime)).
				SetMaxIdleConns(config.GetConf().Mysql.ConnMaxIdle).
				SetMaxOpenConns(config.GetConf().Mysql.ConnMaxOpen))
		if err != nil {
			return fmt.Errorf("init mysql failed: %w", err)
		}
	}

	query.SetDefault(DB)
	return nil
}
