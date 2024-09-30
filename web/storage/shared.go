package storage

import (
	"context"
	"database/sql"
	"log"

	"github.com/cayo-rodrigues/nff/web/utils"
)

type Scanner interface {
	Scan(args ...any) error
}

type ScannableModel interface {
	Values() []any
}

func Scan(row Scanner, models ...ScannableModel) error {
	values := []any{}
	for _, model := range models {
		values = append(values, model.Values()...)
	}
	return row.Scan(values...)
}

func WithTransaction(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Error starting transaction: ", err)
		return utils.InternalServerErr
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			log.Println("Transaction rollback due to panic: ", p)
			panic(p)  // re-throw panic after rollback
		} else if err != nil {
			tx.Rollback()
			log.Println("Transaction rollback due to error: ", err)
		} else {
			err = tx.Commit()
			if err != nil {
				log.Println("Error committing transaction: ", err)
			}
		}
	}()

	err = fn(tx)
	return err
}
