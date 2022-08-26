package goodbus

import (
	"fmt"
	"testing"
	"time"
)

func Test_goodbus(t *testing.T) {
	app := NewApplication("a-demo-appid")
	app.SetErrorHandler(func(err error) {
		t.Error(err)
	})
	if err := app.Auth("127.0.0.1:4222", "nats_client", "Pwa43zr2kS"); err != nil {
		t.Fatal(err)
	}
	if err := app.JoinChannel("mychannelabcd", func(msg Message) {
		fmt.Println("Received:", msg.String())
	}); err != nil {
		t.Fatal(err)
	}
	if err := app.Publish([]byte("turn-off-light")); err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)
}
