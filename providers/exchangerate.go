package providers

import (
	"strconv"
	"time"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"io/ioutil"
	"log"
	"fmt"
	"encoding/json"
	"errors"
	"forklol-collector/config"
	"forklol-collector/db"
)

type ExchangeRateFetcher struct {
	symbol string
	pubkey string
	secret string
}

func NewExchangeRateFetcher(symbol string) *ExchangeRateFetcher {
	pf := ExchangeRateFetcher{}
	pf.init(symbol)

	return &pf
}

// GetExchangeRate returns the exchange rate of a coin at a certain point in time
func (p ExchangeRateFetcher) GetExchangeRate(height, timestamp uint64) (float64, error) {
	if height <= config.CHAINSPLIT_HEIGHT {
		return 0.0, nil
	}

	er := struct {
		ExchangeRate float64 `db:"er"`
	}{}

	if err := db.GetDB().Get(&er, "SELECT price FROM prices WHERE coin = ? AND height = ?", p.symbol, height); err != nil {
		log.Printf("No db price found, fetching from BitcoinAverage.com\n")

		a, err := p.fetchHistoricalExchangeRate(timestamp)
		if err != nil {
			log.Printf("Price could be fetched for %s @ %d\n", p.symbol, height)
			return 0.0, err
		}

		p.setPrice(height, a)
		er.ExchangeRate = a
	}

	return er.ExchangeRate, nil
}

func (p ExchangeRateFetcher) init(sym string) {
	p.symbol = sym
	p.pubkey = config.Options().BTCAVG_PUBKEY
	p.secret = config.Options().BTCAVG_SECRET
}

func (p ExchangeRateFetcher) getFromBitcoinAverage(url string) ([]byte, error) {
	payload := strconv.Itoa(int(time.Now().Unix())) + "." + p.pubkey

	mac := hmac.New(sha256.New, []byte(p.secret))
	mac.Write([]byte(payload))
	digest := hex.EncodeToString(mac.Sum(nil))
	sig := payload + "." + digest

	c := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-signature", sig)

	resp, err := c.Do(req)
	if err != nil {
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return body, nil
}

func (p ExchangeRateFetcher) setPrice(height uint64, price float64) {
	db.GetDB().Exec("INSERT INTO prices (coin, height, price) VALUES(?,?,?)", p.symbol, height, price)
}

func (p ExchangeRateFetcher) fetchHistoricalExchangeRate(timestamp uint64) (float64, error) {
	url := fmt.Sprintf("https://apiv2.bitcoinaverage.com/indices/global/history/%sUSD?at=%d", p.symbol, timestamp)
	body, err := p.getFromBitcoinAverage(url)

	data := struct {
		Average float64 `json:"average"`
	}{}

	if config.Options().DEBUG == true {
		time.Sleep(time.Millisecond * 750)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0.0, errors.New("Could not decode price response.")
	}

	return data.Average, nil
}
