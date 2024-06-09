package attackclient

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/tsinghua-cel/attacker-service/rpc"
	"github.com/tsinghua-cel/attacker-service/types"
	"sync/atomic"
)

// Client defines typed wrappers for the Ethereum RPC API.
type Client struct {
	c      *rpc.Client
	uuid   string
	valIdx int
	info   atomic.Value
}

// Dial connects a client to the given URL.
func Dial(rawurl string, valIdx int) (*Client, error) {
	return DialContext(context.Background(), rawurl, valIdx)
}

// DialContext connects a client to the given URL with context.
func DialContext(ctx context.Context, rawurl string, valIdx int) (*Client, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return NewClient(c, valIdx), nil
}

// NewClient creates a client that uses the given RPC client.
func NewClient(c *rpc.Client, valIdx int) *Client {
	client := &Client{
		c:      c,
		uuid:   uuid.NewString(),
		valIdx: valIdx,
	}
	info := client.clientInfo()
	client.info.Store(info)
	return client
}

// Close closes the underlying RPC connection.
func (ec *Client) Close() {
	ec.c.Close()
}

// Client gets the underlying RPC client.
func (ec *Client) Client() *rpc.Client {
	return ec.c
}

func (ec *Client) clientInfo() string {
	if v := ec.info.Load(); v != nil {
		return v.(string)
	}

	info := types.ClientInfo{
		UUID:           ec.uuid,
		ValidatorIndex: ec.valIdx,
	}
	d, _ := json.Marshal(info)
	return string(d)

}
