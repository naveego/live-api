// +build ignore

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/naveego/live-api"
	"github.com/Sirupsen/logrus"
)

const (
	SERVER_ADDR = "127.0.0.1:9127"
)

func main() {

	logrus.SetLevel(logrus.DebugLevel)

	cli, err := live.NewTCPClient(SERVER_ADDR, "3242", "TEST")
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

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