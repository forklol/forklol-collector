package stats

import "forklol-collector/bitcoin"

type StatBuilder struct {
	Coin bitcoin.Coin
}

// NewStatBuilder return a new StatBuilder
func NewStatBuilder(coin bitcoin.Coin) *StatBuilder {
	return &StatBuilder{
		Coin: coin,
	}
}

func (b StatBuilder) CreateStat()
