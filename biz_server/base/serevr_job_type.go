package base

import (
	"strconv"
	"strings"
)

// ServerJobType 服务器职责类型
type ServerJobType int32

const (
	ServerJobTypeLogin ServerJobType = iota + 1
	ServerJobTypeGame
)

func (serverJobType ServerJobType) ToString() string {
	switch serverJobType {
	case ServerJobTypeLogin:
		return "LOGIN"
	case ServerJobTypeGame:
		return "GAME"
	default:
		return ""
	}
}

func StringToServerJobType(strVal string) ServerJobType {
	if strings.EqualFold(strVal, "LOGIN") ||
		strings.EqualFold(strVal, strconv.Itoa(int(ServerJobTypeLogin))) {
		return ServerJobTypeLogin
	}

	if strings.EqualFold(strVal, "GAME") ||
		strings.EqualFold(strVal, strconv.Itoa(int(ServerJobTypeGame))) {
		return ServerJobTypeGame
	}

	panic("无法转换为服务器类型职责！")
}

func StringToServerJobTypeArray(strVal string) []ServerJobType {
	if len(strVal) <= 0 {
		return nil
	}

	strArray := strings.Split(strVal, ",")
	var enumArray []ServerJobType
	for _, currStr := range strArray {
		serverJobType := StringToServerJobType(currStr)
		enumArray = append(enumArray, serverJobType)
	}
	return enumArray
}
