package serverManager

import (
	"fmt"
	"net/rpc"

	"github.com/Konstantin8105/CalculixRPCserver/serverCalculix"
)

// ViewTable - view table of servers
func (s *ServerManager) ViewTable() (result string) {
	result += fmt.Sprintf("------|----------------------|----------------------|----------------------|\n")
	result += fmt.Sprintf("%4v. | %20v | %20v | %20v |\n", "№", "Server address", "Amount of processors", "Allowable ccx")
	result += fmt.Sprintf("------|----------------------|----------------------|----------------------|\n")
	for i, ip := range s.ipServers {

		client, err := rpc.DialHTTP("tcp", ip)
		if err != nil {
			return
		}
		//
		var amount serverCalculix.Amount
		err = client.Call("Calculix.MaxAllowableTasks", "", &amount)
		if err != nil {
			fmt.Println("err = ", err)
			return
		}
		//
		var amountTask serverCalculix.Amount
		err = client.Call("Calculix.AmountTasks", "", &amountTask)
		if err != nil {
			fmt.Println("err = ", err)
			return
		}
		//
		var check serverCalculix.ChechCCXResult
		err = client.Call("Calculix.CheckCCX", "", &check)
		if err != nil {
			fmt.Println("err = ", err)
			return
		}
		//
		err = client.Close()
		if err != nil {
			fmt.Println("err = ", err)
			return
		}

		result += fmt.Sprintf("%4v. | %20v | %20v | %20v |\n", i, ip, fmt.Sprintf("%v/%v", amountTask.A, amount.A), check.A)
	}
	result += fmt.Sprintf("------|----------------------|----------------------|----------------------|\n")
	return result
}
