package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"reflect"
	"sync"
	"time"
	"unsafe"
)

var dataMap map[string]*dataNode
var listMap map[string][]*dataNode
var dataMapLock,listMapLock *sync.RWMutex
func initGoKeyValue () {
	dataMap = make (map[string]*dataNode)
	listMap = make (map[string][]*dataNode)
	dataMapLock = &sync.RWMutex{}
	listMapLock = &sync.RWMutex{}
	timerChan = make (chan bool)
	timerMutex = &sync.Mutex{}
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
	t,_ := toTypePtr(value)
	dataMap [key] = &dataNode{
		Type:    t,
		Data: value,
	}
}
func GetString (key string) (string) {
	if _,ok := dataMap [key];!ok {
		return ""
	}
	return dataMap [key].Data.(string)
}
func GetInt (key string) (int) {
	if _,ok := dataMap [key];!ok {
		return 0
	}
	return dataMap [key].Data.(int)
}
func GetInt64 (key string) (int64) {
	if _,ok := dataMap [key];!ok {
		return 0
	}
	return dataMap [key].Data.(int64)
}
func GetBool (key string) (bool) {
	if _,ok := dataMap [key];!ok {
		return false
	}
	return dataMap [key].Data.(bool)
}
func InsertList (key string,v interface {}) {
	if _,ok := listMap [key];!ok {
		listMap [key] = make ([]*dataNode,0)
	}
	t,_ := toTypePtr(v)
	listMap[key] = append(listMap[key], &dataNode{
		Type:    t,
		Data: v,
	})
}
func PopList (key string) (v interface{}){
	if _,ok := listMap [key];!ok {
		return nil
	}
	if len (listMap [key]) == 0 {
		return nil
	}
	v = listMap [key][0].Data
	listMap [key] = listMap[key][1:]
	return
}
func ToString () (string,error) {
	out := make (map [string]interface{})
	func () {
		dataMapLock.RLock()
		defer dataMapLock.RUnlock()
		out ["data"] = dataMap
	} ()
	func () {
		listMapLock.RLock()
		defer listMapLock.RUnlock()
		out ["list"] = listMap
	} ()
	t,err := json.Marshal(out)
	if err != nil {
		return "",err
	}
	return string(t),nil
}
func SaveToFile (fileName string) error {
	if _,err := os.Stat(fileName);err == nil {
		os.Remove(fileName)
	}
	f,err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	out := make (map [string]interface{})
	func () {
		dataMapLock.RLock()
		defer dataMapLock.RUnlock()
		out ["data"] = dataMap
	} ()
	func () {
		listMapLock.RLock()
		defer listMapLock.RUnlock()
		out ["list"] = listMap
	} ()
	t,err := json.Marshal(out)
	if err != nil {
		return err
	}
	f.Write(t)
	return nil
}
func TopOfList (key string) interface{} {
	listMapLock.RLock()
	defer listMapLock.RUnlock()
	if _,ok := listMap [key];!ok {
		return nil
	}
	if len (listMap [key]) == 0 {
		return nil
	}
	ret := listMap [key][0].Data
	listMap [key] = listMap [key][1:]
	return ret
}
func LoadFromFile (fileName string) error {
	buf,err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	var data interface{}
	json.Unmarshal(buf,&data)
	t := data.(map [string]interface{})
	func () {
		dataMapLock.Lock()
		defer dataMapLock.Unlock()
		tmp := t ["data"].(map [string]interface{})
		for k,v := range tmp {
			t := v.(map [string]interface{})
			ty := (int)(math.Floor(t ["Type"].(float64)))
			if ty == TypeInt64 {
				t ["Data"] = (int64)(t ["Data"].(float64))
			}
			if ty == TypeInt {
				t ["Data"] = (int)(t ["Data"].(float64))
			}
			dataMap [k] = &dataNode{
				Type: (int)(math.Floor(t ["Type"].(float64))),
				Data: t ["Data"],
			}
		}
	} ()
	func () {
		listMapLock.Lock()
		defer listMapLock.Unlock()
		tmp := t ["list"].(map [string]interface{})
		for k,v := range tmp {
			tt := v.([]interface{})
			if _,ok := listMap [k];!ok {
				listMap [k] = make ([]*dataNode,0)
			}
			for _,v1 := range tt {
				t := v1.(map [string]interface{})
				ty := (int)(math.Floor(t ["Type"].(float64)))
				if ty == TypeInt64 {
					t ["Data"] = (int64)(t ["Data"].(float64))
				}
				if ty == TypeInt {
					t ["Data"] = (int)(t ["Data"].(float64))
				}
				listMap[k] = append(listMap[k], &dataNode{
					Type: ty,
					Data: t ["Data"],
				})
			}
		}
	} ()
	return nil
}
var timerChan chan bool
var timerRunning = false
var timerMutex *sync.Mutex
/*
设置定时保存 当second为-1时表示取消定时
initSave表示是否在启动定时保存的时候进行一次保存
 */
func InitSaveTimer (second int,initSave bool,saveFile string) error {
	timerMutex.Lock()
	defer timerMutex.Unlock()
	if second == -1 {
		timerChan <- true
		return nil
	}
	//开始运行
	if timerRunning {
		timerChan <- true
	}
	go func (saveFile string,second int) {
		if initSave {
			SaveToFile(saveFile)
		}
		for {
			select {
				case <- timerChan:	//要退出了
					fmt.Println("end2")
					return
				case <-time.After(time.Duration(second * 1000)):
					SaveToFile(saveFile)
			}
		}
	} (saveFile,second)
	return nil
}