package clientCalculix

// ClientCalculix - RPC client of Calculix
type ClientCalculix struct {
	tasks   []task
	manager ServerManager
}

type task struct {
	inpBody  string
	ipServer string
	datBody  string
	err      error
}

/*
// Calculate - calculation
func (c *ClientCalculix) Calculate(inpBody []string) (datBody []string) {
	c.updateIPServers()
	for _, inp := range inpBody {

	}
}
*/

// GetIPServers - send list of server ip
//func (c *ClientCalculix) GetIPServers() []string {
//	c.updateIPServers()
//	return c.ipServers
//}

/*
// scan IP adresess

// saving IP servers

// create queue of calculations

// try to found free servers

// waiting

*/
