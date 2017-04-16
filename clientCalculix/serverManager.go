package clientCalculix

import (
	"bufio"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strings"
	"sync"

	"github.com/Konstantin8105/CalculixRPCserver/serverCalculix"
)

const (
	serverManagerFile string = "serverIpManager.txt"
)

// ServerManager - manager of servers
type ServerManager struct {
	ipServers []string // IP allowable servers
}

// NewServerManager - create a new
func NewServerManager() (s *ServerManager) {
	s = new(ServerManager)
	err := s.openFile()
	if err == nil {
		return
	}
	adresess, err := s.generateServersIP()
	if err != nil {
		fmt.Println("You haven`t local network")
		os.Exit(0)
	}
	s.updateIPServers(adresess)
	if len(s.ipServers) == 0 {
		fmt.Println("Cannot search any servers")
		os.Exit(0)
	}
	err = s.saveFile()
	if err != nil {
		fmt.Println("Error to save : ", err)
	}
	return s
}

// GetIPServers - get all ip servers
func (s *ServerManager) GetIPServers() []string {
	return s.ipServers
}

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

		result += fmt.Sprintf("%4v. | %20v | %20v | %20v |\n", i, ip, amount.A, check.A)
	}
	result += fmt.Sprintf("------|----------------------|----------------------|----------------------|\n")
	return result
}

func (s *ServerManager) openFile() (err error) {
	// check file is exist
	if _, err := os.Stat(serverManagerFile); os.IsNotExist(err) {
		return fmt.Errorf("Cannot find %v file : %v", serverManagerFile, err)
	}
	// open file
	inFile, err := os.Open(serverManagerFile)
	if err != nil {
		return err
	}
	defer func() {
		errFile := inFile.Close()
		if errFile != nil {
			if err != nil {
				err = fmt.Errorf("%v ; %v", err, errFile)
			} else {
				err = errFile
			}
		}
	}()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		s.ipServers = append(s.ipServers, scanner.Text())
	}
	serverCalculix.RemoveDuplicates(&s.ipServers)
	return nil
}

func (s *ServerManager) saveFile() (err error) {
	// check file is exist
	if _, err := os.Stat(serverManagerFile); os.IsNotExist(err) {
		// create file
		newFile, err := os.Create(serverManagerFile)
		if err != nil {
			return err
		}
		err = newFile.Close()
		if err != nil {
			return err
		}
	}

	// open file
	f, err := os.OpenFile(serverManagerFile, os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	for _, line := range s.ipServers {
		_, err = fmt.Fprintln(f, line)
		if err != nil {
			return fmt.Errorf("Cannot write to file: %v", err)
		}
	}
	return nil
}

func (s *ServerManager) generateServersIP() (addresses []string, err error) {

	port := ":1234"

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return addresses, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ad := strings.Split(ipnet.IP.String(), ".")
				for i := 1; i <= 255; i++ {
					addresses = append(addresses, fmt.Sprintf("%v.%v.%v.%v%v", ad[0], ad[1], ad[2], i, port))
				}
			}
		}
	}

	serverCalculix.RemoveDuplicates(&addresses)

	return addresses, nil
}

func (s *ServerManager) updateIPServers(addresses []string) {
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

	for _, address := range addresses {
		wg.Add(1)
		go func(address string) {
			defer wg.Done()
			if s.checkIP(address) {
				ch <- address
			}
		}(address)
	}
	wg.Wait()
	close(ch)
	<-quit
	s.ipServers = ipServers
}

func (s *ServerManager) checkIP(ip string) bool {
	client, err := rpc.DialHTTP("tcp", ip)
	if err != nil {
		return false
	}
	//var amount int
	var amount serverCalculix.Amount
	err = client.Call("Calculix.MaxAllowableTasks", "", &amount)
	if err != nil {
		fmt.Println("err = ", err)
		return false
	}
	if amount.A < 0 {
		fmt.Println("Cannot allowable tasks less zero")
		return false
	}
	err = client.Close()
	if err != nil {
		fmt.Println("err = ", err)
		return false
	}
	return true
}
