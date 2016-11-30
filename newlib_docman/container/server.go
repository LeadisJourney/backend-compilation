package main

import (
	"os"
	"net"
	"log"
	"strings"
	"os/exec"
)

// Not clean
func main() {
	breq := make([]byte, 512)
	
	errf, _ := os.Create("/root/host/error")
	errl := log.New(errf,
                "ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	defer errf.Close()

	stdout, err := os.Create("/root/host/stdout")
	if err != nil {
		errl.Println(err)
		return
	}
	stderr, err := os.Create("/root/host/stderr")
	if err != nil {
		errl.Println(err)
		return
	}	
	
	conn, err := net.Dial("unix", "/root/host/host.sock")
	if err != nil {
		errl.Println(err)
		return
	}
	defer conn.Close()
	for {
		_, err := conn.Read(breq)
		if err != nil {
			errl.Println(err)
			return
		}
		req := strings.Split(string(breq), " ")
		stdout.Truncate(0)
		stderr.Truncate(0)
		if req[1] == "EXECUTION" {
			if req[0] == "C" {
				cmd := exec.Command("make", "realclean",  "-C", "/root/host/exercise")
				err = cmd.Start()
				if err != nil {
					errl.Println(err)
					return
				}
				cmd.Wait()
				cmd = exec.Command("make", "-C", "/root/host/exercise")
				cmd.Stdout = stdout
				cmd.Stderr = stderr
				err = cmd.Start()
				if err != nil {
					errl.Println(err)
					return
				}
				cmd.Wait()
			} else if req[0] == "CPP" {
				cmd := exec.Command("make", "realclean",  "-C", "/root/host/exercise")
				err = cmd.Start()
				if err != nil {
					errl.Println(err)
					return
				}
				cmd.Wait()
				cmd = exec.Command("make", "-C", "/root/host/exercise")
				cmd.Stdout = stdout
				cmd.Stderr = stderr
				err = cmd.Start()
				if err != nil {
					errl.Println(err)
					return
				}
				cmd.Wait()
			}
			if _, err = os.Stat("/root/host/exercise/res"); err == nil {
				stdout.Truncate(0)
				stderr.Truncate(0)

				err = os.Chdir("/root/host/exercise")
				if err != nil {
					errl.Println(err)
					return
				}
				cmd := exec.Command("./res")
				cmd.Stdout = stdout
				cmd.Stderr = stderr
				err = cmd.Start()
				if err != nil {
					errl.Println(err)
					return
				}
				cmd.Wait()
			} else {
				errl.Println(err)
				return
			}
		} else if req[1] == "COMPILATION" {
			if req[0] == "C" {
				cmd := exec.Command("make", "realclean",  "-C", "/root/host/exercise")
				err = cmd.Start()
				if err != nil {
					errl.Println(err)
					return
				}
				cmd.Wait()
				cmd = exec.Command("make", "-C", "/root/host/exercise")
				cmd.Stdout = stdout
				cmd.Stderr = stderr
				err = cmd.Start()
				if err != nil {
					errl.Println(err)
					return
				}
				cmd.Wait()
			} else if req[0] == "CPP" {
				cmd := exec.Command("make", "realclean",  "-C", "/root/host/exercise")
				err = cmd.Start()
				if err != nil {
					errl.Println(err)
					return
				}
				cmd.Wait()
				cmd = exec.Command("make", "-C", "/root/host/exercise")
				cmd.Stdout = stdout
				cmd.Stderr = stderr
				err = cmd.Start()
				if err != nil {
					errl.Println(err)
					return
				}
				cmd.Wait()
			}
			if _, err = os.Stat("/root/host/exercise/res"); err == nil {
				stdout.Truncate(0)
				stderr.Truncate(0)
				
				err = os.Chdir("/root/host/exercise")
				if err != nil {
					errl.Println(err)
					return
				}
				cmd := exec.Command("./res")
				cmd.Stdout = stdout
				cmd.Stderr = stderr
				err = cmd.Start()
				if err != nil {
					errl.Println(err)
					return
				}
				cmd.Wait()
			} else {
				errl.Println(err)
				return
			}			
		}
		stdout.Sync()
		stderr.Sync()
		conn.Write([]byte("1"))
	}
}
