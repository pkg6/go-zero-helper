package orm

import (
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Config 数据库连接配置
//
//	DBType: 数据库类型，支持 mysql/postgres/sqlite（不区分大小写），默认为 mysql
//	DSN: 数据库连接字符串，根据 DBType 不同格式如下：
//	      - MySQL:  user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
//	      - PostgreSQL: host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai
//	      - SQLite:  ./test.db 或 :memory:
type Config struct {
	DBType          string // 数据库类型: mysql/postgres/sqlite
	DSN             string // 连接串
	MaxIdleConns    int    `json:",default=10"`    // 最大空闲连接数，默认 10
	MaxOpenConns    int    `json:",default=100"`   // 最大打开连接数，默认 100
	ConnMaxLifetime int    `json:",default=3600"` // 连接最长生命周期（秒），默认 3600
}

// MustOpen 打开数据库连接，失败时 panic
//
//	使用场景：初始化阶段或测试代码
//
//	示例：
//
//	  db := orm.MustOpen(&orm.Config{DBType: "mysql", DSN: "root:root@tcp(localhost:3306)/test"})
func MustOpen(conf *Config) *gorm.DB {
	db, err := Open(conf)
	if err != nil {
		panic(err)
	}
	return db
}

// Open 打开数据库连接，支持 MySQL、PostgreSQL、SQLite
//
//	默认参数：
//	  - MaxIdleConns: 10
//	  - MaxOpenConns: 100
//	  - ConnMaxLifetime: 3600（秒）
//
//	示例：
//
//	  // MySQL
//	  db, _ := orm.Open(&orm.Config{DBType: "mysql", DSN: "root:root@tcp(localhost:3306)/test"})
//
//	  // PostgreSQL
//	  db, _ := orm.Open(&orm.Config{
//	      DBType: "postgres",
//	      DSN:    "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable",
//	  })
//
//	  // SQLite
//	  db, _ := orm.Open(&orm.Config{DBType: "sqlite", DSN: "./test.db"})
func Open(conf *Config) (*gorm.DB, error) {
	// 设置默认连接池参数
	if conf.MaxIdleConns == 0 {
		conf.MaxIdleConns = 10
	}
	if conf.MaxOpenConns == 0 {
		conf.MaxOpenConns = 100
	}
	if conf.ConnMaxLifetime == 0 {
		conf.ConnMaxLifetime = 3600
	}

	// 根据数据库类型创建对应的 Dialector
	var dialector gorm.Dialector
	switch strings.ToLower(conf.DBType) {
	case "sqlite":
		dialector = sqlite.Open(conf.DSN)
	case "postgres":
		dialector = postgres.Open(conf.DSN)
	default:
		dialector = mysql.Open(conf.DSN)
	}

	// 打开数据库连接
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: &Logger{},
	})
	if err != nil {
		return nil, err
	}

	// 获取底层 sql.DB 实例
	sdb, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 配置连接池
	sdb.SetMaxIdleConns(conf.MaxIdleConns)
	sdb.SetMaxOpenConns(conf.MaxOpenConns)
	sdb.SetConnMaxLifetime(time.Second * time.Duration(conf.ConnMaxLifetime))

	return db, nil
}
