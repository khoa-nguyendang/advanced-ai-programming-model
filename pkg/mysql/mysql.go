package mysql

import (
	"aapi/config"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

const (
	maxOpenConns    = 60
	connMaxLifetime = 120
	maxIdleConns    = 30
	connMaxIdleTime = 20
)

// Return new Mysql db instance
func NewMysqlDB(c *config.Config) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true",
		c.Mysql.MysqlUser,
		c.Mysql.MysqlPassword,
		c.Mysql.MysqlHost,
		c.Mysql.MysqlPort,
		c.Mysql.MysqlDbname,
	)
	log.Default().Printf("mysql datasoure: %v", dataSourceName)
	db, err := sqlx.Connect(c.Mysql.MysqlDriver, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigration(c *config.Config) error {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true",
		c.Mysql.MysqlUser,
		c.Mysql.MysqlPassword,
		c.Mysql.MysqlHost,
		c.Mysql.MysqlPort,
		c.Mysql.MysqlDbname,
	)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Default().Printf("RunMigration.error1: %v", err)
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Default().Printf("RunMigration.error2: %v", err)
	}
	log.Default().Printf("RunMigration: %v", dataSourceName)
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "mysql", driver)
	if err != nil {
		log.Default().Printf("RunMigration.error3: %v", err)
		return err
	}
	m.Up()

	return err
}
