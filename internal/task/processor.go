package task

import (
	"github.com/golang/glog"
	lua "github.com/yuin/gopher-lua"
	luaCrypto "gtask/lua_module/crypto"
	luaJson "gtask/lua_module/json"
	luaRequest "gtask/lua_module/request"
	luaTime "gtask/lua_module/time"
)

type TimeoutFunc func(string)

type Processor struct {
	trigger int64
	count   int64
	bReset  bool // 是否会被系统重置
	bLoop   bool // 是否循环(自动重置count)
	bExit   bool // 是否退出处理器
	code    string
	env     *lua.LState
}

func NewProcessor(code string, trigger int64, bReset, bLoop, bExit bool) *Processor {
	env := lua.NewState()
	env.SetGlobal("md5", env.NewFunction(luaCrypto.Str2MD5))
	env.SetGlobal("base64UrlSafe", env.NewFunction(luaCrypto.UrlSafeBase64))
	env.SetGlobal("base64", env.NewFunction(luaCrypto.Base64))
	env.SetGlobal("hmac", env.NewFunction(luaCrypto.Str2HMACWithSHA1))
	env.SetGlobal("sha1", env.NewFunction(luaCrypto.Str2SHA1))
	env.SetGlobal("httpGet", env.NewFunction(luaRequest.Get))
	env.SetGlobal("httpPost", env.NewFunction(luaRequest.Post))
	env.SetGlobal("now", env.NewFunction(luaTime.Now))
	env.SetGlobal("jsonMarshal", env.NewFunction(luaJson.JsonMarshal))
	env.SetGlobal("jsonUnMarshal", env.NewFunction(luaJson.JsonUnMarshal))

	err := env.DoString(code)
	if err != nil {
		glog.Error(err)
		return nil
	}
	return &Processor{
		count:   0,
		trigger: trigger,
		bReset:  bReset,
		bLoop:   bLoop,
		bExit:   bExit,
		code:    code,
		env:     env,
	}
}

func (p *Processor) Run(key string, count int64) (state bool) {
	err := p.env.CallByParam(lua.P{
		Fn:      p.env.GetGlobal("processor"),
		NRet:    1,
		Protect: true,
		Handler: nil,
	}, lua.LString(key), lua.LNumber(count))
	if err != nil {
		glog.Error(err)
	}
	ret := p.env.Get(-1)
	if sta, ok := ret.(lua.LBool); ok {
		if sta == lua.LTrue {
			state = true
		}
	}
	return
}

func (p *Processor) Release() {
	p.env.Close()
}
