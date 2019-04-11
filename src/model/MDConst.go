package model

/*
===========================================================================
Dev Aim
===========================================================================

===========================================================================
Change List
===========================================================================
Version
0.8.7
***** UDP、TCP与UART协议转化
**** 修改DDExpress -> 串口服务器取消，固定式机器人该用UART协议直接通讯（类似电机控制）

===========================================================================
*/

/* ================================ DEFINE ================================ */
type autorunS struct {
	ADProxy bool
	SDExampleInterface bool
}

type proxyS struct {
	ProxyHub map[string]interface{}
}

type fileS struct {
	Chmod    int
	TempPath string
}

type requestS struct {
	DefaultHeader map[string][]string
}

type serverS struct {
	Host       string
	Port       int
	StaticPath string
	UploadPath string
	XPoweredBy string
	ACAO       bool
}

type tlsServerS struct {
	Open        bool
	TLSPort     int
	TLSCertPath string
}

type wsParamS struct {
	Interval   int
	BufferSize int
}

type tcpParamS struct {
	Interval   int
	BufferSize int
}

type udpParamS struct {
	MaxLen     int
	Interval   int
	BufferSize int
}

type uartParamS struct {
	Interval   int
	BufferSize int
}

type intervalS struct {
	HZ25TimerInterval      int
	HZ8TimerInterval       int
	HZ4TimerInterval       int
	HZ2TimerInterval       int
	HZ1TimerInterval       int
	CommanderTimerInterval int
	LooperTimerInterval    int
	SystemTimerInterval    int
	RetryTimerInterval     int
	MicroMsgTokenInterval  int
}

/* ================================ PUBLIC ================================ */

type Const struct {
	RunEnv     int
	Version    string
	NeuronId   string
	NeuronMail string

	AutorunConfig autorunS
	ErrorCode     map[int]string

	Proxy       proxyS
	File        fileS
	HTTPRequest requestS
	HTTPServer  serverS
	HTTPS       tlsServerS

	WSParam   wsParamS
	TCPParam  tcpParamS
	UDPParam  udpParamS
	UartParam uartParamS
	Interval  intervalS
}

/* 构造本体 */
func (*Const) Ontology() Const {
	return Const{
		0,
		"0.8.7",
		"ChronusProxy",
		"adumandix@gmail.com",
		/* 自启动配置 */
		autorunS{
			true,
			true,
		},
		/* 错误代码 */
		map[int]string{
			100: "Success",
			101: "System Running",
			102: "System ShutDown",
			103: "Process ShutDown",
			104: "Process Timeout",
			105: "Return Stack",

			200: "Failed",
			201: "Interface Banned",
			202: "JSON Error",
			203: "Command Error",
			204: "System Error",
			205: "File Error",
			206: "Buffer Error",
			207: "Request Error",
			208: "Auth Error",
			209: "Encode/Decode Error",
			210: "TCP Conn Error",
			211: "UDP Conn Error",
			212: "Url Error",
			213: "Path Error",
			214: "Websocket Error",
			215: "Transform Error",
			216: "IOReader Error",
			217: "Platform Error",
			218: "Exec Error",
			219: "Gzip Error",
			220: "Null Error",
			221: "ContentType Error",
			222: "UART Error",

			300: "Database Disconnected",
			301: "Query Error",
			302: "TransAction Error",
			303: "RollBack Error",
			304: "Commit Error",
			305: "Data Analyze Error",

			400: "Redis Disconnected",
			401: "Redis Error",
		},
		proxyS{
			map[string]interface{}{
				"TCP": map[string]interface{}{},
				"UDP": map[string]interface{}{},
				"TCP2UDP": map[string]interface{}{},
				"UDP2TCP": map[string]interface{}{},
				"UART2UDP": map[string]interface{}{},
			},
		},
		fileS{
			0766,
			"/static/temp/",
		},
		requestS{
			map[string][]string{
				"User-Agent": {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36"},
			},
		},
		serverS{
			"0.0.0.0",
			8800,
			"/static",
			"/static/upload",
			"Chronus Express",
			/* 跨域标识 */
			false,
		},
		tlsServerS{
			false,
			8443,
			"/tls/tls",
		},
		wsParamS{
			12e4,
			2 << 20,
		},
		tcpParamS{
			12e4,
			2 << 12,
		},
		udpParamS{
			1,
			12e4,
			2 << 12,
		},
		uartParamS{
			12e4,
			2 << 12,
		},
		intervalS{
			40,
			125,
			250,
			500,
			1000,
			100,
			3000,
			30000,
			5000,
			7200000,
		},
	}
}
