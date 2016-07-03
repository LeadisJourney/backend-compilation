package main

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"strings"
	"io/ioutil"
	"encoding/json"
	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
)

type Stats struct {
	
}

// CPU TEST

func getCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val
				if i == 4 {
					idle = val
				}
			}
			return
		}
	}
	return
}


// END CPU TEST


func calculateCPUPercent(c *Container, cpu types.CPUStats) float64 {
	cpuPercent := 0.0
	cpuDelta := float64(cpu.CPUUsage.TotalUsage) - float64(c.PCPU)
	systemDelta := float64(cpu.SystemUsage) - float64(c.PSys)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(cpu.CPUUsage.PercpuUsage)) * 100.0
	}
	c.PCPU = cpu.CPUUsage.TotalUsage
	c.PSys = cpu.SystemUsage
	return cpuPercent
}

func (cli *Client) ReadStats() {
	var stats types.Stats

	f, _ := os.Create("/home/exploit/backend-compilation/go/src/docman/stats")
	defer f.Close()
	
	for {
		//os.Truncate("/home/exploit/backend-compilation/go/src/docman/stats", 0)

		idle0, total0 := getCPUSample()
		
		for _, value := range cli.Cont {
			r, _ := cli.Pcli.ContainerStats(context.Background(), value.ID, false)
			b, _ := ioutil.ReadAll(r)
			json.Unmarshal(b, &stats)

			f.WriteString("USERID: "+value.UserID+"\n")

			// TODO Mutlitple core
			// for i := 0; i < len(stats.CPUStats.CPUUsage.PercpuUsage); i++ {
			// 	f.WriteString("CPU: "+strconv.FormatUint(stats.CPUStats.CPUUsage.PercpuUsage[i], 10)+"\n")
			// }
			
			//f.WriteString("MEMORY USAGE: "+strconv.FormatUint(stats.MemoryStats.Usage, 10)+"\n")

			f.WriteString("CPU: "+strconv.FormatFloat(calculateCPUPercent(value, stats.CPUStats), 'f', 3, 64)+"\n")			
		}

		idle1, total1 := getCPUSample()
		idleTicks := float64(idle1 - idle0)
		totalTicks := float64(total1 - total0)
		cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks
		f.WriteString("HOST CPU: "+strconv.FormatFloat(cpuUsage,'f', 3, 64)+"\n")
		
		f.WriteString("\n\n")
		time.Sleep(1000 * time.Millisecond)
	}	
}
