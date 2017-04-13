package clientCalculix

import (
	"fmt"
	"net/rpc"
	"sync"

	"github.com/Konstantin8105/CalculixRPCserver/serverCalculix"
)

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
func (c *ClientCalculix) GetIPServers() []string {
	c.updateIPServers()
	return c.ipServers
}

func (c *ClientCalculix) updateIPServers() {
	var ipServers []string
	ch := make(chan string)
	quit := make(chan int)

	var wg sync.WaitGroup

	go func() {
		for c := range ch {
			ipServers = append(ipServers, c)
		}
		quit <- 1
	}()

	for i := 1; i <= 255; i++ {
		wg.Add(1)
		serverAddress := fmt.Sprintf("%v%v:1234", c.IPPrototype, i)
		go func() {
			defer wg.Done()
			checkIP(serverAddress, ch)
		}()
	}
	wg.Wait()
	close(ch)
	<-quit
	c.ipServers = ipServers
}

func checkIP(ip string, ch chan<- string) {
	client, err := rpc.DialHTTP("tcp", ip)
	if err != nil {
		return
	}
	//var amount int
	var amount serverCalculix.Amount
	err = client.Call("Calculix.MaxAllowableTasks", "", &amount)
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	if amount.A < 0 {
		fmt.Println("Cannot allowable tasks less zero")
		return
	}
	err = client.Close()
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	ch <- ip
}

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
