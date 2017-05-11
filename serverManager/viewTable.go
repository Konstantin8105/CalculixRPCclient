package serverManager

import (
	"fmt"
	"net/rpc"
	"time"

	"github.com/Konstantin8105/CalculixRPCserver/serverCalculix"
)

// ViewTable - view table of servers
func (s *ServerManager) ViewTable() (result string) {
	result += fmt.Sprintf("------|----------------------|----------------------|----------------------|----------------------|\n")
	result += fmt.Sprintf("%4v. | %20v | %20v | %20v | %20v |\n", "№", "Server address", "Amount of processors", "ServerId", "Allowable ccx")
	result += fmt.Sprintf("------|----------------------|----------------------|----------------------|----------------------|\n")

	ips := s.ipServers

	for i, ip := range ips {

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
		var id serverCalculix.ServerName
		err = client.Call("Calculix.GetServerName", "", &id)
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

		result += fmt.Sprintf("%4v. | %20v | %20v | %20v | %20v |\n", i, ip, fmt.Sprintf("%v/%v", amountTask.A, amount.A), id.A, check.A)
	}
	result += fmt.Sprintf("------|----------------------|----------------------|----------------------|----------------------|\n")
	return result
}

// ViewServerPerformance - view time of connect to server
func (s *ServerManager) ViewServerPerformance() (result string) {

	ips := s.ipServers
	amountTests := 10000

	result += fmt.Sprintf("=====================\n")
	for _, ip := range ips {

		client, err := rpc.DialHTTP("tcp", ip)
		if err != nil {
			return
		}
		//
		result += fmt.Sprintf("Server IP: %v\n", ip)
		start := time.Now()
		for j := 0; j < amountTests; j++ {
			var amount serverCalculix.Amount
			err = client.Call("Calculix.MaxAllowableTasks", "", &amount)
			if err != nil {
				fmt.Println("err = ", err)
				return
			}
		}
		dTime := time.Now().Sub(start).Nanoseconds()
		result += fmt.Sprintf("Average time = %.3v microseconds\n", float64(dTime)/float64(amountTests)/float64(1000))
		result += fmt.Sprintf("Summary time = %.3v milliseconds\n", float64(dTime)/float64(1000*1000))
		result += fmt.Sprintf("=====================\n")
		//
		err = client.Close()
		if err != nil {
			fmt.Println("err = ", err)
			return
		}
	}

	return result
}
