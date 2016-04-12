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
	
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/client"

)

type Container struct {
	ID, UserID, Volume string
	Stream types.HijackedResponse
	UnixSock net.Conn
	Time time.Time
}

type Client struct {
	Pcli *client.Client
	Cont []Container
	c chan int
}

func NewClient(c chan int) (*Client) {
	var cli Client
	var err error
	var build  types.ImageBuildResponse
	
	cli.c = c
	
	cli.Pcli, err = client.NewEnvClient()
	if err != nil {
		Error.Println(err)
		return nil
	}

	cwd, _ := os.Getwd()
	ctx, err := os.Open(cwd+"/dockerfile.tar.gz")
	if err != nil {
		Error.Println(err)
		return nil
	}
	build, err = cli.Pcli.ImageBuild(context.Background(), types.ImageBuildOptions{Tags: []string{"leadis_image"}, Remove: true, Context: ctx, SuppressOutput: false})
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

func (cli *Client) ExecuteProgram(UserID, code, lang, types string) (string) {
	var res string
	
	// fmt.Println("\nChecking if", UserID, "Container Exists...")
	idx := findContainer(UserID, cli.Cont)
	if idx == -1 {
		// fmt.Println("Adding", UserID, "Container...")
		err := cli.AddContainer(UserID)
		if err == false {
			return ""
		}
		
		// fmt.Println("Starting", UserID, "Container...")
		err = cli.StartContainer(UserID)
		if err == false {
			return ""
		}
		idx = findContainer(UserID, cli.Cont)
	}
	
	cli.Cont[idx].Time = time.Now()

	// fmt.Println("Copying Code to", UserID, "Container...")
	err := cli.CopytoContainer(idx, code, lang)
	if err == false {
		return ""
	}

	// fmt.Println("Compiling", UserID, "Code...")
	err = cli.CompileRequest(idx, lang, types)
	if err == false {
		return ""
	}

	// fmt.Println("Getting", UserID, "Response...\n")
	res, err = cli.GetResponse(idx)
	if err == false {
		return ""
	}
	cli.c <- 1
	return res
}


// Check error value
func (cli *Client) GetResponse(idx int) (string, bool) {
	tmp := make([]byte, 1)
	cli.Cont[idx].UnixSock.Read(tmp)
	res, _ := ioutil.ReadFile(cli.Cont[idx].Volume+"/stdout")
	return string(res), true
}

func (cli *Client) CompileRequest(idx int, Lang, Req string) (bool) {
	var buf bytes.Buffer
	
	buf.WriteString(strings.ToUpper(Lang)+" "+strings.ToUpper(Req)+" "+cli.Cont[idx].UserID)
	cli.Cont[idx].UnixSock.Write(buf.Bytes())
	return true
}

func (cli *Client) DeleteContainer(idx int) (bool) {
	err := cli.Pcli.ContainerStop(context.Background(), cli.Cont[idx].ID, 0)
	if err != nil {
		Error.Println(err)
		return false
	}
	
	cli.Cont[idx] = cli.Cont[len(cli.Cont)-1]
	cli.Cont = cli.Cont[:len(cli.Cont)-1]
	
	return true	
}

func (cli *Client) CopytoContainer(idx int, code, lang string) (bool) {
	f, err := os.Create(cli.Cont[idx].Volume+"/main."+strings.ToLower(lang))
	if err != nil {
		Error.Println(err)
		return false
	}
	_, err = f.Write([]byte(code))
	if err != nil {
		Error.Println(err)
		return false
	}
	return true
}

func (cli *Client) StartContainer(UserID string) (bool) {
	idx := findContainer(UserID, cli.Cont)
	if idx == -1 {
		return false
	}
	l, err := net.Listen("unix", cli.Cont[idx].Volume+"/host.sock")
	if err != nil {
		Error.Println(err)
		return false
	}
	err = cli.Pcli.ContainerStart(context.Background(), cli.Cont[idx].ID)
	if err != nil {
		Error.Println(err)
		return false
	}

	cli.Cont[idx].UnixSock, err = l.Accept()
	if err != nil {
		Error.Println(err)
		return false
	}	
	return true
}

func (cli *Client) AddContainer(UserID string) (bool) {
	var cont Container
	var vol types.ContainerJSON
	
	cont.UserID = UserID
	resp, err := cli.Pcli.ContainerCreate(context.Background(), initConfig(), nil, nil, "")
	if err != nil {
		Error.Println(err)
		return false
	}
	cont.ID = resp.ID
	cont.Stream, err = cli.Pcli.ContainerAttach(context.Background(), types.ContainerAttachOptions{ContainerID: cont.ID, Stream: true, Stdout: true, Stderr: true})
	if err != nil {
		Error.Println(err)
		return false
	}
	vol, err = cli.Pcli.ContainerInspect(context.Background(), cont.ID)
	if err != nil {
		Error.Println(err)
		return false
	}
	cont.Volume = vol.Mounts[0].Source	
	cli.Cont = append(cli.Cont, cont)
	return true
}
