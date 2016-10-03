package main

import (
	"fmt"
	"net"
	"log"
	// "encoding/binary"
)

const (
	LIST  =iota
	DELETE  = iota
	t2  = iota
	t3  = iota
	t4  = iota
	t5  = iota
	t6  = iota
	t7  = iota
)

func main() {
	// buf := make([]byte, 8)
	l, err := net.Listen("unix", "/tmp/test12.sock")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	for {
		fmt.Println([]byte(t2))
		// binary.BigEndian.PutUint64(buf, t6)
		// n := binary.BigEndian.Uint64(buf)
		// fmt.Println(buf, n)		
		// binary.BigEndian.PutUint64(buf, t2)
		// n = binary.BigEndian.Uint64(buf)
		// fmt.Println(buf, n)
		// binary.BigEndian.PutUint64(buf, t3)
		// n = binary.BigEndian.Uint64(buf)
		// fmt.Println(buf, n)
		// binary.BigEndian.PutUint64(buf, t4)
		// n = binary.BigEndian.Uint64(buf)
		// fmt.Println(buf, n)
		// binary.BigEndian.PutUint64(buf, t5)
		// n = binary.BigEndian.Uint64(buf)
		// fmt.Println(buf, n)
		// binary.BigEndian.PutUint64(buf, t6)
		// n = binary.BigEndian.Uint64(buf)
		// fmt.Println(buf, n)

		fd, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		_ = fd
		
	}
}
