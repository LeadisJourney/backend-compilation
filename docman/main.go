package main

import (
	"io"
	"os"
	"fmt"
	"log"
	"io/ioutil"

	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/api/types/container"
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
	//return &container.Config{Image: "leadis_image", Volumes: mount, Cmd: strslice.StrSlice{"touch /root/host/test"}, AttachStdout: true}
}

func findContainer(UserID string, conts []Container) (int) {
	for idx := range conts {
		if conts[idx].UserID == UserID {
			return idx
		}
	}
	return -1
}

func CopyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourcefile.Close()
	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destfile.Close()
	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}
	}
	return
}

func CopyDir(source string, dest string) (err error) {
	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}
	// create dest dir
	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}
	directory, _ := os.Open(source)
	objects, err := directory.Readdir(-1)
	for _, obj := range objects {
		sourcefilepointer := source + "/" + obj.Name()
		destinationfilepointer := dest + "/" + obj.Name()
		if obj.IsDir() {
			// create sub-directories - recursively
			err = CopyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// perform copy
			err = CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return
}

func main() {
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	Listener()
}
