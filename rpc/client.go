package rpc

import (
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"errors"
)

type call struct {
	JsonRPC string `json:"jsonrpc"`
	Id      string `json:"id"`
	Method  string `json:"method"`
	Params  interface{} `json:"params"`
}

type Client struct {
	url  string
	user string
	pass string
}

// NewClient returns a new Client object to call rpc methods with
func NewClient(url, user, pass string) Client {
	return Client{
		url:  url,
		user: user,
		pass: pass,
	}
}

// Call performs a method call to bitcoind rpc and returns the result as a *[]byte
func (c Client) Call(method string, params interface{}) (*[]byte, error) {
	call := call{
		JsonRPC: "1.0",
		Id:      "forklol",
		Method:  method,
		Params:  params,
	}

	j, err := json.Marshal(call)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(j))
	req.SetBasicAuth(c.user, c.pass)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 500 {
		return nil, errors.New("bad request")
	}

	return &body, nil
}

// GetLastBlock returns the height and blockhash of the last (or "best") block
func (c Client) GetLastBlock() (uint64, string, error) {
	j, err := c.Call("getblockchaininfo", []string{})
	if err != nil {
		return 0, "", err
	}

	t := struct {
		Result struct {
			Height    uint64 `json:"blocks"`
			BlockHash string `json:"bestblockhash"`
		} `json:"result"`
	}{}

	if err := json.Unmarshal(*j, &t); err != nil {
		return 0, "", err
	}

	return t.Result.Height, t.Result.BlockHash, nil
}

// GetBlockHash returns the blockhash of a block at the specified height
func (c Client) GetBlockHash(height uint64) (string, error) {
	params := []uint64{height}

	j, err := c.Call("getblockhash", params)
	if err != nil {
		return "", err
	}

	t := struct {
		Hash string `json:"result"`
	}{}

	if err := json.Unmarshal(*j, &t); err != nil {
		return "", err
	}

	return t.Hash, nil
}

// Block is used by rpc.GetBlock()
type Block struct {
	Height     uint64 `json:"height"`
	Hash       string `json:"hash"`
	Size       uint64 `json:"size"`
	Weight     uint64 `json:"weight"`
	Time       uint64 `json:"time"`
	MedianTime uint64 `json:"mediantime"`
	Difficulty float64 `json:"difficulty"`
}

// GetBlock returns a some basic information about a block with the given blockhash
func (c Client) GetBlock(blockhash string) (*Block, error) {
	j, err := c.Call("getblock", []string{blockhash})
	if err != nil {
		return nil, err
	}

	t := struct {
		Result Block `json:"result"`
	}{}

	if err = json.Unmarshal(*j, &t); err != nil {
		return nil, err
	}

	return &t.Result, nil
}

// BlockStats is used by rpc.GetBlockStats()
type BlockStats map[string][]interface{}

// GetBlockStats will return statistics about the block with the given height. RPC method "getblockstats" will be used.
func (c Client) GetBlockStats(height uint64) (*BlockStats, error) {
	params := []uint64{height, height}

	j, err := c.Call("getblockstats", params)
	if err != nil {
		return nil, err
	}

	t := struct {
		Result BlockStats `json:"result"`
	}{}

	if err := json.Unmarshal(*j, &t); err != nil {
		return nil, err
	}

	return &t.Result, nil
}
