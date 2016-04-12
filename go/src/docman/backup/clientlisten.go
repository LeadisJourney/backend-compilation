package main

import (
	"fmt"
	"time"
	"bufio"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

const Timeout int = 20

type Response struct {
	Status, Result string
	Errors, Warnings []string
}

type UserInfo struct {
	UserId string
	RequestId string
	Code string
	Language string
	Type string
}

type DocMan struct {
	cli *Client
	buf *bufio.Reader
	timeout int
}

func (d *DocMan) TimoutContainers() {
	fmt.Println("A container has timedout")
	for idx := range d.cli.Cont {
		fmt.Println(idx)
		t := time.Since(d.cli.Cont[idx].Time)
		if int(t.Seconds()) >= Timeout  {
			Trace.Println("Deleting", d.cli.Cont[idx].UserID)
			fmt.Println("Deleting", d.cli.Cont[idx].UserID)
			d.cli.DeleteContainer(idx)
			goto ENDLOOP
		}
	}
ENDLOOP:
	d.UpdateTimeout()
}

func (d *DocMan) UpdateTimeout() {
	Trace.Println("Updating Timeout")
	d.timeout = -1
	for idx := range d.cli.Cont {
		t := time.Since(d.cli.Cont[idx].Time)
		if d.timeout == -1 || d.timeout > (Timeout - int(t.Seconds())) {
			d.timeout = Timeout - int(t.Seconds())
		}
	}
	if d.timeout == -1 {
		d.timeout = Timeout
	}
	Trace.Println("New time:", d.timeout)
}

func (d *DocMan) CheckTime(c chan int) {
	for {
		select {
		case <- c:
			d.UpdateTimeout()			
		case <- time.After(time.Second * time.Duration(d.timeout)):
			Trace.Println("Timeout")
			d.TimoutContainers()
		}
	}
}

func (d *DocMan) Handler(w http.ResponseWriter, r *http.Request) {
	var user UserInfo
	var res Response
	
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error.Println(err)
	}
	err = json.Unmarshal(body, &user)
	res.Result = d.cli.ExecuteProgram(user.UserId, user.Code, user.Language, user.Type)
	res.Status = "OK"
	// fmt.Println(res.Result+"\n")
	//b, _ := json.Marshal(res)
	fmt.Fprintf(w, "%s: %s\n", user.UserId, res.Result)
	Trace.Println("Result sent to", user.UserId)
	fmt.Println("Result sent to", user.UserId)
}

func Listener() {
	var d DocMan

	c := make(chan int)
	d.timeout = Timeout
	d.cli = NewClient(c)
	http.HandleFunc("/v0.1/ce/status", d.Handler)
	go d.CheckTime(c)
	err := http.ListenAndServe(":2222", nil)
	if err != nil {
		Error.Println(err)
		return
	}
}
