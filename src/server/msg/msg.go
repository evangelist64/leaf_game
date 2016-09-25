package msg

import (
	"github.com/name5566/leaf/network/json"
)

// 使用默认的 JSON 消息处理器（默认还提供了 protobuf 消息处理器）
var Processor = json.NewProcessor()

func init() {
	Processor.Register(&LoginReq{})
	Processor.Register(&DoMatchReq{})
	Processor.Register(&FireActionReq{})
	Processor.Register(&LoginRep{})
	Processor.Register(&DoMatchRep{})
	Processor.Register(&FireActionRep{})
}
