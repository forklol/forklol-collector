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
func (b StatBuilder) GetStatByPreset(preset StatPreset, ctype CompactType, from, to, step uint64) (*[]Value, error) {
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

func (b StatBuilder) GetStepSize(from, to, parts uint64) uint64 {
	return (to - from) / parts
}
