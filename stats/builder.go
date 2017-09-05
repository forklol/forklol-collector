package stats

import (
	"forklol-collector/bitcoin"
)

type StatBuilder struct {
	Coin bitcoin.Coin
}

// NewStatBuilder return a new StatBuilder
func NewStatBuilder(coin bitcoin.Coin) *StatBuilder {
	return &StatBuilder{
		Coin: coin,
	}
}

// GetStatByPreset returns []Value based on a predefined preset from stats.presets.go
func (b StatBuilder) GetStatByPreset(preset StatPreset, ctype CompactType, from, to, step int64) (*[]Value, error) {
	s := b.getStatisticFetcherFromKind(preset)
	return s.GetValues(s.Compacter(ctype, from, to, step), preset.Method, preset.Type), nil
}

func (b StatBuilder) getStatisticFetcherFromKind(preset StatPreset) StatisticFetcher {
	switch preset.Kind {
	case KIND_BLOCK:
		return NewBlockStatistic(b.Coin.Symbol, preset.Property)
	case KIND_DETAIL:
		return NewDetailStatistic(b.Coin.Symbol, preset.Property)
	case KIND_HASHRATE:
		return NewHashrateStatistic(b.Coin.Symbol, preset.Property)
	default:
		panic("Kind not found")
	}
}

func GetStepSize(from, to, parts int64) int64 {
	return (to - from) / parts
}

// GetCompacterIndex returns an array of integeres representing every step in the compacted set
func GetCompacterIndex(from, to, step int64) *[]int64 {
	compacter := newCompacter(COMPACT_HEIGHT, from, to, step, "b")

	index := make([]int64, 0)
	for n := compacter.From; n < compacter.To; n += compacter.StepSize {
		index = append(index, n)
	}

	return &index
}
