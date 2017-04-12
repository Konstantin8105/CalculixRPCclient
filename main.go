package main

import (
	"fmt"

	"github.com/Konstantin8105/CalculixRPCclient/clientCalculix"
)

func main() {
	client := clientCalculix.ClientCalculix{IPPrototype: "192.168.0."}
	fmt.Println("Servers:\n", client.GetIPServers())
}
