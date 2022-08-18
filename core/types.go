package goodbus

import "encoding/json"

/*
*
* 数据
*
 */
type Message struct {
	Subject string
	Header  map[string][]string
	Data    []byte
}

func (m *Message) String() string {
	bites, _ := json.Marshal(m)
	return string(bites)
}

/*
*
* 应用
*
 */
type IApplication interface {
	// 连接
	Auth(host, username, password string) error
	// 需要检查权限
	JoinChannel(string, func(msg Message)) error
	//
	SetErrorHandler(func(error))
	//
	Publish([]byte) error
	//
	Close()
}
