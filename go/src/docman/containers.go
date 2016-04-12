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
	if _, ok := cli.Cont[UserID]; !ok {
		// TEST
		// if len(cli.Cont) >= 1 {
		// 	cli.OldestContainer()
		// }
		// END TEST
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
	}

	cli.Cont[UserID].Time = time.Now()

	// fmt.Println("Copying Code to", UserID, "Container...")
	err := cli.CopytoContainer(UserID, code, lang)
	if err == false {
		return ""
	}

	// fmt.Println("Compiling", UserID, "Code...")
	err = cli.CompileRequest(UserID, lang, types)
	if err == false {
		return ""
	}

	// fmt.Println("Getting", UserID, "Response...\n")
	res, err = cli.GetResponse(UserID)
	if err == false {
		return ""
	}
	cli.c <- 1
	return res
}

// Check error value
func (cli *Client) GetResponse(UserID string) (string, bool) {
	tmp := make([]byte, 1)
	cli.Cont[UserID].UnixSock.Read(tmp)
	res, _ := ioutil.ReadFile(cli.Cont[UserID].Volume+"/stdout")
	return string(res), true
}

func (cli *Client) CompileRequest(UserID, Lang, Req string) (bool) {
	var buf bytes.Buffer
	
	buf.WriteString(strings.ToUpper(Lang)+" "+strings.ToUpper(Req)+" "+cli.Cont[UserID].UserID)
	cli.Cont[UserID].UnixSock.Write(buf.Bytes())
	return true
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

func (cli *Client) DeleteContainer(UserID string) (bool) {
	err := cli.Pcli.ContainerStop(context.Background(), cli.Cont[UserID].ID, 0)
	if err != nil {
		Error.Println(err)
		return false
	}
	delete(cli.Cont, UserID)
	return true	
}

func (cli *Client) CopytoContainer(UserID, code, lang string) (bool) {
	f, err := os.Create(cli.Cont[UserID].Volume+"/main."+strings.ToLower(lang))
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
	l, err := net.Listen("unix", cli.Cont[UserID].Volume+"/host.sock")
	if err != nil {
		Error.Println(err)
		return false
	}
	err = cli.Pcli.ContainerStart(context.Background(), cli.Cont[UserID].ID)
	if err != nil {
		Error.Println(err)
		return false
	}

	cli.Cont[UserID].UnixSock, err = l.Accept()
	if err != nil {
		Error.Println(err)
		return false
	}	
	return true
}

func (cli *Client) AddContainer(UserID string) (bool) {
	var vol types.ContainerJSON
	var cont Container
	
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
	cli.Cont[UserID] = &cont
	
	return true
}
