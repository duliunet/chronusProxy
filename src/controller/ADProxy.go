/*
===========================================================================

===========================================================================
*/

package controller

import (
	"model"
	"modules/serial"
	"net/http"
)

/* ================================ DEFINE ================================ */
type ProxyS struct {
	Const struct {
		tag  string
		root string
	}
	isStarted bool
	neuron    *NeuronS
	mux       *http.ServeMux

	proxyConfig      map[string]interface{}
	tcp2tcpStopCArr  []chan bool
	udp2udpStopCArr  []chan bool
	udp2tcpStopCArr  []chan bool
	tcp2udpStopCArr  []chan bool
	uart2udpStopCArr []chan bool
}

/* ================================ PRIVATE ================================ */
/* 注册服务 */
func (mProxy *ProxyS) main() {
	/* Var */
	mProxy.tcp2tcpStopCArr = make([]chan bool, 0, 10)
	mProxy.udp2udpStopCArr = make([]chan bool, 0, 10)
	mProxy.udp2tcpStopCArr = make([]chan bool, 0, 10)
	mProxy.tcp2udpStopCArr = make([]chan bool, 0, 10)
	mProxy.uart2udpStopCArr = make([]chan bool, 0, 10)
	/* Func */
	mProxy.readConfig()
}

/* ================================ INTERFACE ================================ */

/* ================================ PROCESS ================================ */

/* 读取端口转发配置文件 */
func (mProxy *ProxyS) readConfig() {
	proxyHub := mProxy.neuron.Brain.Const.Proxy.ProxyHub
	if mProxy.neuron.Brain.CheckIsNull(proxyHub) {
		mProxy.Log("readConfig", "[proxyHub] -> Null")
		return
	}
	mProxy.proxyConfig = proxyHub
}

/* 基于TCP协议的端口转发 */
func (mProxy *ProxyS) runTCP2TCP() {
	if !mProxy.neuron.Brain.CheckIsNull(mProxy.proxyConfig["TCP"]) {
		for k, v := range mProxy.proxyConfig["TCP"].(map[string]interface{}) {
			stopC := make(chan bool)
			go mProxy.neuron.Express.TCPForward(k, v.(string), stopC)
			mProxy.tcp2tcpStopCArr = append(mProxy.tcp2tcpStopCArr, stopC)
		}
	}
}

func (mProxy *ProxyS) killTCP2TCP() {
	if !mProxy.neuron.Brain.CheckIsNull(mProxy.tcp2tcpStopCArr) {
		for _, v := range mProxy.tcp2tcpStopCArr {
			mProxy.neuron.Brain.ClearInterval(v)
		}
	}
}

/* 基于UDP协议的端口转发 */
func (mProxy *ProxyS) runUDP2UDP() {
	if !mProxy.neuron.Brain.CheckIsNull(mProxy.proxyConfig["UDP"]) {
		for k, v := range mProxy.proxyConfig["UDP"].(map[string]interface{}) {
			stopC := make(chan bool)
			go mProxy.neuron.Express.UDPForward(k, v.(string), stopC)
			mProxy.udp2udpStopCArr = append(mProxy.udp2udpStopCArr, stopC)
		}
	}
}

func (mProxy *ProxyS) killUDP2UDP() {
	if !mProxy.neuron.Brain.CheckIsNull(mProxy.udp2udpStopCArr) {
		for _, v := range mProxy.udp2udpStopCArr {
			mProxy.neuron.Brain.ClearInterval(v)
		}
	}
}

/* UDP转TCP协议 */
func (mProxy *ProxyS) runUDP2TCP() {
	if !mProxy.neuron.Brain.CheckIsNull(mProxy.proxyConfig["UDP2TCP"]) {
		for k, v := range mProxy.proxyConfig["UDP2TCP"].(map[string]interface{}) {
			stopC := make(chan bool)
			go mProxy.neuron.Express.UDP2TCPForward(k, v.(string), stopC)
			mProxy.udp2tcpStopCArr = append(mProxy.udp2tcpStopCArr, stopC)
		}
	}
}

func (mProxy *ProxyS) killUDP2TCP() {
	if !mProxy.neuron.Brain.CheckIsNull(mProxy.udp2tcpStopCArr) {
		for _, v := range mProxy.udp2tcpStopCArr {
			mProxy.neuron.Brain.ClearInterval(v)
		}
	}
}

/* TCP转UDP协议 */
func (mProxy *ProxyS) runTCP2UDP() {
	if !mProxy.neuron.Brain.CheckIsNull(mProxy.proxyConfig["TCP2UDP"]) {
		for k, v := range mProxy.proxyConfig["TCP2UDP"].(map[string]interface{}) {
			stopC := make(chan bool)
			go mProxy.neuron.Express.TCP2UDPForward(k, v.(string), stopC)
			mProxy.tcp2udpStopCArr = append(mProxy.tcp2udpStopCArr, stopC)
		}
	}
}

func (mProxy *ProxyS) killTCP2UDP() {
	if !mProxy.neuron.Brain.CheckIsNull(mProxy.tcp2udpStopCArr) {
		for _, v := range mProxy.tcp2udpStopCArr {
			mProxy.neuron.Brain.ClearInterval(v)
		}
	}
}

/* UART转UDP协议 */
func (mProxy *ProxyS) runUART2UDP() {
	if !mProxy.neuron.Brain.CheckIsNull(mProxy.proxyConfig["UART2UDP"]) {
		for k, v := range mProxy.proxyConfig["UART2UDP"].(map[string]interface{}) {
			stopC := make(chan bool)
			option := serial.OpenOptions{
				PortName:        v.(map[string]interface{})["PortName"].(string),
				BaudRate:        uint(v.(map[string]interface{})["BaudRate"].(float64)),
				DataBits:        uint(v.(map[string]interface{})["DataBits"].(float64)),
				StopBits:        uint(v.(map[string]interface{})["StopBits"].(float64)),
				MinimumReadSize: uint(v.(map[string]interface{})["MinimumReadSize"].(float64)),
			}
			go mProxy.neuron.Express.UART2UDPForward(k, option, stopC)
			mProxy.uart2udpStopCArr = append(mProxy.uart2udpStopCArr, stopC)
		}
	}
}

func (mProxy *ProxyS) killUART2UDP() {
	if !mProxy.neuron.Brain.CheckIsNull(mProxy.uart2udpStopCArr) {
		for _, v := range mProxy.uart2udpStopCArr {
			mProxy.neuron.Brain.ClearInterval(v)
		}
	}
}

/* ================================ TOOL ================================ */

/* ================================ SERVICE ================================ */

/* 构造服务 */
func (mProxy *ProxyS) service() {
	mProxy.runTCP2TCP()
	mProxy.runUDP2UDP()
	mProxy.runUDP2TCP()
	mProxy.runTCP2UDP()
	mProxy.runUART2UDP()
}

/* 析构服务 */
func (mProxy *ProxyS) serviceKiller() {
	mProxy.killTCP2TCP()
	mProxy.killUDP2UDP()
	mProxy.killUDP2TCP()
	mProxy.killTCP2UDP()
	mProxy.killUART2UDP()
}

/* ================================ PUBLIC ================================ */

/* 构造本体 */
func (mProxy *ProxyS) Ontology(neuron *NeuronS, mux *http.ServeMux, root string) *ProxyS {
	mProxy.neuron = neuron
	mProxy.mux = mux
	mProxy.Const.tag = root[1:]
	mProxy.Const.root = root
	mProxy.isStarted = neuron.Brain.Const.AutorunConfig.ADProxy
	if mProxy.isStarted {
		mProxy.neuron.Brain.SafeFunction(mProxy.main, func(err interface{}) {
			mProxy.neuron.Brain.MessageHandler(mProxy.Const.tag, "Main", 204, err)
		})
		mProxy.StartService()
	} else {
		mProxy.StopService()
	}
	return mProxy
}

/* 返回开关量 */
func (mProxy *ProxyS) IsStarted() bool {
	return mProxy.isStarted
}

/* 启动服务 */
func (mProxy *ProxyS) StartService() {
	mProxy.isStarted = true
	go mProxy.neuron.Brain.SafeFunction(mProxy.service, func(err interface{}) {
		mProxy.neuron.Brain.MessageHandler(mProxy.Const.tag, "StartService", 204, err)
	})
}

/* 停止服务 */
func (mProxy *ProxyS) StopService() {
	mProxy.isStarted = false
	go mProxy.neuron.Brain.SafeFunction(mProxy.serviceKiller, func(err interface{}) {
		mProxy.neuron.Brain.MessageHandler(mProxy.Const.tag, "StopService", 204, err)
	})
}

/* 打印信息 */
func (mProxy *ProxyS) Log(title string, content interface{}) {
	mProxy.neuron.Brain.LogGenerater(model.LogInfo, mProxy.Const.tag, title, content)
}
