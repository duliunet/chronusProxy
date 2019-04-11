package model

import "sync"

/* 通信中心 */
type ConnHub struct {
	/* 此Tag为Hub的名称 */
	Tag string
	/* [ip] -> SocketConn */
	SyncMap ConnSyncMap
}

type ConnQHub struct {
	/* 此Tag为QHub的名称 */
	Tag string
	/* [ip] -> SocketConn */
	ConnQ *QueueS/* [SocketConn] */
}

/* 并发安全map */
type ConnSyncMap struct {
	Map map[string/*[Request().RemoteAddr]*/]interface{/*[model.SocketConn]*/}
	sync.RWMutex
}

/* Socket客户端 */
type SocketConn struct {
	/* 此Tag为配置文件中NeuronId */
	Tag  string
	Conn interface{ /* ws -> [*websocket.Conn] | tcp -> [net.Conn] | */}
}