package main

import (
	"fmt"
	goodbus "goodbus/core"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT)
	app := goodbus.NewApplication("a-demo-appid")
	if err := app.Auth("127.0.0.1:4222", "nats_client", "Pwa43zr2kS"); err != nil {
		panic(err)
	}
	app.SetErrorHandler(func(err error) {
		fmt.Println(err)
	})
	if err := app.JoinChannel("mychannelabcd", func(msg goodbus.Message) {
		fmt.Println("Received:", msg.String())

	}); err != nil {
		panic(err)
	}
	signal := <-c
	fmt.Printf("Received stop signal:%v", signal)
	os.Exit(1)
}
