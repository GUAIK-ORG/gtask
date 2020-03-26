package request

import (
	"github.com/golang/glog"
	lua "github.com/yuin/gopher-lua"
	"io/ioutil"
	"net/http"
	"strings"
)

func Get(L *lua.LState) int {
	url := L.ToString(1)
	headers := L.ToTable(2)
	client := &http.Client{}
	req, err := http.NewRequest("GET",url,nil)
	if err != nil {
		glog.Error(err)
		L.Push(lua.LString(""))
		L.Push(lua.LFalse)
	}
	headers.ForEach(func(key,value lua.LValue){
		req.Header.Set(key.String(),value.String())
	})
	resp, err := client.Do(req)
	if err != nil {
		glog.Error(err)
		L.Push(lua.LString(""))
		L.Push(lua.LFalse)
	}
	defer resp.Body.Close()
	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		L.Push(lua.LString(""))
		L.Push(lua.LFalse)
	}
	L.Push(lua.LString(string(body)))
	L.Push(lua.LTrue)
	return 2
}

func Post(L *lua.LState) int {
	url := L.ToString(1)
	headers := L.ToTable(2)
	data := L.ToString(3)
	client := &http.Client{}
	req, err := http.NewRequest("POST",url,strings.NewReader(data))
	if err != nil {
		glog.Error(err)
		L.Push(lua.LString(""))
		L.Push(lua.LFalse)
	}
	headers.ForEach(func(key,value lua.LValue){
		req.Header.Set(key.String(),value.String())
	})
	resp, err := client.Do(req)
	if err != nil {
		glog.Error(err)
		L.Push(lua.LString(""))
		L.Push(lua.LFalse)
	}
	defer resp.Body.Close()
	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		L.Push(lua.LString(""))
		L.Push(lua.LFalse)
	}
	L.Push(lua.LString(string(body)))
	L.Push(lua.LTrue)
	return 2
}
