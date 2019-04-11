package model

/* Saas服务接口 */
type ExpressI interface {
	// 获取服务启动状态
	IsStarted() bool
	// 启动服务
	StartService()
	// 停止服务
	StopService()
}

/* WS通信接口 */
type WebsocketI interface {
	// 服务器容器
	WSHub() ConnHub
	// 消息解析模块
	GMessageHandler(interface{}, interface{})
}