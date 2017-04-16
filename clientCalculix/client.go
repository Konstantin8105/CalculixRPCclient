package clientCalculix

import (
	"fmt"
	"net/rpc"
	"time"

	"github.com/Konstantin8105/CalculixRPCserver/serverCalculix"
)

// ClientCalculix - RPC client of Calculix
type ClientCalculix struct {
	Manager ServerManager
}

func (c *ClientCalculix) getServer() (client *rpc.Client, err error) {
	addresses := c.Manager.GetIPServers()
	for {
		for _, address := range addresses {
			client, err = rpc.DialHTTP("tcp", address)
			if err != nil {
				return nil, fmt.Errorf("Cannot run : %v", err)
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
		time.Sleep(time.Second * 2)
	}
RUN:
	return client, err
}
