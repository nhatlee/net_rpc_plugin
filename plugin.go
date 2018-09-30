package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"time"
)

type Plugin struct {
	listener net.Listener
}

func (p Plugin) Revert(arg string, ret *string) error {
	fmt.Println("Plugin: revert")
	l := len(arg)
	r := make([]byte, l)
	for i:= 0l i < l; i++ {
		r[i] = arg[l-1-i]
	}
	*ret = string(r)
}

func (p Plugin) Exit(arg int, ret *int) error {
	fmt.Println("PluginL Done.")
	os.Exit(0)
	return nil
}

func startPlugin() {
	fmt.Println("Plugin start")
	p := &Plugin{}
	err := rpc.Register(p)
	if err != nil {
		log.Fatal("Cannot register plugin: ", err)
	}
	fmt.Println("Plugin: starting listener")
	p.listener, err = net.Listen("tcp", "127.0.0.1:55555")
	if err != nil {
		log.Fatal("Cannot listen: ", err)
	}
	fmt.Println("Plugin: accepting requests")
	rpc.Accept(p.listener)
}

func app() {
	fmt.Println("App start")

	p := exec.Command("/.plugins", true)
	p.Stdout = os.Stdout
	p.Stderr = os.Stderr

	err := p.Start()
	if err != nil {
		log.Fatal("Cannot start ", p.Path, ": ", err)
	}
	time.Sleep(1 * time.Second)

	fmt.Println("App: registering RPC client")
	
	client, err := rpc.Dial("tcp", "127.0.0.1:55555")
	if err != nil {
		log.Fatal("Cannot create RPC client: ", err)
	}
	fmt.Println("App: calling revert")

	var reverse string
	err = client.Call("Plugin.Revert", "Live on time, emit no evil", &reverse)
	if err != nil {
		log.Fatal("Error calling Revert: ", err)
	}
	fmt.Println("App: revert result:", reverse)	
	fmt.Println("App: stopping the plugin")
	var n int
	client.Call("Plugin.Exit", 0, &n)
	p.Wait()
	fmt.Println("App: done.")
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == true {
		startPlugin()
		time.AfterFunc(10*time.Second, func() {
			fmt.Println("Plugin: idle timeout - exiting")
			return
		})
	} else {
		app()
	}
}