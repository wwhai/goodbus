# goodbus: 面向应用的物联网消息总线
![1661501655788](image/readme/1661501655788.png)
## 简介

goodbus是基于Nats封装的一套应用层消息总线系统，主要用来分发物联网设备数据。
## 使用
```go
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
	if err := app.Auth("127.0.0.1:4222", "nats_client", "password"); err != nil {
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

```