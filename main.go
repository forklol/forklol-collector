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
	"encoding/json"
	"io/ioutil"
)

var coins []bitcoin.Coin

func main() {
	Init()
	stats.LoadPresets()
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

	if config.Options().DEBUG {
		now := uint64(time.Now().Unix())

		builder := stats.NewStatBuilder(coins[0])
		stepSize := builder.GetStepSize(now-(30*24*3600), now, 144)

		allStats := make([]*[]stats.Value, 0)

		for n, v := range (*stats.GetPresets()) {
			stats, _ := builder.GetStatByPreset(v, stats.COMPACT_TIME, now-(30*24*3600), now, stepSize)
			allStats = append(allStats, stats)

		}
		j, _ := json.Marshal(allStats)
		ioutil.WriteFile("/tmp/forklol2.json", j, 0644)
		return
	}

	syncers := make([]bitcoin.ChainSync, 0)

	done := make(chan bool)

	// initial sync
	for _, coin := range coins {
		sync := bitcoin.NewChainSync(coin)
		go sync.Sync(done)

		syncers = append(syncers, sync)
	}

	for n := 0; n < len(syncers); n++ {
		<-done
	}
	close(done)

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

		for n := 0; n < len(syncers); n++ {
			<-done
		}
		close(done)
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
