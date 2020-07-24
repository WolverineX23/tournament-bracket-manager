/*

 */

package models

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jinzhu/gorm/dialects/sqlite"
)

type DB struct {
	DB               *gorm.DB
	dbType           string
	dbConnectionPath string
}

// NewDB initializes a DB object
func NewDB(dpType, dbConnectionPath string) *DB {
	return &DB{
		dbType:           dpType,
		dbConnectionPath: dbConnectionPath,
	}
}

// Connect initiates a new connection with the given connection parameters
// for the database. It returns an error in case the connection fails.
func (db *DB) Connect() error {
	var dbConnection *gorm.DB
	var err error
	if db.dbType == "sqlite3" {
		dbConnection, err = gorm.Open(sqlite.Open(db.dbConnectionPath), &gorm.Config{})
	} else if db.dbType == "postgres" {
		dbConnection, err = gorm.Open(postgres.Open(db.dbConnectionPath), &gorm.Config{})
	} else {
		err = errors.New("invalid dbtype")
	}
	if err != nil {
		return fmt.Errorf("failed to initialize database, got error: %v", err)
	}
	dbConnection.AutoMigrate(&Match{})
	//dbConnection.LogMode(true)
	db.DB = dbConnection
	return err
}
