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

const (
	RTime time.Duration = 5 * time.Second
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

func (cli *Client) ExecuteProgram(UserID, code, lang, types string) (string, error) {
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
	err := cli.CopytoContainer(UserID, code, lang)
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
	return res, nil
}

// Check error value
func (cli *Client) GetResponse(UserID string) (string, error) {
	tmp := make([]byte, 1)

	t := time.Now()
	t = t.Add(RTime)
	err := cli.Cont[UserID].UnixSock.SetReadDeadline(t)
	if err != nil {
		return "", err
	}	
	cli.Cont[UserID].UnixSock.Read(tmp)
	res, _ := ioutil.ReadFile(cli.Cont[UserID].Volume+"/stdout")
	return string(res), nil
}

func (cli *Client) CompileRequest(UserID, Lang, Req string) (error) {
	var buf bytes.Buffer
	
	buf.WriteString(strings.ToUpper(Lang)+" "+strings.ToUpper(Req)+" "+cli.Cont[UserID].UserID)
	cli.Cont[UserID].UnixSock.Write(buf.Bytes())
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
	err := cli.Pcli.ContainerStop(context.Background(), cli.Cont[UserID].ID, 0)
	if err != nil {
		return err
	}
	delete(cli.Cont, UserID)
	return nil	
}

func (cli *Client) CopytoContainer(UserID, code, lang string) (error) {
	f, err := os.Create(cli.Cont[UserID].Volume+"/main."+strings.ToLower(lang))
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
	err = cli.Pcli.ContainerStart(context.Background(), cli.Cont[UserID].ID)
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
	resp, err := cli.Pcli.ContainerCreate(context.Background(), initConfig(), nil, nil, "")
	if err != nil {
		return err
	}
	cont.ID = resp.ID
	cont.Stream, err = cli.Pcli.ContainerAttach(context.Background(), types.ContainerAttachOptions{ContainerID: cont.ID, Stream: true, Stdout: true, Stderr: true})
	if err != nil {
		return err
	}
	vol, err = cli.Pcli.ContainerInspect(context.Background(), cont.ID)
	if err != nil {
		return err
	}

	cont.Volume = vol.Mounts[0].Source		
	cli.Cont[UserID] = &cont
	
	return nil
}
