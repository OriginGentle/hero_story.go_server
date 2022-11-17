package base

type BizServerData struct {
	ServerId           int32           `json:"serverId"`
	ServerAddr         string          `json:"serverAddr"`
	ServerJobTypeArray []ServerJobType `json:"serverJobTypeArray"`
	LoadCount          int32           `json:"loadCount"`
}
