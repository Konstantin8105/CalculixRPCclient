package main

import (
	"fmt"

	"github.com/Konstantin8105/CalculixRPCclient/clientCalculix"
)

func main() {
	manager := clientCalculix.NewServerManager("192.168.0.")
	fmt.Println(manager.ViewTable())

	/*


		client := clientCalculix.ClientCalculix{IPPrototype: "192.168.0."}
		//client := clientCalculix.ClientCalculix{IPPrototype: "192.168.5."}
		ips := client.GetIPServers()
		fmt.Println("Servers:\n", ips)

		for _, ip := range ips {
			client, err := rpc.DialHTTP("tcp", ip)
			if err != nil {
				return
			}
			var ccx serverCalculix.ChechCCXResult
			err = client.Call("Calculix.CheckCCX", "", &ccx)
			if err != nil {
				fmt.Println("err = ", err)
				return
			}
			fmt.Println("ccx = ", ccx)
		}

	*/
}
