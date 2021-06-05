package core

import (
	"log"
	"net"
)

func GetFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "127.0.0.1:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}

	return
}

func CloseProxies(proxies []TorProxy) {
	for _, proxy := range proxies {
		log.Println("Cleanup proxy")
		proxy.Close()
	}
}
