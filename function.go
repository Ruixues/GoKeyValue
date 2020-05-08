package main

import (
	"reflect"
	"unsafe"
)

var dataMap map[string]*dataNode
var listMap map[string][]*dataNode
func initGoKeyValue () {
	dataMap = make (map[string]*dataNode)
	listMap = make (map[string][]*dataNode)
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
func toTypePtr (value interface{}) (t int,ptr unsafe.Pointer) {
	t = interfaceToType(value)
	if t == UnsupportedType {
		panic("不支持的类型:" + reflect.TypeOf(value).String())
		return
	}
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
	return
}
func SetKey (key string,value interface {}) {
	t,ptr := toTypePtr(value)
	dataMap [key] = &dataNode{
		Type:    t,
		Pointer: ptr,
	}
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
func InsertList (key string,v interface {}) {
	if _,ok := listMap [key];!ok {
		listMap [key] = make ([]*dataNode,0)
	}
	t,ptr := toTypePtr(v)
	listMap[key] = append(listMap[key], &dataNode{
		Type:    t,
		Pointer: ptr,
	})
}
func PopList (key string) (v interface{}){
	if _,ok := listMap [key];!ok {
		return nil
	}
	if len (listMap [key]) == 0 {
		return nil
	}
	v = *(listMap [key][0].Pointer)
	listMap [key] = listMap[key][1:]
	return
}