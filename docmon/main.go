package main

import (
	"net"
	"fmt"
	"log"
	"time"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"encoding/binary"

	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const (
	LIST = iota
	DELETE =iota
)

type Containers struct {
	Cpu float64
	Mem uint64
	Name string
	Action int
}

type MonMan struct {
	Unix net.Conn
	User, ID []string
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

func NewMonMan() (m *MonMan) {
	c, err := net.Dial("unix", "/tmp/monitor.sock")
	if err != nil {
		fmt.Println(err)
		return
	}	
	return &MonMan{Unix: c}
}

func (m *MonMan)GetContainer() {
	bl := make([]byte, 8)
	
	b := []byte{LIST}
	fmt.Println("SENDING REQUEST")
	m.Unix.Write(b)
	fmt.Println("SENT REQUEST")
	_, err := m.Unix.Read(bl)
	if err != nil {
		fmt.Println(err)
	}
	l := binary.BigEndian.Uint64(bl)
	bc := make([]byte, l)
	_, err = m.Unix.Read(bc)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(bc)
	json.Unmarshal(bc, m)
	// fmt.Println(m)
}

func UpdateList(cold, cnew []Containers) ([]Containers) {
	for k := range cold {
		cold[k].Action = 1
		for k2 := range cnew {
			if cnew[k2].Name == cold[k].Name {
				cold[k].Action = 0
				fmt.Println("Keep container")
				break
			}
		}
		if cold[k].Action == 1 {
			fmt.Println("Remove container")
		}
	}
	fmt.Println("Cold List: ", cold)
	return cnew
}

func Handlews(ws *websocket.Conn) {
	var stats types.Stats
	var oc []Containers	
	// TEST
	// m := NewMonMan()
	
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
	}
	options := types.ContainerListOptions{}

	for {
		fmt.Println("maine loop")
		c := []Containers{}
		containers, err := cli.ContainerList(context.Background(), options)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Containers: :", containers)
		for _, cont := range containers {
			fmt.Println("container loop")
			r, _ := cli.ContainerStats(context.Background(), cont.ID, false)
			b, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(b, &stats)
			tmp := Containers{calculateCPUPercent(stats.PreCPUStats, stats.CPUStats), stats.MemoryStats.Usage, cont.ID, 0}
			c = append(c, tmp)
		}
		fmt.Println("New List: ", c)
		fmt.Println("Old List: ", oc)
		
		tmp2 := c
		UpdateList(oc, c)
		oc = tmp2
		b, _ := json.Marshal(c)
		ws.Write(b)

		// m.GetContainer()
		// NewMonMan()
		
		time.Sleep(1000 * time.Millisecond)
	}
	// END TEST
}


func Handlews2(ws *websocket.Conn) {
	var stats types.Stats
	m := NewMonMan()
	var c []Containers
	var oc []Containers


	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
	}
	for {
		fmt.Println("Elements copied: ", copy(oc, c))
		m.GetContainer()
		for n, cont := range m.ID {
			r, _ := cli.ContainerStats(context.Background(), cont, false)
			b, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(b, &stats)
			tmp := Containers{calculateCPUPercent(stats.PreCPUStats, stats.CPUStats), stats.MemoryStats.Usage, m.User[n], 0}
			c = append(c, tmp)
		}
		fc := UpdateList(oc, c)
		fmt.Println("New List: ", c)
		fmt.Println("Old List: ", oc)
		fmt.Println("Final List: ", fc)
		b, _ := json.Marshal(fc)
		ws.Write(b)
		time.Sleep(1000 * time.Millisecond)
	}
}


func main() {
	
	http.Handle("/", http.FileServer(http.Dir(".")))	
	http.Handle("/test", websocket.Handler(Handlews2))


	if err := http.ListenAndServe(":2222", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
