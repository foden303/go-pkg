package kcp

const (
	MsgTypePing uint32 = iota + 1
	MsgTypePong
	MsgTypeData
	MsgTypeAck
	MsgTypeError
	MsgTypeHeartbeat
)
