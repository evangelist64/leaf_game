package msg

type LoginReq struct {
	Name string
}

type LoginRep struct {
	Result string
}

type DoMatchReq struct {
}

type DoMatchRep struct {
	Enemy_name string
	Result     string
}

type FireActionReq struct {
}

type FireActionRep struct {
	Result string
}
