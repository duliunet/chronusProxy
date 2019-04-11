package model

import (
	"net"
	"time"
)

/* Log Type */
const (
	LogInfo int = iota
	LogDebug
	LogTrace
	LogWarn
	LogError
	LogCritical
)

/* Neuron Service Type */
type ServiceS struct {
	Tag   string
	Alive chan bool
}

/* 时间信息 */
type TimeS struct {
	YMD           string
	Week          string
	Time          string
	Timestamp     string
	TimestampMill string
	TimestampNano string
	Datetime      time.Time
}

/* ICMP包数据结构 */
type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

/* UDP通信消息数据结构 */
type UDPPacket struct {
	Addr *net.UDPAddr
	Msg []byte
}

/* 内部消息 */
type MessageS struct {
	Code    int
	Message string
	Data    interface{}
}