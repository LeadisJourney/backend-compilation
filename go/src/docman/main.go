package main

import (
	"io"
	"os"
	"log"
	"io/ioutil"

	"github.com/docker/engine-api/types/strslice"
	"github.com/docker/engine-api/types/container"
)

var (
	Trace   *log.Logger
	Info   *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	
	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func initConfig() (config *container.Config) {
	mount := map[string]struct{}{"/root/host": {}}
	return &container.Config{Image: "leadis_image", Volumes: mount, Cmd: strslice.StrSlice{"/root/server.py"}, AttachStdout: true}
}

func findContainer(UserID string, conts []Container) (int) {
	for idx := range conts {
		if conts[idx].UserID == UserID {
			return idx
		}
	}
	return -1
}

func main() {
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	Listener()
}
