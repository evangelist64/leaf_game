package internal

import (
	"fmt"
)

func init() {
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
}

func rpcNewAgent(args []interface{}) {
	fmt.Println("new agent")
	//a := args[0].(gate.Agent)
	//_ = a

}

func rpcCloseAgent(args []interface{}) {
	fmt.Println("close agent")
	//a := args[0].(gate.Agent)
	//_ = a
}
