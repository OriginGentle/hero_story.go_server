package base

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ICmdContext 自定义指令上下文接口
type ICmdContext interface {

	// BindUserId 绑定用户id
	BindUserId(val int64)

	// GetUserId 获取用户id
	GetUserId() int64

	// GetClientIpAddr 获取客户端 IP 地址
	GetClientIpAddr() string

	// Write 写出消息对象
	Write(msgObj protoreflect.ProtoMessage)

	// SendError 发送错误消息
	SendError(errorCode int, errorInfo string)

	// Disconnect 断开连接
	Disconnect()
}
