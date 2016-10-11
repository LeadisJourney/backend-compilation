package main

import (
	"os"
	"net"
)

func main() {
	var req [512]byte
	
	conn, err := net.Dial("unix", "/root/host/host.sock")
	if err != nil {
		return
	}
	for {
		_, err := conn.Read(req)
		f, _ := os.Create("/root/host/test")
		f.Write(req)
		
	}
}
