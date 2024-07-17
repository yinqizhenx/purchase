package utils

import (
	"reflect"
	"runtime"
	"strings"
)

func GetMethodName(i interface{}) string {
	// 使用反射获取值对象
	val := reflect.ValueOf(i)

	// 检查是否是函数类型
	if val.Kind() != reflect.Func {
		return ""
	}
	// 获取包名函数名
	fullName := runtime.FuncForPC(val.Pointer()).Name()
	// 获取函数名
	parts := strings.Split(fullName, ".")
	return parts[len(parts)-1]
}
