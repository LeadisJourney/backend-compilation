package main

import (
	"os"
	"fmt"
	"net"
	"time"
	"bytes"
	"strings"
	"io/ioutil"
	"golang.org/x/net/context"
	
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const (
	RTime time.Duration = 60 * time.Second
)

type Container struct {
	ID, UserID, Volume string
	Stream types.HijackedResponse
	UnixSock net.Conn
	Time time.Time
	PCPU uint64
	PSys uint64
}


type Client struct {
	Pcli *client.Client
	Cont map[string]*Container
	c chan int
}

func NewClient(c chan int) (*Client) {
	var cli Client
	var err error
	var build  types.ImageBuildResponse
	
	cli.c = c
	cli.Cont = make(map[string]*Container)

	cli.Pcli, err = client.NewEnvClient()
	if err != nil {
		Error.Println(err)
		return nil
	}

	cwd, _ := os.Getwd()
	ctx, err := os.Open(cwd+"/Dockerfile.tar.gz")
	if err != nil {
		Error.Println(err)
		return nil
	}
	build, err = cli.Pcli.ImageBuild(context.Background(), ctx, types.ImageBuildOptions{Tags: []string{"leadis_image"}, Remove: true, Context: ctx, SuppressOutput: false})
	if err != nil {
		Error.Println(err)
		return nil
	}

	// Test
	b, _ := ioutil.ReadAll(build.Body)
	fmt.Println(string(b))
	// End Test
	
	return &cli
}

func (cli *Client) ExecuteProgram(UserID, code, lang, types, ex string) (string, error) {
	var res string
	
	// fmt.Println("\nChecking if", UserID, "Container Exists...")
	if _, ok := cli.Cont[UserID]; !ok {
		// TEST
		// if len(cli.Cont) >= 1 {
		// 	cli.OldestContainer()
		// }
		// END TEST
		// fmt.Println("Adding", UserID, "Container...")
		err := cli.AddContainer(UserID)
		if err != nil {
			Error.Println(err)
			return "", err
		}
		
		// fmt.Println("Starting", UserID, "Container...")
		err = cli.StartContainer(UserID)
		if err != nil {
			Error.Println(err)
			return "", err
		}
	}

	cli.Cont[UserID].Time = time.Now()

	// fmt.Println("Copying Code to", UserID, "Container...")
	err := cli.CopytoContainer(UserID, code, lang, ex)
	// fmt.Println("\nCODE BEGIN\n", code, "\nCODE END\n\n")
	if err != nil {
		Error.Println(err)
		return "", err
	}

	// fmt.Println("Compiling", UserID, "Code...")
	err = cli.CompileRequest(UserID, lang, types)
	if err != nil {
		Error.Println(err)
		return "", err
	}

	// fmt.Println("Getting", UserID, "Response...\n")
	res, err = cli.GetResponse(UserID)
	if err != nil {
		Error.Println(err)
		return "", err
	}
	cli.c <- 1

	
	// BEGIN STATS TEST
	//cli.ReadStats(UserID)
	// END STATS TEST
	return res, nil
}

// Check error value
func (cli *Client) GetResponse(UserID string) (string, error) {
	tmp := make([]byte, 4)

	t := time.Now()
	t = t.Add(RTime)
	err := cli.Cont[UserID].UnixSock.SetReadDeadline(t)
	if err != nil {
		return "", err
	}	
	cli.Cont[UserID].UnixSock.Read(tmp)
	
	// TEST
	// fmt.Println(tmp)
	// END TEST
	
	res, _ := ioutil.ReadFile(cli.Cont[UserID].Volume+"/stdout")
	return string(res), nil
}

func (cli *Client) CompileRequest(UserID, Lang, Req string) (error) {
	var buf bytes.Buffer

	fmt.Println(strings.ToUpper(Lang)+" "+strings.ToUpper(Req)+" "+cli.Cont[UserID].UserID)
	
	buf.WriteString(strings.ToUpper(Lang)+" "+strings.ToUpper(Req)+" "+cli.Cont[UserID].UserID)
	_, err := cli.Cont[UserID].UnixSock.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}


func (cli *Client) OldestContainer() {
	old := ""
	if len(cli.Cont) > 0 {
		for key, _ := range cli.Cont {
			if old == "" || time.Since(cli.Cont[old].Time) < time.Since(cli.Cont[key].Time) {
				old = key
			}
		}
	}
	cli.DeleteContainer(old)
}

func (cli *Client) DeleteContainer(UserID string) (error) {
	timeout := time.Second * 0
	err := cli.Pcli.ContainerStop(context.Background(), cli.Cont[UserID].ID, &timeout)
	if err != nil {
		return err
	}
	delete(cli.Cont, UserID)
	return nil	
}

func (cli *Client) CopytoContainer(UserID, code, lang, ex string) (error) {
	//fmt.Println("Exersice name", ex)
	
	os.Remove(cli.Cont[UserID].Volume+"/exercise")
	CopyDir(ex, cli.Cont[UserID].Volume+"/exercise")
	f, err := os.Create(cli.Cont[UserID].Volume+"/exercise/src/"+ex+"."+strings.ToLower(lang))
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(code))
	if err != nil {
		return err
	}
	return nil
}

func (cli *Client) StartContainer(UserID string) (error) {
	l, err := net.Listen("unix", cli.Cont[UserID].Volume+"/host.sock")
	if err != nil {
		return err
	}
	err = cli.Pcli.ContainerStart(context.Background(), cli.Cont[UserID].ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	cli.Cont[UserID].UnixSock, err = l.Accept()
	if err != nil {
		return err
	}
	return nil
}

func (cli *Client) AddContainer(UserID string) (error) {
	var vol types.ContainerJSON
	var cont Container
	
	cont.UserID = UserID
	// REMOVE
	cont.PCPU = 0
	cont.PSys = 0
	// END REMOVE
	resp, err := cli.Pcli.ContainerCreate(context.Background(), initConfig(), nil, nil, "leadis journey")
	if err != nil {
		return err
	}
	cont.ID = resp.ID
	cont.Stream, err = cli.Pcli.ContainerAttach(context.Background(), cont.ID, types.ContainerAttachOptions{Stream: true, Stdout: true, Stderr: true})
	if err != nil {
		return err
	}
	vol, err = cli.Pcli.ContainerInspect(context.Background(), cont.ID)
	if err != nil {
		return err
	}

	cont.Volume = vol.Mounts[0].Source
	fmt.Println(cont.Volume)
	cli.Cont[UserID] = &cont
	Info.Println("Creating New Container, ID: ", UserID)
	return nil
}
