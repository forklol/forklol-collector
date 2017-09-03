package config

const (
	CHAINSPLIT_TIMESTAMP = 1501593374
	CHAINSPLIT_WORK = 32729585000856628.00
	CHAINSPLIT_HEIGHT = 478558
)

type options struct {
	DEBUG                bool
	DB_CONNECTION_STRING string
	BTCAVG_PUBKEY        string
	BTCAVG_SECRET        string

	RPC_BTC  string
	RPC_BCH  string
	RPC_TBTC string
	RPC_ELM  string
	RPC_LQD  string
}

var opts options

func Options() *options {
	return &opts
}
