package main

import (
	"bytes"
	"client/msg"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
)

var conn net.Conn

func connectToServer() net.Conn {
	conn, err := net.Dial("tcp", "127.0.0.1:3563")
	if err != nil {
		panic(err)
	}
	return conn
}

func sendData(msg_body interface{}) {
	fmt.Println(msg_body)
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
	fmt.Println(string(data))
	// len + data
	m := make([]byte, 2+len(data))

	// 默认使用大端序
	binary.BigEndian.PutUint16(m, uint16(len(data)))

	copy(m[2:], data)
	// 发送消息
	conn.Write(m)
}

func tryMatch() {
	var msg_body = msg.DoMatchReq{}
	sendData(msg_body)
}

func doSelectAction(command string) {
	var msg_body = msg.SelectActionReq{}
	msg_body.Action = command
	sendData(msg_body)
}

func main() {
	conn = connectToServer()
	for {
		var command string
		fmt.Scan(&command)

		switch command {
		case "match":
			tryMatch()
			break
		case "1", "2", "3":
			doSelectAction(command)
			break
		default:
			fmt.Println("unknown command")
			break
		}
	}
}
