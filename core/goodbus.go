package goodbus

import (
	"errors"

	"github.com/nats-io/nats.go"
	"github.com/plgd-dev/kit/v2/log"
)

type GoodbusApplication struct {
	connected    bool
	connection   *nats.Conn
	errorhandler func(error)
	appid        string
}

func NewApplication(appid string) IApplication {
	return &GoodbusApplication{appid: appid, connected: false}
}

// 连接
func (app *GoodbusApplication) Auth(host, username, password string) error {
	connection, err := nats.Connect(host, func(o *nats.Options) error {
		o.Name = app.appid
		o.User = username
		o.Password = password
		o.AllowReconnect = true
		return nil
	})
	if err != nil {
		log.Error(err)
		return err
	}
	connection.SetClosedHandler(func(c *nats.Conn) {
		log.Error("Closed")
		if app.errorhandler != nil {
			app.errorhandler(errors.New("Closed"))
		}
	})
	connection.SetReconnectHandler(func(c *nats.Conn) {
		log.Error("Try Reconnect")
		if app.errorhandler != nil {
			app.errorhandler(errors.New("Reconnect"))
		}
	})
	connection.SetDisconnectHandler(func(c *nats.Conn) {
		log.Error("Disconnect")
		if app.errorhandler != nil {
			app.errorhandler(errors.New("Disconnect"))
		}
	})
	app.connection = connection
	app.connected = true
	return nil

}

//
func (app *GoodbusApplication) JoinChannel(channel string, callback func(msg Message)) error {
	if app.connected {
		_, err := app.connection.Subscribe(channel, func(m *nats.Msg) {
			callback(Message{
				Subject: m.Subject,
				Header:  m.Header,
				Data:    m.Data,
			})
			m.Ack()
		})
		if err != nil {
			log.Error(err)
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
	if app.connected {
		app.connection.Drain()
		app.connection.Close()
	}

}

//
func (app *GoodbusApplication) Publish(data []byte) error {
	if app.connected {
		return app.connection.Publish(app.appid, data)
	}
	return nil
}
