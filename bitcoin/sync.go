package bitcoin

import (
	"log"
	"forklol-collector/db"
	"forklol-collector/rpc"
)

type ChainSync struct {
	Coin Coin
}

func NewChainSync(coin Coin) ChainSync {
	return ChainSync{
		Coin: coin,
	}
}

// SyncAll brings the database up to date with the bitcoind chain. It does not handle reorgs.
func (c ChainSync) Sync() {
	prevHeight, prevHash, err := db.GetLastBlock(c.Coin.Symbol)
	if err != nil {
		log.Printf("Could not get last %s block from database: %s\n", c.Coin.Symbol, err.Error())
		return
	}

	client := c.Coin.RPCClient()

	height, hash, err := client.GetLastBlock()
	if err != nil {
		log.Printf("Could not get last block from %s chain: %s.\n", c.Coin.Symbol, err.Error())
		return
	}

	if prevHeight < height {
		// @TODO check for reorg
		log.Printf("Syncing %s chain to block %d (from %d, %d blocks)\n", c.Coin.Symbol, height, prevHeight, height-prevHeight)

		c.syncFromHeight(prevHeight, height)

	} else if prevHash != hash {
		// @TODO check for reorg or new best
	} else {
		// no new block(s) found
	}
}

// syncFromHeight will get new blocks from bitcoind and pass them to handleNewBlock for processing
func (c ChainSync) syncFromHeight(prevHeight, height uint32) {
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

		log.Printf("Handling new %s block %d, %s (%d left)", c.Coin.Symbol, h, hash, height-h)
		c.handleNewBlock(block)

	}
}

// handleNewBlock will insert the block into the database and collect some more information about it after
func (c ChainSync) handleNewBlock(block *rpc.Block) error {
	tx := db.GetDB().MustBegin()

	_, err := db.InsertBlock(
		tx,
		c.Coin.Symbol,
		block.Hash,
		block.Height,
		block.Size,
		block.Weight,
		block.Time,
		block.Difficulty,
	)

	if err != nil {
		log.Printf("Could not insert block into database: %s\n", err.Error())
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Could not commit db transactions.\n")
		return err
	}

	return nil
}
