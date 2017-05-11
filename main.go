package main

import (
	"fmt"

	"github.com/Konstantin8105/CalculixRPCclient/serverManager"
)

func main() {
	manager := serverManager.NewServerManager()
	fmt.Println(manager.ViewTable())
	fmt.Println(manager.ViewServerPerformance())
}
