// +build ignore

package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/naveego/live-api"
	"github.com/Sirupsen/logrus"
)

const (
	SERVER_ADDR = "ws://127.0.0.1:8888"
)

func main() {

	logrus.SetLevel(logrus.DebugLevel)
	
	cli, err := live.NewWebSocketClient(SERVER_ADDR, "23432", "TEST")
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, os.Interrupt)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("Connected")
	<-done
	fmt.Println("Exiting")

}

