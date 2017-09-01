package bitcoin

import "forklol-collector/rpc"

type Coin struct {
	Symbol  string
	RPCUrl  string
	RPCUser string
	RPCPass string

	rpc *rpc.Client
}

func (c Coin) RPCClient() *rpc.Client {
	if c.rpc == nil {
		client := rpc.NewClient(c.RPCUrl, c.RPCUser, c.RPCPass)
		c.rpc = &client
	}
	return c.rpc
}
