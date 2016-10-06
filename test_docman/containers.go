package main

import (
	"os"
	"fmt"
	"net"
	"time"
	"bytes"
	"errors"
	"strings"
	"io/ioutil"
	"golang.org/x/net/context"
	
	"github.com/fsouza/go-dockerclient"
)

const (
	RTime time.Duration = 60 * time.Second
)

type Container struct {
	ID, UserID, Volume string
	UnixSock net.Conn
	Time time.Time
	PCPU uint64
	PSys uint64
}


type Client struct {
	Pcli *docker.Client
	Cont map[string]*Container
	c chan int
}

func NewClient(c chan int) (*Client) {
	var cli Client
	var err error
	
	cli.c = c
	cli.Cont = make(map[string]*Container)

	cli.Pcli, err = docker.NewClientFromEnv()
	if err != nil {
		Error.Println(err)
		return nil
	}

	cwd, _ := os.Getwd()
	// ctx, err := os.Open(cwd+"/Dockerfile.tar.gz")
	// if err != nil {
	// 	Error.Println(err)
	// 	return nil
	// }
	err = cli.Pcli.BuildImage(docker.BuildImageOptions{Dockerfile: cwd+"/Dockerfile.tar.gz", Name: "leadis_image", SuppressOutput: false, Context: context.Background(),})
	if err != nil {
		Error.Println(err)
		return nil
	}
	
	return &cli
}

// Creates and prepares Docker conatiner then compile and execute in container
func (cli *Client) ExecuteProgram(UserID, code, lang, types, ex string) (string, error) {
	var res string
	
	if _, ok := cli.Cont[UserID]; !ok {
		// TEST
		// if len(cli.Cont) >= 1 {
		// 	cli.OldestContainer()
		// }
		// END TEST
		err := cli.AddContainer(UserID)
		if err != nil {
			return "", err
		}
		
		err = cli.StartContainer(UserID)
		if err != nil {
			return "", err
		}
	}

	cli.Cont[UserID].Time = time.Now()

	err := cli.CopytoContainer(UserID, code, lang, ex)
	if err != nil {
		return "", err
	}

	err = cli.CompileRequest(UserID, lang, types)
	if err != nil {
		return "", err
	}

	res, err = cli.GetResponse(UserID)
	if err != nil {
		return "", err
	}
	cli.c <- 1

	
	// BEGIN STATS TEST
	// cli.ReadStats(UserID)
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
		Error.Println(err)
		return "", errors.New("Internal Error!")
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
		Error.Println(err)
		return errors.New("Internal Error!")
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
	err := cli.Pcli.StopContainer(cli.Cont[UserID].ID, 0)
	if err != nil {
		Error.Println(err)
		return errors.New("Internal Error!")
	}
	delete(cli.Cont, UserID)
	return nil	
}

func (cli *Client) CopytoContainer(UserID, code, lang, ex string) (error) {
	os.Remove(cli.Cont[UserID].Volume+"/exercise")
	CopyDir(ex, cli.Cont[UserID].Volume+"/exercise")
	f, err := os.Create(cli.Cont[UserID].Volume+"/exercise/src/"+ex+"."+strings.ToLower(lang))
	if err != nil {
		Error.Println(err)
		return errors.New("Internal Error!")
	}
	_, err = f.Write([]byte(code))
	if err != nil {
		Error.Println(err)
		return errors.New("Internal Error!")
	}
	return nil
}

func (cli *Client) StartContainer(UserID string) (error) {
	l, err := net.Listen("unix", cli.Cont[UserID].Volume+"/host.sock")
	if err != nil {
		Error.Println(err)
		return errors.New("Internal Error!")
	}
	err = cli.Pcli.StartContainer(cli.Cont[UserID].ID, nil)
	if err != nil {
		Error.Println(err)
		return errors.New("Internal Error!")
	}

	cli.Cont[UserID].UnixSock, err = l.Accept()
	if err != nil {
		Error.Println(err)
		return errors.New("Internal Error!")
	}
	return nil
}

func (cli *Client) AddContainer(UserID string) (error) {
	var cont Container
	
	cont.UserID = UserID
	// REMOVE
	cont.PCPU = 0
	cont.PSys = 0
	// END REMOVE
	resp, err := cli.Pcli.CreateContainer(docker.CreateContainerOptions{"leadis journey", initConfig(), nil, nil, context.Background()})
	if err != nil {
		Error.Println(err)
		return errors.New("Internal Error!")
	}
	cont.ID = resp.ID
	err = cli.Pcli.AttachToContainer(docker.AttachToContainerOptions{Container: cont.ID, Stream: true, Stdout: true, Stderr: true})
	if err != nil {
		Error.Println(err)
		return errors.New("Internal Error!")
	}
	vol, err := cli.Pcli.InspectContainer(cont.ID)
	if err != nil {
		Error.Println(err)
		return errors.New("Internal Error!")
	}

	//cont.Volume = vol.Mounts[0].Source
	cont.Volume = vol.Mounts[0].Source
	fmt.Println(cont.Volume)
	cli.Cont[UserID] = &cont
	Info.Println("Creating New Container, ID: ", UserID)
	return nil
}
