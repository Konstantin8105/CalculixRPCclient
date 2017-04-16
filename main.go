package main

import (
	"fmt"

	"github.com/Konstantin8105/CalculixRPCclient/clientCalculix"
)

func main() {
	manager := clientCalculix.NewServerManager()
	fmt.Println(manager.ViewTable())
}
