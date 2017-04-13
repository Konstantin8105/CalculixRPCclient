package clientCalculix

// ClientCalculix - RPC client of Calculix
type ClientCalculix struct {
	tasks       []task
	IPPrototype string   // example: "192.168.5." or "192.168.0."
	ipServers   []string // IP allowable servers
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


args := &server.Args{7,8}
var reply int
err = client.Call("Arith.Multiply", args, &reply)
if err != nil {
	log.Fatal("arith error:", err)
}
fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)

// Asynchronous call
quotient := new(Quotient)
divCall := client.Go("Arith.Divide", args, quotient, nil)
replyCall := <-divCall.Done	// will be equal to divCall
// check errors, print, etc.


*/
