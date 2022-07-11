package base

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ICmdContext interface {
	BindUserId(val int64)

	GetUserId() int64

	GetClientIpAddr() string

	Write(msgObj protoreflect.ProtoMessage)

	SendError(errorCode int, errorInfo string)

	Disconnect()
}
