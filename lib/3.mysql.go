package lib

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	DB, _ = sql.Open("mysql", Config.GetString("job.mysql"))
	DB.SetMaxOpenConns(50)
	DB.SetMaxIdleConns(10)
}

type TransactionHandler func(tx *sql.Tx) error

func DoTransaction(handler TransactionHandler) error {
	tx, err := DB.Begin()
	if err != nil {
		Logger.Errorln(err.Error())
		return err
	}
	err = handler(tx)
	if err != nil {
		Logger.Errorf("Transaction rollback, %s", err.Error())
		err = tx.Rollback()
		if err != nil {
			Logger.Errorln(err.Error())
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		Logger.Errorln(err.Error())
		return err
	}
	return nil
}
