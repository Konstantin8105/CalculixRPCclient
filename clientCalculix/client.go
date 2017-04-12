package clientCalculix

import (
	"fmt"
	"net/rpc"

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
	fmt.Println("updateIPServers")
	var ipServers []string
	ch := make(chan string)
	quit := make(chan int)

	go func() {
		fmt.Println("In channel")
		for c := range ch {
			fmt.Println("channel.c = ", c)
			ipServers = append(ipServers, c)
		}
		fmt.Println("quit <- 1")
		quit <- 1
	}()

	for i := 8; i <= 10; i++ {
		fmt.Println("i = ", i)
		go func(port int) {
			//fmt.Println("============")
			serverAddress := fmt.Sprintf("%v%v", c.IPPrototype, port)
			fmt.Println("serverAddress = ", serverAddress)
			client, err := rpc.DialHTTP("tcp", serverAddress+":1234")
			if err != nil {
				fmt.Println("err = ", err)
				//continue
				return
			}
			///calculix := serverCalculix.NewCalculix()
			fmt.Println("Create calculix")
			//var amount int
			var amount serverCalculix.Amount
			err = client.Call("Calculix.MaxAllowableTasks", "", &amount)
			if err != nil {
				fmt.Println("err = ", err)
				//continue
				return
			}
			fmt.Println("amount task on server = ", amount)
			if amount.A < 0 {
				fmt.Println("err = ", err)
				//continue
				return
			}
			/*
					var ttt string
					err = client.Call("Calculix.AmountTasks", "empty", &ttt) //&amount)
					fmt.Println(err, ttt)
					if err != nil {
						fmt.Println("err = ", err)
						//continue
						return
					}
				fmt.Println("amount task on server = ", amount)
				if amount < 0 {
					fmt.Println("err = ", err)
					//continue
					return
				}
			*/
			fmt.Println("close client")
			err = client.Close()
			if err != nil {
				fmt.Println("err = ", err)
				//continue
				return
			}
			fmt.Println("Add ip server")
			//ipServers = append(ipServers, serverAddress)
			fmt.Println("============")
			ch <- serverAddress
		}(i)
	}
	fmt.Println("<-quit")
	<-quit
	fmt.Println("close(ch)")
	close(ch)
	fmt.Println("============")
	fmt.Println("List of server ip: ", ipServers)
	fmt.Println("============")
	c.ipServers = ipServers
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
