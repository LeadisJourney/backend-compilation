package main

import (
	"fmt"
	"log"
	"time"
	"io/ioutil"
	"net/http"
	"encoding/json"

	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/client"
)

type Containers struct {
	Cpu float64
	Mem uint64
	Name string
	Action int
}

func calculateCPUPercent(pcpu, cpu types.CPUStats) float64 {
	cpuPercent := 0.0
	cpuDelta := float64(cpu.CPUUsage.TotalUsage) - float64(pcpu.CPUUsage.TotalUsage)
	systemDelta := float64(cpu.SystemUsage) - float64(pcpu.SystemUsage)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(cpu.CPUUsage.PercpuUsage)) * 100.0
	}
	return cpuPercent
}

func Handlews(ws *websocket.Conn) {
	var stats types.Stats
	var pre []Containers
	// c1 := []Containers{{1.1, 1.1, "hello"},{2.2, 2.2, "world"}}
	// c2 := []Containers{{2.2, 2.2, "hello"}}
	// b, _ := json.Marshal(c1)
	// b2, _ := json.Marshal(c2)
	// fmt.Println(string(b))

	// TEST
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
	}
	options := types.ContainerListOptions{}

	for {
		c := []Containers{}
		containers, err := cli.ContainerList(context.Background(), options)
		if err != nil {
			fmt.Println(err)
		}
		for _, cont := range containers {
			r, _ := cli.ContainerStats(context.Background(), cont.ID, false)
			b, _ := ioutil.ReadAll(r)
			json.Unmarshal(b, &stats)
			tmp := Containers{calculateCPUPercent(stats.PreCPUStats, stats.CPUStats), stats.MemoryStats.Usage, cont.ID}
			c = append(c, tmp)
		}
		fmt.Println(c)
		b, _ := json.Marshal(c)
		ws.Write(b)		
		time.Sleep(1000 * time.Millisecond)
	}
	// END TEST
	
	// for {
	// 	ws.Write(b)
	// 	time.Sleep(3000 * time.Millisecond)
	// 	ws.Write(b2)
	// 	time.Sleep(3000 * time.Millisecond)
	// }
	// ws.Write(b)
	// _ = b2
}

func main() {
	
	http.Handle("/", http.FileServer(http.Dir(".")))	
	http.Handle("/test", websocket.Handler(Handlews))


	if err := http.ListenAndServe(":2222", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
