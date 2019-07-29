package util

import (
	"encoding/json"
	"runtime/debug"
)

//Obj2Json 对象转为 json 编码
func Obj2Json(obj interface{}) string {
	bj, _ := json.Marshal(obj)
	return string(bj)
}

func Obj2JsonIndent(obj interface{}) string {
	bj, _ := json.MarshalIndent(obj, "    ", "    ")
	return string(bj)
}

func Err2JsonObj(err error) map[string]interface{} {
	return map[string]interface{}{
		"code":     500,
		"err_code": 500,
		//"ver":      constant.Version,
		"error": err.Error(),
		"stack": string(debug.Stack()),
		"data":  nil,
	}
}

func ErrC2JsonObj(c int, err error) map[string]interface{} {
	return map[string]interface{}{
		"code":     c,
		"err_code": c,
		//"ver":      constant.Version,
		"error": err.Error(),
		"stack": string(debug.Stack()),
		"data":  nil,
	}
}
