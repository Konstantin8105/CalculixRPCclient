package serverManager

import (
	"bufio"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Konstantin8105/CalculixRPCserver/serverCalculix"
)

// UpdateServers - updating list of servers
// by default - check only saved server address in file
func (s *ServerManager) UpdateServers() {
	serverList, err := s.openFile()
	if err == nil {
		if len(serverList) == 0 {
			// remove file
			err = s.removeFile()
			if err != nil {
				fmt.Println("Cannot remove file:", serverManagerFile, "\nErr = ", err)
				os.Exit(0)
			}
			// recursive try update server list
			s.UpdateServers()
			return
		}
		// check all ip from file
		for _, address := range serverList {
			if !s.checkIP(address) {
				// remove file
				err = s.removeFile()
				if err != nil {
					fmt.Println("Cannot remove file:", serverManagerFile, "\nErr = ", err)
					os.Exit(0)
				}
				s.UpdateServers()
				return
			}
		}
		// refresh server list
		s.setIPServers(serverList)
		return
	}
	// file is not exist
	s.updateServersFromNet()
}

func (s *ServerManager) updateServersFromNet() {
	adresess, err := s.generateServersIP()
	if err != nil {
		fmt.Println("You haven`t local network: ", err)
		os.Exit(0)
	}
	serverList := s.foundIPServers(adresess)
	if len(serverList) == 0 {
		// wait and try again
		timer := time.NewTimer(time.Second * 10)
		<-timer.C
		s.UpdateServers()
		return
	}
	err = s.saveFile(serverList)
	if err != nil {
		fmt.Println("Error to save: ", err)
	}
	s.setIPServers(serverList)
}

func (s *ServerManager) setIPServers(serverList []string) {
	var mutex = &sync.Mutex{}
	mutex.Lock()

	s.ipServers = serverList

	mutex.Unlock()
}

func (s *ServerManager) removeFile() (err error) {
	if _, err := os.Stat(serverManagerFile); os.IsExist(err) {
		return os.Remove(serverManagerFile)
	}
	return nil
}

func (s *ServerManager) openFile() (serverList []string, err error) {
	// check file is exist
	if _, err := os.Stat(serverManagerFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("Cannot find %v file : %v", serverManagerFile, err)
	}
	// open file
	inFile, err := os.Open(serverManagerFile)
	if err != nil {
		return nil, err
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
		serverList = append(serverList, scanner.Text())
	}
	serverCalculix.RemoveDuplicates(&serverList)
	return serverList, nil
}

func (s *ServerManager) saveFile(serverList []string) (err error) {
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
	for _, line := range serverList {
		_, err = fmt.Fprintf(f, "%v\n", line)
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

func (s *ServerManager) foundIPServers(addresses []string) (serverList []string) {

	ch := make(chan string)
	quit := make(chan int)

	var wg sync.WaitGroup

	go func() {
		for c := range ch {
			serverList = append(serverList, c)
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

	return s.removeDublicateName(serverList)
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

func (s *ServerManager) removeDublicateName(addresses []string) (serverList []string) {
	type ipWithName struct {
		ip   string
		name string
	}

	var ipName []ipWithName

	ch := make(chan ipWithName)
	quit := make(chan int)

	var wg sync.WaitGroup

	go func() {
		for c := range ch {
			ipName = append(ipName, c)
		}
		quit <- 1
	}()

	for _, address := range addresses {
		wg.Add(1)
		go func(address string) {
			defer wg.Done()
			client, err := rpc.DialHTTP("tcp", address)
			if err != nil {
				return
			}
			//var amount int
			var id serverCalculix.ServerName
			err = client.Call("Calculix.GetServerName", "", &id)
			if err != nil {
				fmt.Println("err = ", err)
				return
			}
			ch <- ipWithName{ip: address, name: id.A}
		}(address)
	}
	wg.Wait()
	close(ch)
	<-quit

checkAgain:
	for i := range ipName {
		for j := range ipName {
			if i != j {
				if ipName[i].name == ipName[j].name {
					ipName = append(ipName[:i], ipName[i+1:]...)
					goto checkAgain
				}
			}
		}
	}

	for i := range ipName {
		serverList = append(serverList, ipName[i].ip)
	}

	return serverList
}
