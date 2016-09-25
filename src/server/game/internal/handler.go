package internal

import (
	"container/list"
	"fmt"
	"reflect"
	"server/msg"
	"strconv"
	"time"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
)

//user state
const (
	state_normal      = iota
	state_in_matching = iota
	state_game_start  = iota
	state_fired       = iota
)

type game struct {
	user_agent_1 gate.Agent
	user_agent_2 gate.Agent
	start_time   int64
	result       string
}

var match_queue *list.List

func init() {
	handler(&msg.LoginReq{}, handleLogin)
	handler(&msg.DoMatchReq{}, handleDoMatch)
	handler(&msg.FireActionReq{}, handleFireAction)
	match_queue = list.New()
	go doMatchLoop()
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handleLogin(args []interface{}) {
	property_map := make(map[string]interface{})

	m := args[0].(*msg.LoginReq)
	a := args[1].(gate.Agent)
	property_map["state"] = state_normal
	property_map["name"] = m.Name
	a.SetUserData(property_map)

	log.Debug("user login:%v", m.Name)
	a.WriteMsg(&msg.LoginRep{Result: "login ok"})
}

func handleDoMatch(args []interface{}) {
	a := args[1].(gate.Agent)
	match_queue.PushBack(a)
	property_map := a.UserData().(map[string]interface{})
	property_map["state"] = state_in_matching

	log.Debug("user do match:%v", property_map["name"])
}
func handleFireAction(args []interface{}) {
	//m := args[0].(*msg.FireActionReq)
	a := args[1].(gate.Agent)
	property_map := a.UserData().(map[string]interface{})
	if property_map["state"] == state_game_start {
		property_map["time"] = time.Now().UnixNano()
		property_map["state"] = state_fired
	} else if property_map["state"] == state_fired {
		a.WriteMsg(&msg.FireActionRep{Result: "already fired"})
	} else {
		a.WriteMsg(&msg.FireActionRep{Result: "not in game"})
	}
}

func doMatchLoop() {
	for {
		if match_queue.Len() >= 2 {
			elem1 := match_queue.Front()
			user_agent_1 := elem1.Value.(gate.Agent)
			match_queue.Remove(elem1)

			elem2 := match_queue.Front()
			user_agent_2 := elem2.Value.(gate.Agent)
			match_queue.Remove(elem2)

			property_map_1 := user_agent_1.UserData().(map[string]interface{})
			property_map_2 := user_agent_2.UserData().(map[string]interface{})

			user_agent_1.WriteMsg(&msg.DoMatchRep{Enemy_name: property_map_2["name"].(string), Result: "match ok"})
			user_agent_2.WriteMsg(&msg.DoMatchRep{Enemy_name: property_map_1["name"].(string), Result: "match ok"})

			//begin a game
			game_one := new(game)
			game_one.user_agent_1 = user_agent_1
			game_one.user_agent_2 = user_agent_2
			property_map_1["time"] = 0
			property_map_2["time"] = 0
			property_map_1["state"] = state_game_start
			property_map_2["state"] = state_game_start
			startGame(game_one)

		} else {
			time.Sleep(time.Second)
		}
	}
}

func startGame(game_one *game) {
	fmt.Println("start game")
	game_one.start_time = time.Now().UnixNano()

	//prepare 11s, timeout 4s
	time.AfterFunc(time.Second*15, func() {
		property_map_1 := game_one.user_agent_1.UserData().(map[string]interface{})
		property_map_2 := game_one.user_agent_2.UserData().(map[string]interface{})

		var time_1 int64 = 0
		var time_2 int64 = 0
		if property_map_1["time"] != 0 {
			time_1 = property_map_1["time"].(int64)
		}
		if property_map_2["time"] != 0 {
			time_2 = property_map_2["time"].(int64)
		}
		var time_diff_1 int64 = time_1 - game_one.start_time
		var time_diff_2 int64 = time_2 - game_one.start_time
		var user_foul_1 bool = time_diff_1 < 11*1000000000
		var user_foul_2 bool = time_diff_2 < 11*1000000000
		var time_diff_after_begin_1 int64 = time_diff_1 - 11*1000000000
		var time_diff_after_begin_2 int64 = time_diff_2 - 11*1000000000

		var str_user_1 string = property_map_1["name"].(string) + " fired after begin:" + strconv.FormatInt(time_diff_after_begin_1, 10) + "ns\n"
		var str_user_2 string = property_map_2["name"].(string) + " fired after begin:" + strconv.FormatInt(time_diff_after_begin_2, 10) + "ns\n"
		if user_foul_1 {
			str_user_1 = property_map_1["name"].(string) + " fouled\n"
		}
		if user_foul_2 {
			str_user_2 = property_map_2["name"].(string) + " fouled\n"
		}

		if user_foul_1 && user_foul_2 {
			game_one.result = "Result:draw\n"
		} else if user_foul_1 && !user_foul_2 {
			game_one.result = "Result:" + property_map_2["name"].(string) + " win\n"
		} else if !user_foul_1 && user_foul_2 {
			game_one.result = "Result:" + property_map_1["name"].(string) + " win\n"
		} else {
			if time_diff_1 < time_diff_2 {
				game_one.result = "Result:" + property_map_1["name"].(string) + " win\n"
			} else if time_diff_1 > time_diff_2 {
				game_one.result = "Result:" + property_map_2["name"].(string) + " win\n"
			} else {
				game_one.result = "Result:draw\n"
			}
		}
		game_one.result = "---------Game Result---------\n" + game_one.result + str_user_1 + str_user_2

		game_one.user_agent_1.WriteMsg(&msg.FireActionRep{Result: game_one.result})
		game_one.user_agent_2.WriteMsg(&msg.FireActionRep{Result: game_one.result})

		property_map_1["state"] = state_normal
		property_map_2["state"] = state_normal
	})
}
