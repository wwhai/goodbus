package goodbus

import (
	"errors"
	"time"

	"github.com/nats-io/nats.go"
)

type GoodbusApplication struct {
	connection   *nats.Conn
	errorhandler func(error)
	appid        string
}

func NewApplication(appid string) IApplication {
	return &GoodbusApplication{appid: appid}
}

// 连接
func (app *GoodbusApplication) Auth(host, username, password string) error {
	connection, err := nats.Connect(host, func(o *nats.Options) error {
		o.Name = app.appid
		o.User = username
		o.Password = password
		o.AllowReconnect = true
		o.ReconnectWait = 5 * time.Second
		return nil
	})
	if err != nil {
		return err
	}
	connection.SetErrorHandler(func(c *nats.Conn, s *nats.Subscription, err error) {
		if app.errorhandler != nil {
			app.errorhandler(err)
		}
	})
	connection.SetClosedHandler(func(c *nats.Conn) {
		if app.errorhandler != nil {
			app.errorhandler(errors.New("closed"))
		}
	})
	connection.SetReconnectHandler(func(c *nats.Conn) {
		if app.errorhandler != nil {
			app.errorhandler(errors.New("reconnect"))
		}
	})
	connection.SetDisconnectHandler(func(c *nats.Conn) {
		if app.errorhandler != nil {
			app.errorhandler(errors.New("disconnect"))
		}
	})
	app.connection = connection
	return nil

}

//
func (app *GoodbusApplication) JoinChannel(channel string, callback func(msg Message)) error {
	if app.connection.IsClosed() {
		_, err := app.connection.Subscribe(channel, func(m *nats.Msg) {
			callback(Message{
				Subject: m.Subject,
				Header:  m.Header,
				Data:    m.Data,
			})
			m.Ack()
		})
		if err != nil {
			app.errorhandler(err)
			return err
		}
		return nil
	} else {
		return errors.New("disconnected")
	}

}

//
func (app *GoodbusApplication) SetErrorHandler(errorhandler func(error)) {
	app.errorhandler = errorhandler
}

//
func (app *GoodbusApplication) Close() {
	if app.connection.IsConnected() {
		app.connection.Drain()
		app.connection.Close()
	}

}

//
// 发送数据
//
func (app *GoodbusApplication) Publish(data []byte) error {
	if app.connection.IsConnected() {
		return app.connection.Publish(app.appid, data)
	}
	return nil
}
