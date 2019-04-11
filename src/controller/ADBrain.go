package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"model"
	"modules/logs/logger"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
	"time"
)

/* ================================ DEFINE ================================ */

type BrainS struct {
	tag       string
	Const     model.Const
}

/* 构造本体 */
func (brain *BrainS) Ontology() *BrainS {
	brain.tag = "Brain"
	brain.Const = brain.Const.Ontology()
	return brain
}

/* 判断是否为空 */
func (brain *BrainS) CheckIsNull(i interface{}) bool {
	if brain.Const.RunEnv != 0 {
		defer func() {
			if err := recover(); err != nil {
				brain.MessageHandler(brain.tag, "CheckIsNull", 204, err)
			}
		}()
	}
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Struct, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
		return false
	case reflect.Slice, reflect.Array:
		if v.IsNil() {
			return true
		} else {
			if v.Len() == 0 {
				return true
			} else {
				return false
			}
		}
	case reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface:
		if v.IsNil() {
			return true
		} else {
			return false
		}
	case reflect.Chan:
		if v.IsNil() {
			return true
		}
		switch i.(type) {
		case chan bool:
			select {
			case _, ok := <-i.(chan bool):
				return !ok
			default:
			}
		case chan int:
			select {
			case _, ok := <-i.(chan int):
				return !ok
			default:
			}
		case chan interface{}:
			select {
			case _, ok := <-i.(chan interface{}):
				return !ok
			default:
			}
		}
	}
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}

/* 结构化日志记录 */
func (brain *BrainS) LogGenerater(logtype int, model string, function string, content interface{}) {
	if brain.CheckIsNull(content) {
		content = ""
	}
	timenow := "[" + time.Now().Format("2006-01-02 15:04:05") + "] "
	if function != "" {
		function = "_" + function
	}

	switch logtype {
	case 0:
		logger.Info(timenow + "[Info] " + brain.Const.NeuronId + "[" + brain.Const.Version + "]" + " - [" + model + function + "] => " + fmt.Sprintf("%+v", content))
	case 1:
		logger.Debug(timenow + "[Debug] " + brain.Const.NeuronId + "[" + brain.Const.Version + "]" + " - [" + model + function + "] => " + fmt.Sprintf("%+v", content))
	case 2:
		logger.Trace(timenow + "[Trace] " + brain.Const.NeuronId + "[" + brain.Const.Version + "]" + " - [" + model + function + "] => " + fmt.Sprintf("%+v", content))
	case 3:
		logger.Warn(timenow + "[Warn] " + brain.Const.NeuronId + "[" + brain.Const.Version + "]" + " - [" + model + function + "] => " + fmt.Sprintf("%+v", content))
	case 4:
		logger.Error(timenow + "[Error] " + brain.Const.NeuronId + "[" + brain.Const.Version + "]" + " - [" + model + function + "] => " + fmt.Sprintf("%+v", content))
	case 5:
		logger.Critical(timenow + "[Critical] " + brain.Const.NeuronId + "[" + brain.Const.Version + "]" + " - [" + model + function + "] => " + fmt.Sprintf("%+v", content))
	}
}

/* 信息处理 */
func (brain *BrainS) MessageHandler(tag string, function string, code int, data interface{}) model.MessageS {
	// Code & Data Output
	logtype := model.LogInfo
	if code > 100 && code < 200 {
		logtype = model.LogWarn
	} else if code >= 200 {
		logtype = model.LogError
	}
	message := brain.Const.ErrorCode[code]
	if brain.CheckIsNull(message) {
		message = brain.Const.ErrorCode[200]
	}
	msgs := model.MessageS{
		Code:    code,
		Message: message,
		Data:    data,
	}
	brain.LogGenerater(logtype, tag, function, msgs)
	return msgs
}

/* 同步TryCatch实现 */
func (brain *BrainS) SafeFunction(next func(), callback ...func(err interface{})) {
	if brain.Const.RunEnv > 0 {
		defer func() {
			if err := recover(); err != nil {
				errR := err
				// 捕获堆栈信息
				if brain.Const.RunEnv < 2 {
					var buf [102400]byte
					n := runtime.Stack(buf[:], false)
					errR = fmt.Sprintf("[%v]\r\n%v", err, string(buf[:n]))
				}
				brain.LogGenerater(model.LogCritical, brain.tag, "SafeFunction", fmt.Sprintf("{%v -> %v}", brain.GetFuncName(next), errR))
				for _, v := range callback {
					v(errR)
				}
			}
		}()
	}
	next()
}

/* 构造可用路径 */
func (brain *BrainS) PathAbs(dirPath string) string {
	dirPath = strings.Replace(dirPath, "../", "", -1)
	return path.Dir(os.Args[0]) + dirPath
}

/* 文件读 */
func (brain *BrainS) FileReader(filePath string) (int, interface{}) {
	var codeR int
	var dataR interface{}
	brain.SafeFunction(func() {
		dirPath := path.Dir(filePath)
		if brain.PathExists(dirPath) {
			fileBuffer, err := ioutil.ReadFile(filePath)
			if err != nil {
				codeR = 205
				dataR = err
			} else {
				codeR = 100
				dataR = fileBuffer
			}
		} else {
			codeR = 205
			dataR = errors.New("[FileReader] file path Error")
		}
	})
	return codeR, dataR
}

/* 文件写 */
func (brain *BrainS) FileWriter(filePath string, data []byte) (int, interface{}) {
	var codeR int
	var dataR interface{}
	brain.SafeFunction(func() {
		dirPath := path.Dir(filePath)
		if code, err := brain.PathCreate(dirPath); code == 100 {
			err := ioutil.WriteFile(filePath, data, os.FileMode(brain.Const.File.Chmod))
			if err != nil {
				codeR = 205
				dataR = err
			} else {
				codeR = 100
				dataR = nil
			}
		} else {
			codeR = code
			dataR = err
		}
	})
	return codeR, dataR
}

/* 反射获取方法名 */
func (brain *BrainS) GetFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

/* 查询文件夹 */
func (brain *BrainS) PathExists(dirPath string) bool {
	_, err := os.Stat(dirPath)
	if err != nil {
		return false
	} else {
		return true
	}
}

/* 永久循环 */
func (brain *BrainS) SetInterval(next func() (int, interface{}), callback func(code int, data interface{}), interval int, stopC chan bool, containers ...map[string]interface{}) {
	nextName := brain.GetFuncName(next)
	if brain.CheckIsNull(stopC) {
		brain.LogGenerater(model.LogError, brain.tag, "SetInterval", fmt.Sprintf("%s -> Lack of Stop Channel", nextName))
		return
	}
	defer close(stopC)
	if brain.Const.RunEnv < 2 {
		brain.LogGenerater(model.LogDebug, brain.tag, "SetInterval", nextName)
	}
	endC := make(chan map[int]interface{})
	msgC := make(chan map[int]interface{})
	defer close(endC)
	defer close(msgC)
	// Timer Init
	duration := time.Duration(interval) * time.Millisecond
	mTimer := time.NewTimer(duration)
	defer mTimer.Stop()
	// Timer Runnable
	for {
		select {
		// Exit Handler
		case data := <-stopC:
			if data {
				callback(103, nextName)
				return
			}
		default:
			go brain.SafeFunction(func() {
				code, data := next()
				if code == 100 {
					msgC <- map[int]interface{}{code: data}
				} else {
					endC <- map[int]interface{}{code: data}
				}
			}, func(err interface{}) {
				endC <- map[int]interface{}{204: err}
			})
		}
		select {
		case data := <-endC:
			// Error Handler
			for k, v := range data {
				callback(k, v)
			}
			return
		case <-mTimer.C:
			// Message Handler
			brain.SafeFunction(func() {
				data := <-msgC
				for k, v := range data {
					callback(k, v)
				}
			})
			// Reset duration
			if len(containers) > 0 {
				count := len(containers[0])
				if count == 0 {
					count = 1
				}
				mduraion := time.Duration(interval/count) * time.Millisecond
				mTimer.Reset(mduraion)
			} else {
				mTimer.Reset(duration)
			}
		}
	}
}

/* 结束永久循环 */
func (brain *BrainS) ClearInterval(stopC chan bool) {
	brain.SafeFunction(func() {
		if !brain.CheckIsNull(stopC) {
			stopC <- true
		}
	})
}

/* 创建文件夹 */
func (brain *BrainS) PathCreate(dirPath string) (int, interface{}) {
	if brain.PathExists(dirPath) {
		return 100, nil
	} else {
		err := os.MkdirAll(dirPath, os.FileMode(brain.Const.File.Chmod))
		if err != nil {
			return 205, err
		} else {
			return 100, nil
		}
	}
}

/* Json -> String2Object */
func (brain *BrainS) JsonDecoder(data []byte, structS ...interface{}) interface{} {
	// Try Decode Struct
	if len(structS) > 0 {
		err := json.Unmarshal(data, structS[0])
		if err == nil {
			return structS[0]
		} else {
			brain.MessageHandler(brain.tag, fmt.Sprintf("JsonDecoder[Struct] -> %s", structS[0]), 202, err)
		}
	}
	// Try Decode MapObject
	var dataObjMap map[string]interface{}
	err := json.Unmarshal(data, &dataObjMap)
	if err == nil {
		return dataObjMap
	} else {
		brain.MessageHandler(brain.tag, fmt.Sprintf("JsonDecoder[Map] -> %s", data), 202, err)
	}
	// Try Decode ArrayObjectc
	var dataObjArray []interface{}
	err = json.Unmarshal(data, &dataObjArray)
	if err == nil {
		return dataObjArray
	} else {
		brain.MessageHandler(brain.tag, fmt.Sprintf("JsonDecoder[Array] -> %s", data), 202, err)
	}
	return nil
}

/* Json -> Object2String */
func (brain *BrainS) JsonEncoder(data interface{}, indent ...bool) []byte {
	switcher := false
	if !brain.CheckIsNull(indent) {
		switcher = indent[0]
	}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(data)
	if err != nil {
		brain.MessageHandler(brain.tag, "JsonEncoder", 202, err)
		return nil
	}
	if switcher {
		var dstBuf bytes.Buffer
		err := json.Indent(&dstBuf, buf.Bytes(), "", "    ")
		if err != nil {
			brain.MessageHandler(brain.tag, "JsonEncoder", 202, err)
			return nil
		}
		return dstBuf.Bytes()
	}
	return buf.Bytes()
}