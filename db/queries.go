package db

import (
	"github.com/jmoiron/sqlx"
)

type Block struct {
	Coin       string `db:"coin"`
	Height     uint64 `db:"height"`
	Hash       string `db:"hash"`
	Difficulty float64 `db:"difficulty"`
	Work       float64 `db:"work"`
	Time       uint64 `db:"time"`
}

// GetLastBlock returns the last block (by height) found in the database
func GetLastBlock(coin string) (uint64, string, error) {
	row := struct {
		Height    uint64 `db:"height"`
		BlockHash string `db:"hash"`
	}{}

	if err := GetDB().Get(&row, "SELECT height, hash FROM blocks WHERE coin = ? ORDER BY height DESC LIMIT 1", coin); err != nil {
		if err.Error() == "sql: no rows in result set" {
			// assume db is fresh, return a block from before the fork
			return 0, "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f", nil
		}
		return 0, "", err
	}

	return row.Height, row.BlockHash, nil
}

// GetBlock will return a single block from the database with the given coin and height
func GetBlock(coin string, height uint64) (*Block, error) {
	blk := Block{}

	if err := GetDB().Get(&blk, "SELECT * FROM blocks WHERE coin = ? AND height = ?", coin, height); err != nil {
		return nil, err
	}

	return &blk, nil
}

// InsertBlock will insert a new row into the blocks table and will return the insert_id
func InsertBlock(tx *sqlx.Tx, coin, hash string, height, time uint64, diff, work float64) (int64, error) {
	qry := "INSERT INTO blocks (coin, hash, height, time, difficulty, work) VALUES(?, ?, ?, ?, ?, ?)"

	r := tx.MustExec(qry, coin, hash, height, time, diff, work)
	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// InsertDetails will insert block related statistics that came from the rpc.GetBlockStats() method
func InsertDetails(tx *sqlx.Tx, params *map[string]interface{}) (int64, error) {
	qry := "INSERT INTO details VALUES(:coin, :height, :avgfee, :avgfeerate, :avgfeerate_old, :ins, :outs, :maxfee," +
	  ":maxfeerate, :maxfeerate_old, :medianfee, :medianfeerate, :medianfeerate_old, :time, :mediantime, :minfee, :minfeerate," +
	  ":minfeerate_old, :reward, :subsidy, :totalfee, :total_size, :txs, :swtotal_size, :swtotal_weight, :swtxs, :total_out," +
	  ":total_weight, :utxo_increase, :utxo_size_inc)"

	r, err := tx.NamedExec(qry, *params)
	if err != nil {
		return 0, err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// InsertDetailsNoSegwit will insert block related statistics that came from the rpc.GetBlockStats() method
func InsertDetailsNoSegwit(tx *sqlx.Tx, params *map[string]interface{}) (int64, error) {
	qry := "INSERT INTO details (coin, height, avgfee, vavgfeerate, inputs, outputs, maxfee, vmaxfeerate, medianfee, vmedianfeerate, " +
	  "mediantime, minfee, reward,subsidy,time,totalout,totalsize,fee,txs,utxoinc,utxosizeinc) " +
	  "VALUES(:coin, :height, :avgfee, :avgfeerate_old, :ins, :outs, :maxfee, :maxfeerate_old, :medianfee, :medianfeerate_old," +
	  ":mediantime, :minfee, :reward, :subsidy, :time, :total_out, :total_size, :totalfee, :txs, :utxo_increase, :utxo_size_inc)"

	r, err := tx.NamedExec(qry, *params)
	if err != nil {
		return 0, err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetBlocksAfter returns an array of blocks that came after a certain time
func GetBlocksAfter(tx *sqlx.Tx, coin string, time uint64) (*[]Block, error) {
	blocks := make([]Block, 0, 2048)
	err := tx.Select(&blocks, "SELECT * FROM blocks WHERE coin = ? AND time >= ? ORDER BY height", coin, time)
	return &blocks, err
}

// BLockGroup is used by GetBlocksAfterGrouped()
type BlockGroup struct {
	Difficulty float64 `db:"difficulty"`
	Time       uint64 `db:"time"`
	Count      uint32 `db:"count"`
}

// GetBlocksAfterGrouped returns an array of blocks that came after a certain time grouped by the block difficulty
func GetBlocksAfterGrouped(tx *sqlx.Tx, coin string, time uint64) (*[]BlockGroup, error) {
	blocks := make([]BlockGroup, 0, 32)
	qry := "SELECT difficulty, MAX(time) as time, COUNT(*) as count FROM blocks WHERE coin = ? AND time >= ? GROUP BY difficulty ORDER BY height"
	err := tx.Select(&blocks, qry, coin, time)
	return &blocks, err
}

// InsertRates will insert hashrates for a certain coin and height
func InsertRates(tx *sqlx.Tx, coin string, height uint64, rates *map[string]float64) error {
	_, err := tx.Exec("INSERT INTO hashrates VALUES(?,?,?,?,?,?,?,?,?)",
		coin,
		height,
		(*rates)["h3"],
		(*rates)["h6"],
		(*rates)["h12"],
		(*rates)["d1"],
		(*rates)["d3"],
		(*rates)["d7"],
		(*rates)["d30"],
	)

	return err
}

// Returns the time of the last block found
func GetLastTime() (uint64, error) {
	var last uint64

	err := GetDB().Get(&last, "SELECT MAX(time) FROM blocks")
	return last, err
}
