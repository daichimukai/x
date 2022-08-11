package testutils

import (
	"net"
	"strings"
	"time"
)

func GetFreePorts(count int) []string {
	var ports []string
	for count > 0 {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		defer l.Close()
		ports = append(ports, strings.ReplaceAll(l.Addr().String(), "127.0.0.1:", ""))
		count--
	}
	return ports
}

func GetFreePort() string {
	return GetFreePorts(1)[0]
}
