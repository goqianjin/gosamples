package onemodel

type PutUserTuneSwitchesReq string

type PutUserTuneSwitchesResp struct {
}

type GetUserTuneSwitchesReq struct {
}

type GetUserTuneSwitchesResp struct {
	TuneSwitches string `json:"tune_switches"`
}
