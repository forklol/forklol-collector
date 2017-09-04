package bitcoin

import (
	"log"
	"forklol-collector/db"
	"forklol-collector/rpc"
	"github.com/jmoiron/sqlx"
	"sync"
	"errors"
	"time"
)

type ChainSync struct {
	Coin   Coin
	TxLock sync.Mutex
}

func NewChainSync(coin Coin) ChainSync {
	return ChainSync{
		Coin: coin,
	}
}

// SyncAll brings the database up to date with the bitcoind chain. It does not handle reorgs.
func (c ChainSync) Sync(done chan bool) {
	prevHeight, prevHash, err := db.GetLastBlock(c.Coin.Symbol)
	if err != nil {
		log.Printf("Could not get last %s block from database: %s\n", c.Coin.Symbol, err.Error())
		done <- false
		return
	}

	client := c.Coin.RPCClient()

	height, hash, err := client.GetLastBlock()
	if err != nil {
		log.Printf("Could not get last block from %s chain: %s.\n", c.Coin.Symbol, err.Error())
		done <- false
		return
	}

	if prevHeight < height {
		// @TODO check for reorg
		log.Printf("Syncing %s chain to block %d (from %d, %d blocks)\n", c.Coin.Symbol, height, prevHeight, height-prevHeight)

		c.syncFromHeight(prevHeight, height)
		done <- true
		return
	} else if prevHash != hash {
		// @TODO check for reorg or new best
	} else {
		// no new block(s) found
	}

	done <- false
}

// syncFromHeight will get new blocks from bitcoind and pass them to handleNewBlock for processing
func (c ChainSync) syncFromHeight(prevHeight, height uint64) {
	for h := prevHeight + 1; h <= height; h++ {
		hash, err := c.Coin.RPCClient().GetBlockHash(h)
		if err != nil {
			log.Printf("Error getting %s blockhash at height %d, aborting sync at this height point.\n", c.Coin.Symbol, h)
			break
		}

		// get block info based on hash
		block, err := c.Coin.RPCClient().GetBlock(hash)
		if err != nil {
			log.Printf("Error getting %s block info at height %d, aborting sync at this height point.\n", c.Coin.Symbol, h)
			break
		}

		log.Printf("\u2794 Handling new %s block %d, %s (%d left)", c.Coin.Symbol, h, hash, height-h)
		start := time.Now()
		err = c.handleNewBlock(block)
		end := time.Now()
		if err != nil {
			log.Printf("\u2718 Error handling %s block %d, skipping other blocks\n", c.Coin.Symbol, block.Height)
			break
		}
		log.Printf("\u2714 New %s block %d handled in %s\n", c.Coin.Symbol, block.Height, end.Sub(start))

	}
}

// handleNewBlock will insert the block into the database and collect some more information about it after
func (c ChainSync) handleNewBlock(block *rpc.Block) error {
	done := make(chan *collectResult)
	if c.Coin.RPCStats == true {
		go c.asyncCollectStats(done, block.Height)
	} else {
		go func() {
			done <- nil
		}()
	}

	tx := db.GetDB().MustBegin()

	prevBlock, err := db.GetBlock(c.Coin.Symbol, block.Height-1)
	work := float64(0.0)

	if err == nil {
		work = prevBlock.Work
	}

	_, err = db.InsertBlock(
		tx,
		c.Coin.Symbol,
		block.Hash,
		block.Height,
		block.Time,
		block.Difficulty,
		work+block.Difficulty,
	)

	if err != nil {
		<-done
		close(done)
		log.Printf("Could not insert block into database: %s\n", err.Error())
		return err
	}

	rates, err := c.determineHashrates(tx, block.Time)
	if err != nil {
		<-done
		close(done)
		log.Printf("Could not determine hashrate of %s block %d\n", c.Coin.Symbol, block.Height)
		tx.Rollback()
		return err
	}

	if err := db.InsertRates(tx, c.Coin.Symbol, block.Height, rates); err != nil {
		<-done
		close(done)
		log.Printf("Could not insert hashrates of %s block %d\n", c.Coin.Symbol, block.Height)
		tx.Rollback()
		return err
	}

	if c.Coin.RPCStats == true {
		select {
		case stats := <-done:
			if stats == nil {
				tx.Rollback()
				return errors.New("Could not get stats through RPC.")
			}

			if c.Coin.SegWit {
				_, err = db.InsertDetails(tx, stats.Stats)
			} else {
				_, err = db.InsertDetailsNoSegwit(tx, stats.Stats)
			}

			if err != nil {
				log.Printf("Could not insert block details database: %s\n", err.Error())
				return err
			}
		}
		close(done)
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Could not commit db transactions: %s.\n", err.Error())
		return err
	}

	return nil
}

func (c ChainSync) asyncCollectStats(done chan *collectResult, height uint64) {
	stats, err := c.collectBlockStats(height)
	if err != nil {
		done <- nil
		return
	}

	done <- &collectResult{
		Stats: stats,
		Error: nil,
	}
}

type collectResult struct {
	Stats *map[string]interface{}
	Error error
}

// collectBlockStats will call rpc.GetBlockStats() and return the statistics for the block at the given height.
func (c ChainSync) collectBlockStats(height uint64) (*map[string]interface{}, error) {
	stats, err := c.Coin.RPCClient().GetBlockStats(height)
	if err != nil {
		return nil, err
	}

	flat := map[string]interface{}{
		"coin": c.Coin.Symbol,
	}

	for key, val := range *stats {
		flat[key] = val[0]
	}

	return &flat, nil
}

func (c ChainSync) determineHashrates(tx *sqlx.Tx, time uint64) (*map[string]float64, error) {
	avgs := map[string]uint64{
		"h3":  time - 3*3600,
		"h6":  time - 6*3600,
		"h12": time - 12*3600,
		"d1":  time - 24*3600,
		"d3":  time - 3*24*3600,
		"d7":  time - 7*24*3600,
		"d30": time - 30*24*3600,
	}

	rates := map[string]float64{}

	for avg, t := range avgs {
		blocks, err := db.GetBlocksAfterGrouped(tx, c.Coin.Symbol, t)
		if err != nil {
			return nil, err
		}

		diffs := make([]db.BlockGroup, 0, 8)
		num := map[float64]uint32{}
		totalwork := float64(0.0)

		// track and compensate for difficulty adjustments
		for _, blk := range *blocks {
			diffs = append(diffs, blk)
			num[blk.Difficulty] = blk.Count
			totalwork += (blk.Difficulty * float64(blk.Count))
		}

		n := 0;
		lastT := t
		compensatedRate := 0.0

		// combine all blocks during difficulty era's into one estimate (factored by how long the era was active)
		for _, blk := range diffs {
			n++

			timeTaken := blk.Time - lastT
			blocksExpected := float64(timeTaken) / 600.0
			blocksActually := num[blk.Difficulty]

			factor := (float64(timeTaken) / float64(time-t))
			compensatedRate += float64(blocksActually) / blocksExpected * blk.Difficulty * factor
			lastT = (blk).Time
		}

		rates[avg] = compensatedRate
	}

	return &rates, nil
}
