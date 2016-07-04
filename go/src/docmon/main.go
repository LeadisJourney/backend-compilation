package main

import (
	"fmt"
	"log"
	// "time"
	"net/http"
	"encoding/json"

	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/client"
)

type Containers struct {
	Cpu, Mem float64
	Name string
}

func Handlews(ws *websocket.Conn) {
	c := []Containers{{1.1, 1.1, "hello"},{2.2, 2.2, "world"}}
	c2 := []Containers{{2.2, 2.2, "hello"}}
	b, _ := json.Marshal(c)
	b2, _ := json.Marshal(c2)
	fmt.Println(string(b))

	// TEST
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
	}
	options := types.ContainerListOptions{}
	containers, err := cli.ContainerList(context.Background(), options)
	if err != nil {
		fmt.Println(err)
	}
	for _, cont := range containers {
		fmt.Println(cont)
	}
	// END TEST
	
	// for {
	// 	ws.Write(b)
	// 	time.Sleep(3000 * time.Millisecond)
	// 	ws.Write(b2)
	// 	time.Sleep(3000 * time.Millisecond)
	// }
	ws.Write(b)
	_ = b2
}

func main() {
	
	http.Handle("/", http.FileServer(http.Dir(".")))	
	http.Handle("/test", websocket.Handler(Handlews))


	if err := http.ListenAndServe(":2222", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
