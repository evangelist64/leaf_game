package main

import (
	"client/network"
	"encoding/json"
	"fmt"
	"time"
)

const (
	state_normal      = iota
	state_in_matching = iota
	state_in_counting = iota
	state_in_game     = iota
)

type UserData struct {
	state int
}

var my_data UserData
var match_timer *time.Timer

func doLogin(name string) {
	var msg_body = network.LoginReq{}
	msg_body.Name = name
	network.CallProc(msg_body, "LoginRep", func(json_data json.RawMessage) {
		var rep = &network.LoginRep{}
		if err := json.Unmarshal(json_data, &rep); err != nil {
			panic(err)
		}
		if rep.Result == "login ok" {
			fmt.Println("-------" + rep.Result + "-------")
			my_data.state = state_normal
			go ClientOp()
		} else {
			fmt.Println("can not login, please try again later")
		}
	})
}

func tryMatch() {
	var msg_body = network.DoMatchReq{}
	network.CallProc(msg_body, "DoMatchRep", func(json_data json.RawMessage) {
		var rep = &network.DoMatchRep{}
		if err := json.Unmarshal(json_data, &rep); err != nil {
			panic(err)
		}
		fmt.Println(rep.Result)
		if rep.Result == "match ok" {
			fmt.Println("-------" + rep.Result + "-------")
			fmt.Println("you will fight againest \"" + rep.Enemy_name + "\"")
			match_timer.Stop()

			startGame()
		} else {
			fmt.Println("can not match anyone, please try again later")
		}
	})

	my_data.state = state_in_matching
	fmt.Println("in matching...")
	match_timer = time.AfterFunc(time.Second*30, onMatchTimeout)
}

func onMatchTimeout() {
	fmt.Println("can not match anyone, please try again later")
	my_data.state = state_normal
}

func startGame() {
	network.RegisterFunc("FireActionRep", func(json_data json.RawMessage) {
		var rep = &network.FireActionRep{}
		if err := json.Unmarshal(json_data, &rep); err != nil {
			panic(err)
		}
		fmt.Println(rep.Result)
		my_data.state = state_normal
	})

	my_data.state = state_in_game
	fmt.Println("press enter button to fire, when you think the count reaches 0")
	count := 10
	count_down_timer := time.NewTicker(time.Second)
	go func() {
		for _ = range count_down_timer.C {
			fmt.Println(count)
			count = count - 1
			if count <= 5 {
				count_down_timer.Stop()
				break
			}
		}
	}()
}

func doFightAction() {
	var msg_body = network.FireActionReq{}
	network.SendProc(msg_body)
}

func ClientOp() {

	fmt.Println("press \"match\" to start a game")
	for {
		var command string
		fmt.Scanln(&command)

		if my_data.state == state_normal && command == "match" {
			tryMatch()
		} else if my_data.state == state_in_game {
			doFightAction()
		} else {
			fmt.Println("can not do command:" + command)
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

	fmt.Println("Welcome to cowboy dule game!")
	fmt.Println("please enter your name:")
	var name string
	fmt.Scanln(&name)
	doLogin(name)

	network.OnRecieveMsg()
}
