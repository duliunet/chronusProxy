/*
===========================================================================
Just like Container
===========================================================================
*/

package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"modules/logs/logger"
)

/* ================================ DEFINE ================================ */
type NeuronS struct {
	Brain        *BrainS
	Express      *ExpressS
}

/* ================================ PRIVATE ================================ */

func (neuron *NeuronS) initLogger() {
	// 默认配置文件Base64
	configDefault, _ := base64.StdEncoding.DecodeString("PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0idXRmLTgiID8+PGxvZ3M+PGluZm8+PGNvbnNvbGUgb3V0cHV0PSJzdGRvdXQiIGZvcmVncm91bmQ9ImdyZWVuIiAvPjxidWZmZXIgc2l6ZT0iMTAiPjxyb3RhdGUgZGlyPSIvbG9nL0luZm8vIiBzaXplPSIzMk0iIC8+PC9idWZmZXI+PC9pbmZvPjxkZWJ1Zz48Y29uc29sZSBvdXRwdXQ9InN0ZG91dCIgZm9yZWdyb3VuZD0iY3lhbiIgLz48YnVmZmVyIHNpemU9IjEwIj48cm90YXRlIGRpcj0iL2xvZy9EZWJ1Zy8iIHNpemU9IjMyTSIgLz48L2J1ZmZlcj48L2RlYnVnPjx0cmFjZT48Y29uc29sZSBvdXRwdXQ9InN0ZG91dCIgZm9yZWdyb3VuZD0id2hpdGUiIC8+PGJ1ZmZlciBzaXplPSIxMCI+PHJvdGF0ZSBkaXI9Ii9sb2cvVHJhY2UvIiBzaXplPSIzMk0iIC8+PC9idWZmZXI+PC90cmFjZT48d2Fybj48Y29uc29sZSBvdXRwdXQ9InN0ZG91dCIgZm9yZWdyb3VuZD0ieWVsbG93IiAvPjxidWZmZXIgc2l6ZT0iNSI+PHJvdGF0ZSBkaXI9Ii9sb2cvV2Fybi8iIHNpemU9IjMyTSIgLz48L2J1ZmZlcj48L3dhcm4+PGVycm9yPjxjb25zb2xlIG91dHB1dD0ic3RkZXJyIiBmb3JlZ3JvdW5kPSJyZWQiIC8+PGJ1ZmZlciBzaXplPSIxIj48cm90YXRlIGRpcj0iL2xvZy9FcnJvci8iIHNpemU9IjMyTSIgLz48L2J1ZmZlcj48L2Vycm9yPjxjcml0aWNhbD48Y29uc29sZSBvdXRwdXQ9InN0ZGVyciIgZm9yZWdyb3VuZD0ibWFnZW50YSIgLz48YnVmZmVyIHNpemU9IjEiPjxyb3RhdGUgZGlyPSIvbG9nL0NyaXRpY2FsLyIgc2l6ZT0iMzJNIiAvPjwvYnVmZmVyPjwvY3JpdGljYWw+PC9sb2dzPg==")
	// 配置文件路径
	configPath := neuron.Brain.PathAbs("/log/config.xml")
	// 读取配置文件
	code, data := neuron.Brain.FileReader(configPath)
	if code == 100 {
		err := logger.InitFromXMLString(string(data.([]byte)))
		if err != nil {
			fmt.Println("[NeuronInit]LoggerLoad => Failed: " + err.Error())
		}
	} else {
		// 读取文件失败
		fmt.Println("[NeuronInit]LoggerLoad => Failed: " + data.(error).Error())
		// 写入默认配置
		fmt.Println("[NeuronInit]Logger => Creating Default Config File")
		code, data := neuron.Brain.FileWriter(configPath, []byte(configDefault))
		if code == 100 {
			err := logger.InitFromXMLFile(configPath)
			if err != nil {
				fmt.Println("[NeuronInit]LoggerWrite => Failed: " + err.Error())
			}
		} else {
			fmt.Printf("[NeuronInit]LoggerWrite => Failed: %+v", data)
		}
	}
}

func (neuron *NeuronS) initConfig() {
	// 默认配置文件
	configDefault := neuron.Brain.Const
	// 配置文件路径
	configPath := neuron.Brain.PathAbs("/config.json")
	// 读取配置文件
	code, data := neuron.Brain.FileReader(configPath)
	if code == 100 {
		// 读取文件成功
		err := json.Unmarshal(data.([]byte), &neuron.Brain.Const)
		if err != nil {
			fmt.Println("[NeuronInit]ConfigLoad => Failed: " + err.Error())
		}
	} else {
		// 读取文件失败
		fmt.Println("[NeuronInit]ConfigLoad => Failed: " + data.(error).Error())

		// 默认生产环境客户端不自动生成config
		if neuron.Brain.Const.RunEnv == 2 {
			fmt.Println("[NeuronInit]ConfigLoad => Failed: You Need file [config.json]")
			return
		}

		// 写入默认配置
		json, err := json.MarshalIndent(configDefault, "", "    ")
		if err != nil {
			fmt.Println("[NeuronInit]ConfigWrite => Failed: " + err.Error())
		}
		code, data := neuron.Brain.FileWriter(configPath, json)
		if code != 100 {
			fmt.Printf("[NeuronInit]ConfigWrite => Failed: %+v", data)
		} else {
			neuron.initConfig()
		}
	}
}

func (neuron *NeuronS) initStatic() {
	// 默认配置文件Base64
	staticDefault, _ := base64.StdEncoding.DecodeString("PCFET0NUWVBFIGh0bWw+CjxodG1sPgo8aGVhZD4KICAgIDx0aXRsZT5XZWxjb21lIHRvIENocm9udXMgRXhwcmVzczwvdGl0bGU+CjwvaGVhZD4KPGJvZHk+CjxoMT5XZWxjb21lIHRvIENocm9udXMgRXhwcmVzczwvaDE+CjxkaXYgc3R5bGU9ImJhY2tncm91bmQ6IHVybCh0ZW1wL2NvZGUuanBnKSBuby1yZXBlYXQ7cGFkZGluZy10b3A6IDIwJSI+PC9kaXY+CjwvYm9keT4KPC9odG1sPg==")
	// 配置文件路径
	staticPath := neuron.Brain.PathAbs("/static/index.html")
	uploadPath := neuron.Brain.PathAbs(neuron.Brain.Const.HTTPServer.UploadPath)
	// 读取配置文件
	code, data := neuron.Brain.FileReader(staticPath)
	if code != 100 {
		// 读取文件失败
		fmt.Println("[NeuronInit]InitStatic => Failed: " + data.(error).Error())
		// 写入默认配置
		fmt.Println("[NeuronInit]InitStatic => Creating Default Index File")
		code, data := neuron.Brain.FileWriter(staticPath, []byte(staticDefault))
		if code != 100 {
			fmt.Println("[NeuronInit]InitStatic Error => ", data.(error).Error())
			return
		}
		code, data = neuron.Brain.PathCreate(uploadPath)
		if code != 100 {
			fmt.Println("[NeuronInit]InitStatic Error => ", data.(error).Error())
			return
		}
		fmt.Println("[NeuronInit]InitStatic => Created")
	}
}

/* ================================ PUBLIC ================================ */
/* 构造本体 */
func (neuron *NeuronS) Ontology() *NeuronS {
	// Brain
	neuron.Brain = new(BrainS).Ontology()
	// Initialize
	neuron.initLogger()
	neuron.initConfig()
	neuron.initStatic()
	// Express
	neuron.Express = new(ExpressS).Ontology(neuron)
	return neuron
}

/* 重新读取配置文件 */
func (neuron *NeuronS) InitConfig() {
	neuron.initConfig()
}