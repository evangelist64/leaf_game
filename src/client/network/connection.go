package network

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"reflect"
)

var conn net.Conn
var proc_map map[string]func(json.RawMessage)

func init() {
	proc_map = make(map[string]func(json.RawMessage))
}

func CallProc(send_msg interface{}, recv_id string, callback func(json.RawMessage)) {
	proc_map[recv_id] = callback
	SendProc(send_msg)
}

func RegisterFunc(recv_id string, callback func(json.RawMessage)) {
	proc_map[recv_id] = callback
}

func SendProc(msg_body interface{}) {
	msg_id := reflect.TypeOf(msg_body).Name()
	req_body_bytes, err := json.Marshal(msg_body)
	if err != nil {
		panic(err)
	}
	b := bytes.Buffer{}
	b.WriteString("{\"" + msg_id + "\":")
	b.WriteString(string(req_body_bytes))
	b.WriteString("}")

	var data = b.Bytes()
	// len + data
	m := make([]byte, 2+len(data))

	// 默认使用大端序
	binary.BigEndian.PutUint16(m, uint16(len(data)))

	copy(m[2:], data)
	// 发送消息
	conn.Write(m)
}

func ConnectToServer() {
	c, err := net.Dial("tcp", "127.0.0.1:3563")
	if err != nil {
		panic(err)
	}
	conn = c
}

func OnRecieveMsg() {
	for {
		var bytes_len [2]byte
		if _, err := io.ReadFull(conn, bytes_len[0:2]); err != nil {
			panic(err)
		}
		msg_len := binary.BigEndian.Uint16(bytes_len[0:2])
		bytes := make([]byte, msg_len)
		if _, err := io.ReadFull(conn, bytes); err != nil {
			panic(err)
		}

		var dat map[string]json.RawMessage
		if err := json.Unmarshal(bytes, &dat); err != nil {
			panic(err)
		}
		for k, v := range dat {
			f := proc_map[k]
			if f != nil {
				f(v)
			} else {
				fmt.Println("can't find func by index " + k)
			}
		}
	}
}

func CloseConn() {
	conn.Close()
}
