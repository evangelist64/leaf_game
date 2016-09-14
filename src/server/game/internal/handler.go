package internal

import (
	"reflect"
	"server/msg"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
)

func init() {
	handler(&msg.LoginReq{}, handleLogin)
	handler(&msg.DoMatchReq{}, handleDoMatch)
	handler(&msg.SelectActionReq{}, handleSelectAction)
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handleLogin(args []interface{}) {
	m := args[0].(*msg.LoginReq)
	a := args[1].(gate.Agent)
	a.SetUserData(m.Name)
	log.Debug("login,name:%v", m.Name)
	a.WriteMsg(&msg.LoginRep{Result: "login ok"})
}

func handleDoMatch(args []interface{}) {
	//m := args[0].(*msg.DoMatchReq)
	a := args[1].(gate.Agent)

	log.Debug("do match,%v", a.UserData())
}
func handleSelectAction(args []interface{}) {
	m := args[0].(*msg.SelectActionReq)
	//a := args[1].(gate.Agent)

	log.Debug("do select,%v", m.Action)
}
