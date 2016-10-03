package main

import (
	"fmt"
	"net"
	"encoding/json"
	"encoding/binary"
)

const (
	LIST = iota
	DELETE = iota
)

type MonCont struct {
	User, ID []string
}

func getAllContainers(cli *Client, u net.Conn) {
	var uc MonCont
	bl := make([]byte, 8)
	
	for _, c := range cli.Cont {
		uc.User = append(uc.User, c.UserID)
		uc.ID = append(uc.ID, c.ID)
	}
	b, _ := json.Marshal(uc)
	fmt.Println("SENDING: ", string(b))
	l := uint64(len(b))
	binary.BigEndian.PutUint64(bl, l)
	u.Write(append(bl, b...))
}

func ManageRequest(cli *Client, u net.Conn) {
	b := make([]byte, 1)
	for {
		u.Read(b)
		fmt.Println("RECEIVED: ", b)
		if b[0] == LIST {
			getAllContainers(cli, u)
		}
	}
	u.Close()
}

func MonListen(cli *Client) {
	l, err := net.Listen("unix", "/tmp/monitor.sock")
	if err != nil {
		Error.Println(err)
	}
	for {
		fd, err := l.Accept()
		if err != nil {
			Error.Println(err)
		}
		go ManageRequest(cli, fd)
	}	
}
