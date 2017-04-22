package serverManager

import (
	"time"
)

const (
	serverManagerFile string = "serverIpManager.txt"
)

// ServerManager - manager of servers
type ServerManager struct {
	ipServers []string // IP allowable servers
}

// NewServerManager - constructor for creating ServerManager
func NewServerManager() (s *ServerManager) {
	s = new(ServerManager)
	s.UpdateServers()

	//function for autoupdate server list in file from network
	ticker := time.NewTicker(time.Minute * 30)
	go func() {
		for _ = range ticker.C {
			s.updateServersFromNet()
		}
	}()

	return s
}

// GetIPServers - get all ip servers
func (s *ServerManager) GetIPServers() []string {
	return s.ipServers
}
