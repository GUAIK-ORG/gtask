package json

import (
	"encoding/json"
	"github.com/golang/glog"
	lua "github.com/yuin/gopher-lua"
)

// 检查Table是否为List
func checkList(value lua.LValue) (b bool) {
	if value.Type().String() == "table" {
		b = true
		value.(*lua.LTable).ForEach(func(k, v lua.LValue) {
			if k.Type().String() != "number" {
				b = false
				return
			}
		})
	}
	return
}

func marshal(data lua.LValue) interface{} {
	switch data.Type() {
	case lua.LTTable:
		if checkList(data) {
			jdata := make([]interface{}, 0)
			data.(*lua.LTable).ForEach(func(key, value lua.LValue) {
				jdata = append(jdata, marshal(value))
			})
			return jdata
		} else {
			jdata := map[string]interface{}{}
			data.(*lua.LTable).ForEach(func(key, value lua.LValue) {
				jdata[key.String()] = marshal(value)
			})
			return jdata
		}
	case lua.LTNumber:
		return float64(data.(lua.LNumber))
	case lua.LTString:
		return string(data.(lua.LString))
	case lua.LTBool:
		return bool(data.(lua.LBool))
	}
	return nil
}

func JsonMarshal(L *lua.LState) int {
	data := L.ToTable(1)
	str, err := json.Marshal(marshal(data))
	if err != nil {
		glog.Error(err)
	}
	L.Push(lua.LString(str))
	return 1
}

func unmarshal(L *lua.LState, data interface{}) lua.LValue {
	switch data.(type) {
	case map[string]interface{}:
		tb := L.NewTable()
		for k, v := range data.(map[string]interface{}) {
			tb.RawSet(lua.LString(k), unmarshal(L, v))
		}
		return tb
	case []interface{}:
		tb := L.NewTable()
		for i, v := range data.([]interface{}) {
			tb.Insert(i+1, unmarshal(L, v))
		}
		return tb
	case float64:
		return lua.LNumber(data.(float64))
	case string:
		return lua.LString(data.(string))
	case bool:
		return lua.LBool(data.(bool))
	}
	return lua.LNil
}

func JsonUnMarshal(L *lua.LState) int {
	str := L.ToString(1)
	jdata := map[string]interface{}{}
	err := json.Unmarshal([]byte(str), &jdata)
	if err != nil {
		glog.Error(err)
	}
	L.Push(unmarshal(L, jdata))
	return 1
}
