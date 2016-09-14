package msg

type LoginReq struct {
	Name string
}

type LoginRep struct {
	Result string
}

type DoMatchReq struct {
}

type SelectActionReq struct {
	Action string
}
