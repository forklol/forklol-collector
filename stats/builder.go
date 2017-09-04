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

type StatKind string

const (
	KIND_BLOCK StatKind = "block"
	KIND_DETAIL StatKind = "detail"
	KIND_HASHRATE StatKind = "hashrate"
	KIND_PRICE StatKind = "price"
)

func (b StatBuilder) GetStatValues(kind StatKind, prop string, from, to, step uint64) *[]Value {

}

