package rpc

import (
	"encoding/json"

	"github.com/tendermint/go-rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

type HTTPClient struct {
	remote   string
	endpoint string
	rpc      *rpcclient.ClientJSONRPC
	ws       *rpcclient.WSClient
}

func New(remote, wsEndpoint string) *HTTPClient {
	return &HTTPClient{
		rpc:      rpcclient.NewClientJSONRPC(remote),
		remote:   remote,
		endpoint: wsEndpoint,
	}
}

func (c *HTTPClient) Status() (*ctypes.ResultStatus, error) {
	tmResult := new(ctypes.TMResult)
	_, err := c.rpc.Call("status", []interface{}{}, tmResult)
	if err != nil {
		return nil, err
	}
	// note: panics if rpc doesn't match.  okay???
	return (*tmResult).(*ctypes.ResultStatus), nil
}

func (c *HTTPClient) ABCIInfo() (*ctypes.ResultABCIInfo, error) {
	tmResult := new(ctypes.TMResult)
	_, err := c.rpc.Call("abci_info", []interface{}{}, tmResult)
	if err != nil {
		return nil, err
	}
	return (*tmResult).(*ctypes.ResultABCIInfo), nil
}

func (c *HTTPClient) ABCIQuery(query []byte) (*ctypes.ResultABCIQuery, error) {
	tmResult := new(ctypes.TMResult)
	_, err := c.rpc.Call("abci_query", []interface{}{query}, tmResult)
	if err != nil {
		return nil, err
	}
	return (*tmResult).(*ctypes.ResultABCIQuery), nil
}

func (c *HTTPClient) ABCIProof(key []byte, height uint64) (*ctypes.ResultABCIProof, error) {
	tmResult := new(ctypes.TMResult)
	_, err := c.rpc.Call("abci_proof", []interface{}{key, height}, tmResult)
	if err != nil {
		return nil, err
	}
	return (*tmResult).(*ctypes.ResultABCIProof), nil
}

func (c *HTTPClient) BroadcastTxCommit(tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	return c.broadcastTX("broadcast_tx_commit", tx)
}

func (c *HTTPClient) BroadcastTxAsync(tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	return c.broadcastTX("broadcast_tx_async", tx)
}

func (c *HTTPClient) BroadcastTxSync(tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	return c.broadcastTX("broadcast_tx_sync", tx)
}

func (c *HTTPClient) broadcastTX(route string, tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	tmResult := new(ctypes.TMResult)
	_, err := c.rpc.Call(route, []interface{}{tx}, tmResult)
	if err != nil {
		return nil, err
	}
	return (*tmResult).(*ctypes.ResultBroadcastTxCommit), nil
}

func (c *HTTPClient) NetInfo() (*ctypes.ResultNetInfo, error) {
	tmResult := new(ctypes.TMResult)
	_, err := c.rpc.Call("net_info", nil, tmResult)
	if err != nil {
		return nil, err
	}
	return (*tmResult).(*ctypes.ResultNetInfo), nil
}

func (c *HTTPClient) DialSeeds(seeds []string) (*ctypes.ResultDialSeeds, error) {
	tmResult := new(ctypes.TMResult)
	// TODO: is this the correct way to marshall seeds?
	_, err := c.rpc.Call("dial_seeds", []interface{}{seeds}, tmResult)
	if err != nil {
		return nil, err
	}
	return (*tmResult).(*ctypes.ResultDialSeeds), nil
}

func (c *HTTPClient) BlockchainInfo(minHeight, maxHeight int) (*ctypes.ResultBlockchainInfo, error) {
	tmResult := new(ctypes.TMResult)
	_, err := c.rpc.Call("blockchain", []interface{}{minHeight, maxHeight}, tmResult)
	if err != nil {
		return nil, err
	}
	return (*tmResult).(*ctypes.ResultBlockchainInfo), nil
}

func (c *HTTPClient) Genesis() (*ctypes.ResultGenesis, error) {
	tmResult := new(ctypes.TMResult)
	_, err := c.rpc.Call("genesis", nil, tmResult)
	if err != nil {
		return nil, err
	}
	return (*tmResult).(*ctypes.ResultGenesis), nil
}

func (c *HTTPClient) Block(height int) (*ctypes.ResultBlock, error) {
	tmResult := new(ctypes.TMResult)
	_, err := c.rpc.Call("block", []interface{}{height}, tmResult)
	if err != nil {
		return nil, err
	}
	return (*tmResult).(*ctypes.ResultBlock), nil
}

func (c *HTTPClient) Validators() (*ctypes.ResultValidators, error) {
	tmResult := new(ctypes.TMResult)
	_, err := c.rpc.Call("validators", nil, tmResult)
	if err != nil {
		return nil, err
	}
	return (*tmResult).(*ctypes.ResultValidators), nil
}

/** websocket event stuff here... **/

// StartWebsocket starts up a websocket and a listener goroutine
// if already started, do nothing
func (c *HTTPClient) StartWebsocket() error {
	var err error
	if c.ws == nil {
		ws := rpcclient.NewWSClient(c.remote, c.endpoint)
		_, err = ws.Start()
		if err == nil {
			c.ws = ws
		}
	}
	return err
}

// StopWebsocket stops the websocket connection
func (c *HTTPClient) StopWebsocket() {
	if c.ws != nil {
		c.ws.Stop()
		c.ws = nil
	}
}

// GetEventChannels returns the results and error channel from the websocket
func (c *HTTPClient) GetEventChannels() (chan json.RawMessage, chan error) {
	if c.ws == nil {
		return nil, nil
	}
	return c.ws.ResultsCh, c.ws.ErrorsCh
}

func (c *HTTPClient) Subscribe(event string) error {
	return c.ws.Subscribe(event)
}

func (c *HTTPClient) Unsubscribe(event string) error {
	return c.ws.Unsubscribe(event)
}
