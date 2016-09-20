package main

import (
	"client/network"
	"encoding/json"
	"fmt"
)

type UserData struct {
	state int //1已登录 2匹配中 3游戏中
	score int
}

func doLogin(name string) {
	var msg_body = network.LoginReq{}
	msg_body.Name = name
	network.CallProc(msg_body, "LoginRep", func(json_data json.RawMessage) {
		var rep = &network.LoginRep{}
		if err := json.Unmarshal(json_data, &rep); err != nil {
			panic(err)
		}
		fmt.Println(rep.Result)
	})
}

func tryMatch() {
	var msg_body = network.DoMatchReq{}
	network.SendProc(msg_body)
}

func doSelectAction(command string) {
	var msg_body = network.SelectActionReq{}
	msg_body.Action = command
	network.SendProc(msg_body)
}

func ClientOp() {
	fmt.Println("enter your name:")
	var name string
	fmt.Scan(&name)

	doLogin(name)

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

func main() {
	network.ConnectToServer()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			network.CloseConn()
		}
	}()
	go ClientOp()
	network.OnRecieveMsg()
}
