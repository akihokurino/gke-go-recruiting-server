package cloudsql

import (
	"context"
	"database/sql"
	"errors"

	"gke-go-recruiting-server/adapter"

	gosql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB(connect string) adapter.DB {
	rawDB, err := sql.Open("mysql", connect)
	if err != nil {
		panic(err)
	}

	conn, err := gorm.Open(
		mysql.New(mysql.Config{Conn: rawDB}),
		&gorm.Config{DisableAutomaticPing: true},
	)

	return func(ctx context.Context) *gorm.DB {
		return conn.WithContext(ctx)
	}
}

func NewTX() adapter.TX {
	return func(db *gorm.DB, fn func(db *gorm.DB) error) error {
		tx := db.Begin()
		if tx.Error != nil {
			return tx.Error
		}
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		if err := fn(tx); err != nil {
			_ = tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			_ = tx.Rollback()
			return nil
		}

		return nil
	}
}

func IsDuplicateError(err error) bool {
	if mysqlErr, ok := err.(*gosql.MySQLError); ok {
		return mysqlErr.Number == 1062
	}
	return false
}

func IsNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)

}
