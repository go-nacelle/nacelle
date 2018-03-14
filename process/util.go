package process

import (
	"fmt"
	"net"
)

func makeListener(port int) (*net.TCPListener, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return nil, err
	}

	return net.ListenTCP("tcp", addr)
}
