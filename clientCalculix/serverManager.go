package clientCalculix

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"sync"

	"github.com/Konstantin8105/CalculixRPCserver/serverCalculix"
)

const (
	serverManagerFile string = "serverIpManager.txt"
)

// ServerManager - manager of servers
type ServerManager struct {
	ipServers   []string
	IPPrototype string
}

// NewServerManager - create a new
func NewServerManager(ipPrototype string) (s *ServerManager) {
	s = new(ServerManager)
	s.IPPrototype = ipPrototype
	err := s.openFile()
	if err == nil {
		return
	}
	s.updateIPServers()
	//TODO: add strategy for case - if len of ip list is zero
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
	result += fmt.Sprintf("%4v.\t%30v\t%20v\n", "№", "Server address", "Amount of processors")
	for i, ip := range s.ipServers {

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
		err = client.Close()
		if err != nil {
			fmt.Println("err = ", err)
			return
		}

		result += fmt.Sprintf("%4v.\t%30v\t%20v\n", i, ip, amount.A)
	}
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

func (s *ServerManager) updateIPServers() {
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
		serverAddress := fmt.Sprintf("%v%v:1234", s.IPPrototype, i)
		go func() {
			defer wg.Done()
			s.checkIP(serverAddress, ch)
		}()
	}
	wg.Wait()
	close(ch)
	<-quit
	s.ipServers = ipServers
}

func (s *ServerManager) checkIP(ip string, ch chan<- string) {
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
