package main

import (
	"fmt"
	"time"
	"bufio"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

const (
	Timeout int = 500
	NCont int = 10
)

type Response struct {
	Status, Result string
	Errors, Warnings []string
	Graphic string
}

type UserInfo struct {
	UserId string
	RequestId string
	Code string
	Language string
	Type string
	Exercise string
}

type DocMan struct {
	cli *Client
	buf *bufio.Reader
	timeout int
}

func (d *DocMan) TimoutContainers() {
	for key, _ := range d.cli.Cont {
		t := time.Since(d.cli.Cont[key].Time)
		if int(t.Seconds()) >= Timeout && int(t.Hours()) < 2087 {
			Trace.Println("Deleting", d.cli.Cont[key].UserID)
			Info.Println("Deleting", d.cli.Cont[key].UserID, d.cli.Cont[key].Time)
			d.cli.DeleteContainer(key)
		}
	}
	d.UpdateTimeout()
}

func (d *DocMan) OldestContainer() {
	old := ""
	if len(d.cli.Cont) > 0 {
		for key, _ := range d.cli.Cont {
			if old == "" || time.Since(d.cli.Cont[old].Time) < time.Since(d.cli.Cont[key].Time) {
				old = key
			}
		}
	}
	d.cli.DeleteContainer(old)
}

func (d *DocMan) UpdateTimeout() {
	Trace.Println("Updating Timeout")
	d.timeout = -1
	for key, _ := range d.cli.Cont {
		t := time.Since(d.cli.Cont[key].Time)
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

// Parse request sent by API
func (d *DocMan) Handler(w http.ResponseWriter, r *http.Request) {
	var user UserInfo
	var res Response
	
	body, err := ioutil.ReadAll(r.Body)

	// TMP
	res.Errors = append(res.Errors, "Error")
	res.Warnings = append(res.Warnings, "Warning")
	// END TMP
	
	defer r.Body.Close()
	if err != nil {
		res.Status = "KO"
		b, _ := json.Marshal(res)
		fmt.Fprintf(w, "%s", b)
		Error.Println(err)
		return
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		res.Status = "KO"
		b, _ := json.Marshal(res)
		fmt.Fprintf(w, "%s", b)
		Error.Println(err)
		return
	}
	Info.Println("Received request from", user.UserId)
	tmp, err := d.cli.ExecuteProgram(user.UserId, user.Code, user.Language, user.Type, user.Exercise)
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", tmp)
	idx := strings.LastIndex(tmp, "\"Graphic\": [")
	idx2 := strings.LastIndex(tmp, "]")
	if idx != -1 && idx2 != -1 {
		res.Graphic = tmp[idx+12:idx2]
		res.Result = tmp[:idx]
	} else {
		res.Result = tmp
	}
	if err == nil {
		res.Status = "OK"
	} else {
		res.Status = "KO"
		d.cli.DeleteContainer(user.UserId)
		res.Errors = append(res.Errors, fmt.Sprint(err))
	}
	Info.Println("Result: ", res)
	b, _ := json.Marshal(res)
	fmt.Fprintf(w, "%s", b)
	Info.Println("Result sent to", user.UserId)
}

func Listener() {
	var d DocMan

	c := make(chan int)
	d.timeout = Timeout
	d.cli = NewClient(c)
	if d.cli == nil {
		return
	}
	// go MonListen(d.cli)
	http.HandleFunc("/v0.1/ce/status", d.Handler)
	go d.CheckTime(c)
	err := http.ListenAndServe(":8443", nil)
	if err != nil {
		Error.Println(err)
		return
	}
}
