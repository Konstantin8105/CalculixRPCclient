package clientCalculix

import (
	"net/rpc"
	"time"

	"github.com/Konstantin8105/CalculixRPCclient/serverManager"
	"github.com/Konstantin8105/CalculixRPCserver/serverCalculix"
)

// ClientCalculix - RPC client of Calculix
type ClientCalculix struct {
	manager serverManager.ServerManager
}

// NewClient - create new client for calculation
func NewClient() (client *ClientCalculix) {
	client = new(ClientCalculix)
	client.manager = *serverManager.NewServerManager()
	return client
}

func (c *ClientCalculix) getServer() (client *rpc.Client, err error) {
	addresses := c.manager.GetIPServers()
	for {
		for _, address := range addresses {
			client, err = rpc.DialHTTP("tcp", address)
			if err != nil {
				c.manager.UpdateServers()
				return c.getServer()
			}
			var amount serverCalculix.Amount
			err = client.Call("Calculix.AmountTasks", "", &amount)
			if err != nil {
				return nil, err
			}
			if amount.A > 0 {
				goto RUN
			}
			err = client.Close()
			if err != nil {
				return nil, err
			}
		}
		time.Sleep(2 * time.Second)
	}
RUN:
	return client, err
}
