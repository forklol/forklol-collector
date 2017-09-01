package db

import (
	"github.com/jmoiron/sqlx"
	"errors"
)

// GetLastBlock returns the last block (by height) found in the database
func GetLastBlock(coin string) (uint32, string, error) {
	row := struct {
		Height    uint32 `db:"height"`
		BlockHash string `db:"hash"`
	}{}

	if err := GetDB().Get(&row, "SELECT height, hash FROM blocks WHERE coin = ? ORDER BY height DESC LIMIT 1", coin); err != nil {
		if err.Error() == "sql: no rows in result set" {
			// assume db is fresh, return a block from before the fork
			return 478001, "000000000000000001211f05ea9cf34d5a9710ca2ab7572b9d79273924e030ae", nil
		}
		return 0, "", err
	}

	return row.Height, row.BlockHash, nil
}

// InsertBlock will insert a new row into the blocks table and will return the insert_id
func InsertBlock(tx *sqlx.Tx, coin, hash string, height, size, weight, time uint32, diff float64) (int64, error) {
	qry := "INSERT INTO blocks (coin, hash, height, size, weight, time, difficulty) VALUES(?, ?, ?, ?, ?, ?, ?)"

	if tx == nil {
		return 0, errors.New("No transaction object provided, cannot insert row into blocks table.")
	}

	r := tx.MustExec(qry, coin, hash, height, size, weight, time, diff)
	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
