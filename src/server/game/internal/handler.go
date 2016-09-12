package internal

import (
	"reflect"
	"server/msg"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
)

func init() {
	// 向当前模块（game 模块）注册 Hello 消息的消息处理函数 handleHello
	handler(&msg.DoMatchReq{}, handleDoMatch)
	handler(&msg.SelectActionReq{}, handleSelectAction)

}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handleDoMatch(args []interface{}) {
	//m := args[0].(*msg.DoMatchReq)
	a := args[1].(gate.Agent)

	log.Debug("%v", a.UserData())
}
func handleSelectAction(args []interface{}) {
	m := args[0].(*msg.SelectActionReq)
	//a := args[1].(gate.Agent)

	log.Debug("%v", m.Action)
}
