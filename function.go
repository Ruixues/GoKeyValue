package main

import "unsafe"

var dataMap map[string]*dataNode
func initGoKeyValue () {
	dataMap = make (map[string]*dataNode)
	
}
func interfaceToType (data interface{}) int {
	switch data.(type) {
	case string:
		return TypeString
	case int:
		return TypeInt
	case int64:
		return TypeInt64
	case bool:
		return TypeBool
	}
	return UnsupportedType
}
func SetKey (key string,value interface {}) {
	t := interfaceToType(value)
	var ptr unsafe.Pointer
	switch value.(type) {
	case string:
		t := value.(string)
		ptr = unsafe.Pointer(&t)
		break
	case int:
		t := value.(int)
		ptr = unsafe.Pointer(&t)
		break
	case bool:
		t := value.(bool)
		ptr = unsafe.Pointer(&t)
		break
	case int64:
		t := value.(int64)
		ptr = unsafe.Pointer(&t)
		break
	}
	dataMap [key].Type = t
	dataMap [key].Pointer = ptr
}
func GetString (key string) (string) {
	if _,ok := dataMap [key];!ok {
		return ""
	}
	return *(*string)(dataMap [key].Pointer)
}
func GetInt (key string) (int) {
	if _,ok := dataMap [key];!ok {
		return 0
	}
	return *(*int)(dataMap [key].Pointer)
}
func GetInt64 (key string) (int64) {
	if _,ok := dataMap [key];!ok {
		return 0
	}
	return *(*int64)(dataMap [key].Pointer)
}
func GetBool (key string) (bool) {
	if _,ok := dataMap [key];!ok {
		return false
	}
	return *(*bool)(dataMap [key].Pointer)
}

