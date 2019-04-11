package main

import (
	"controller"
	"fmt"
	"model"
	"modules/logs/logger"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

/* ================================ DEFINE ================================ */

const tag = "Neuron"

/* Application */
var application struct {
	looperStopC chan bool
	// 神经元指针
	neuron *controller.NeuronS
	// 系统信号量
	sysExitSignalC chan os.Signal
	// 服务系统
	serviceHub *model.QueueS
}

/* ================================ EVENT ================================ */

/* 循环事件 */
func looperEvent() {
	// 开发环境显示服务器数量
	application.neuron.Brain.LogGenerater(model.LogInfo, tag, "LooperEvent", fmt.Sprintf("Service Alive Count -> [%d]", application.serviceHub.Len()))
	// 监视服务
	application.looperStopC = make(chan bool)
	application.neuron.Brain.SetInterval(func() (int, interface{}) {
		for _, v := range application.serviceHub.ToArray() {
			service := v.Value.(*model.ServiceS)
			select {
			case data := <-service.Alive:
				if !data {
					application.neuron.Brain.LogGenerater(model.LogError, tag, "LooperEvent", fmt.Sprintf("Service Dead -> [%s]", service.Tag))
					application.serviceHub.Delete(v)
					// 显示服务器数量
					application.neuron.Brain.LogGenerater(model.LogInfo, tag, "LooperEvent", fmt.Sprintf("Service Alive Count -> [%d]", application.serviceHub.Len()))
				}
			default:
				continue
			}
		}
		return 100, nil
	}, func(code int, data interface{}) {
		if code >= 200 {
			application.neuron.Brain.LogGenerater(model.LogError, tag, "LooperEvent", fmt.Sprintf("Looper Error -> %s", data))
		}
	}, application.neuron.Brain.Const.Interval.HZ1TimerInterval, application.looperStopC)

}

/* 退出事件 */
func exitEvent(exitSignal ...string) {
	if len(exitSignal) == 0 {
		application.neuron.Brain.LogGenerater(model.LogInfo, tag, "", tag+" Stopped Gracefully..")
	} else {
		application.neuron.Brain.LogGenerater(model.LogInfo, tag, "", tag+" Stopped by Signal -> "+fmt.Sprintf("%s", exitSignal)+"..")
	}
	// 停止looperEvent
	application.neuron.Brain.ClearInterval(application.looperStopC)
	// logs保存
	logger.Flush()
	os.Exit(0)
}

/* 监听系统退出信号 */
func sysExitSignalEvent() {
	go func() {
		exitSignal := <-application.sysExitSignalC
		exitEvent(exitSignal.String())
	}()
	// remove SIGPIPE(管道broken) SIGFPE(浮点运算错误) SIGTRAP(断点)
	signal.Notify(application.sysExitSignalC, syscall.SIGABRT, syscall.SIGALRM, syscall.SIGBUS,
		syscall.SIGHUP, syscall.SIGILL, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT,
		syscall.SIGSEGV, syscall.SIGTERM)
}

func init() {
	// 初始化神经元
	application.neuron = new(controller.NeuronS).Ontology()
	// 初始化应用
	application.sysExitSignalC = make(chan os.Signal)
	// 初始化服务容器
	application.serviceHub = new(model.QueueS).New(512)
	// 欢迎信息
	application.neuron.Brain.LogGenerater(model.LogInfo, tag, "", tag+" Starting..")
	// 监听系统信号
	sysExitSignalEvent()
}

/* ================================ SERVICE ================================ */

/* Service Process -> HTTP */
func serviceProcess() *model.ServiceS {
	// service init
	service := new(model.ServiceS)
	service.Tag = "HTTPServant"
	service.Alive = make(chan bool)
	// service run
	go application.neuron.Brain.SafeFunction(func() {
		/* Map of all Services */
		services := make(map[string]interface{})
		/* Construct Interface */
		mux := http.NewServeMux()
		/* Static Interface */
		mux.HandleFunc("/", application.neuron.Express.StaticHandler)

		application.neuron.Brain.LogGenerater(model.LogInfo, tag, service.Tag, "Preparing..")


		/* AD -> Proxy */
		proxyRoot := "/Proxy"
		proxy := new(controller.ProxyS).Ontology(application.neuron, mux, proxyRoot)
		services[proxyRoot] = proxy
		mux.HandleFunc(proxyRoot, func(res http.ResponseWriter, req *http.Request) {
			application.neuron.Express.ConstructService(proxy, proxyRoot, res, req)
		})

		application.neuron.Brain.LogGenerater(model.LogInfo, tag, service.Tag, "Prepared..")

		go protocalHTTP(service, mux)
		if application.neuron.Brain.Const.HTTPS.Open {
			go protocalTLS(service, mux)
		}
	})
	return service
}

func protocalHTTP(service *model.ServiceS, mux *http.ServeMux) {
	// HTTP Listen Port
	listenPort := strconv.Itoa(application.neuron.Brain.Const.HTTPServer.Port)
	listenAddr := application.neuron.Brain.Const.HTTPServer.Host + ":" + listenPort
	application.neuron.Brain.LogGenerater(model.LogInfo, tag, service.Tag, "Listening port -> "+listenPort)
	err := http.ListenAndServe(listenAddr, mux)
	if err != nil {
		application.neuron.Brain.MessageHandler(tag, "Protocal -> HTTP", 204, err)
		protocalTLS(service, mux)
		return
	}
	service.Alive <- false
}

func protocalTLS(service *model.ServiceS, mux *http.ServeMux) {
	// Listen Port
	listenPort := strconv.Itoa(application.neuron.Brain.Const.HTTPS.TLSPort)
	listenAddr := application.neuron.Brain.Const.HTTPServer.Host + ":" + listenPort
	application.neuron.Brain.LogGenerater(model.LogInfo, tag, service.Tag+"[TLS]", "Listening port -> "+listenPort)
	// Get Crt & Key
	crtPath := application.neuron.Brain.PathAbs(application.neuron.Brain.Const.HTTPS.TLSCertPath + ".crt")
	keyPath := application.neuron.Brain.PathAbs(application.neuron.Brain.Const.HTTPS.TLSCertPath + ".key")
	err := http.ListenAndServeTLS(listenAddr, crtPath, keyPath, mux)
	if err != nil {
		application.neuron.Brain.MessageHandler(tag, "Protocal -> TLS", 204, err)
	}
	service.Alive <- false
}

/* ================================ MAIN ================================ */

func main() {
	// push main process -> Neuron
	application.serviceHub.Push(serviceProcess())

	// pprof server
	if application.neuron.Brain.Const.RunEnv < 2 {
		go http.ListenAndServe(fmt.Sprintf("%s:%d", application.neuron.Brain.Const.HTTPServer.Host, application.neuron.Brain.Const.HTTPServer.Port+1), nil)
	}

	looperEvent()
}