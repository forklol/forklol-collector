package main

import (
	"os"
	"flag"
	"fmt"
	"forklol-collector/config"
	"forklol-collector/db"
	"forklol-collector/bitcoin"
	"time"
	"forklol-collector/stats"
	"runtime"
)

var coins []bitcoin.Coin

func main() {
	Init()
	stats.LoadPresets()
	stats.StartDispatcher(runtime.NumCPU() - 2)
	db.InitDB(config.Options().DB_CONNECTION_STRING)

	// FIXME: get from config file
	coins = []bitcoin.Coin{
		{
			Symbol:   "BTC",
			RPCUrl:   "http://192.168.1.2:8332/",
			RPCUser:  "forklol",
			RPCPass:  "forklol",
			RPCStats: true,
			SegWit:   true,
		},
		{
			Symbol:   "BCH",
			RPCUrl:   "http://192.168.1.2:8331/",
			RPCUser:  "forklol",
			RPCPass:  "forklol",
			RPCStats: true,
			SegWit:   false,
		},
	}

	syncers := make([]bitcoin.ChainSync, 0)
	for _, c := range coins {
		syncers = append(syncers, bitcoin.NewChainSync(c))
	}

	RunSyncers(syncers)
}

func RunSyncers(syncers []bitcoin.ChainSync) {
	t := time.NewTicker(time.Second * 5)

	for {
		done := make(chan bool)
		select {
		case <-t.C:
			t.Stop()
			for _, sync := range syncers {
				go sync.Sync(done)
			}
		}

		changed := false
		for n := 0; n < len(syncers); n++ {
			changed = <-done || changed
		}
		close(done)

		if changed {
			stats.Package(coins)
		}

		t = time.NewTicker(time.Second * 5)
	}
}

// init parses arguments and sets config.Options
func Init() {
	// env values
	env_pubkey, _ := os.LookupEnv("FORKLOL_BTCAVG_PUBKEY")
	env_secret, _ := os.LookupEnv("FORKLOL_BTCAVG_SECRET")

	env_dbuser, _ := os.LookupEnv("FORKLOL_DB_USER")
	env_dbpass, _ := os.LookupEnv("FORKLOL_DB_PASS")
	env_dbhost, _ := os.LookupEnv("FORKLOL_DB_HOST")
	env_dbport, _ := os.LookupEnv("FORKLOL_DB_PORT")
	env_dbschm, _ := os.LookupEnv("FORKLOL_DB_SCHEME")

	// set argument flags
	pub := flag.String("pubkey", env_pubkey, "bitcoinaverage.com api public key, defaults to env var FORKLOL_BTCAVG_PUBKEY")
	sec := flag.String("secret", env_secret, "bitcoinaverage.com api secret, defaults to env var FORKLOL_BTCAVG_SECRET")
	dbg := flag.Bool("debug", false, "enable debugging")

	dbuser := flag.String("dbuser", env_dbuser, "mysql user")
	dbpass := flag.String("dbpass", env_dbpass, "mysql password")
	dbhost := flag.String("dbhost", env_dbhost, "mysql host/address")
	dbport := flag.String("dbport", env_dbport, "mysql port")
	dbscheme := flag.String("dbscheme", env_dbschm, "mysql database name/scheme")

	flag.Parse()

	// set config.Optios
	opts := config.Options()

	opts.DEBUG = *dbg
	opts.DB_CONNECTION_STRING = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", *dbuser, *dbpass, *dbhost, *dbport, *dbscheme)
	opts.BTCAVG_PUBKEY = *pub
	opts.BTCAVG_SECRET = *sec
}
