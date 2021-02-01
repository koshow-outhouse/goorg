// This file was generated by goorg/cli/repository at 2021/02/01 11:50:04

package infra

import (
	"context"
	"log"

	"github.com/organization-service/goorg/database"
	"gorm.io/gorm"
)

type (
	dbRepository struct {
		driver database.IDriver
	}

	dbConnection struct {
		driver database.IDriver
	}

	dbTransaction struct {
		connectionReadWrite *gorm.DB
	}
)

// @repository
func NewRepository(driver database.IDriver) repositories.Repository {
	return &dbRepository{
		driver: driver,
	}
}

func (r *dbRepository) NewConnection() (repositories.Connection, error) {
	return &dbConnection{
		driver: r.driver,
	}, nil
}

func (r *dbRepository) MustConnection() repositories.Connection {
	db, err := r.NewConnection()
	if err != nil {
		panic(err)
	}
	return db
}

func (con *dbConnection) Transaction(c context.Context, f func(tx repositories.Transaction) error) error {
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()
	tx := con.driver.ReadWriteConnection(c)
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			panic(err)
		}
	}()
	err = f(&dbTransaction{
		connectionReadWrite: tx,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
