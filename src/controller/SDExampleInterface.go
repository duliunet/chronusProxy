/*
===========================================================================

===========================================================================
*/

package controller

import (
	"model"
	"net/http"
)

/* ================================ DEFINE ================================ */
type ExampleS struct {
	Const struct {
		tag  string
		root string
	}
	isStarted bool
	neuron    *NeuronS
	mux       *http.ServeMux
}

/* ================================ PRIVATE ================================ */
/* 注册服务 */
func (mExample *ExampleS) main() {
}

/* ================================ INTERFACE ================================ */

/* 某服务开关 */
func (mExample *ExampleS) ExampleInterface() {
	mExample.mux.HandleFunc(mExample.Const.root+"/Example", func(res http.ResponseWriter, req *http.Request) {
		mExample.neuron.Express.ConstructInterface(res, req, mExample.isStarted, func() {
			query := mExample.neuron.Express.Req2Query(req)
			if mExample.neuron.Brain.CheckIsNull(query) {
				mExample.neuron.Express.CodeResponse(res, 207, "Lack of Param", "ExampleInterface")
				return
			}
			if !mExample.neuron.Brain.CheckIsNull(query["chronus"]) {
				mExample.neuron.Express.CodeResponse(res, 100, "Hello Chronus Express!", "ExampleInterface")
			} else {
				mExample.neuron.Express.CodeResponse(res, 100, "Hello World!", "ExampleInterface")
			}
		}, func(err interface{}) {
			mExample.neuron.Express.CodeResponse(res, 204, err, "ExampleInterface[ConstructInterface]")
		})
	})
}

/* ================================ PROCESS ================================ */

/* ================================ SQL PROCESS ================================ */

/* ================================ TOOL ================================ */

/* ================================ LOOPER ================================ */

/* ================================ SERVICE ================================ */

/* 构造服务 */
func (mExample *ExampleS) service() {
}

/* 析构服务 */
func (mExample *ExampleS) serviceKiller() {
}

/* ================================ PUBLIC ================================ */

/* 构造本体 */
func (mExample *ExampleS) Ontology(neuron *NeuronS, mux *http.ServeMux, root string) *ExampleS {
	mExample.neuron = neuron
	mExample.mux = mux
	mExample.Const.tag = root[1:]
	mExample.Const.root = root
	mExample.isStarted = neuron.Brain.Const.AutorunConfig.SDExampleInterface
	if mExample.isStarted {
		mExample.neuron.Brain.SafeFunction(mExample.main, func(err interface{}) {
			mExample.neuron.Brain.MessageHandler(mExample.Const.tag, "Main", 204, err)
		})
		mExample.StartService()
	} else {
		mExample.StopService()
	}
	return mExample
}

/* 返回开关量 */
func (mExample *ExampleS) IsStarted() bool {
	return mExample.isStarted
}

/* 启动服务 */
func (mExample *ExampleS) StartService() {
	mExample.isStarted = true
	go mExample.neuron.Brain.SafeFunction(mExample.service, func(err interface{}) {
		mExample.neuron.Brain.MessageHandler(mExample.Const.tag, "StartService", 204, err)
	})
}

/* 停止服务 */
func (mExample *ExampleS) StopService() {
	mExample.isStarted = false
	go mExample.neuron.Brain.SafeFunction(mExample.serviceKiller, func(err interface{}) {
		mExample.neuron.Brain.MessageHandler(mExample.Const.tag, "StopService", 204, err)
	})
}

/* 打印信息 */
func (mExample *ExampleS) Log(title string, content interface{}) {
	mExample.neuron.Brain.LogGenerater(model.LogInfo, mExample.Const.tag, title, content)
}

/* ================================ EVAL INTERFACE ================================ */
