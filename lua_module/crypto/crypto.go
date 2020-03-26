package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	lua "github.com/yuin/gopher-lua"
)

func Str2MD5(L *lua.LState) int {
	str := L.ToString(1)
	h := md5.New()
	h.Write([]byte(str))
	L.Push(lua.LString(hex.EncodeToString(h.Sum(nil))))
	return 1
}


func Str2SHA1(L *lua.LState) int {
	str := L.ToString(1)
	h := sha1.New()
	h.Write([]byte(str))
	L.Push(lua.LString(hex.EncodeToString(h.Sum(nil))))
	return 1
}

func Str2HMACWithSHA1(L *lua.LState) int {
	key := L.ToString(1)
	str := L.ToString(2)
	mac := hmac.New(sha1.New,[]byte(key))
	mac.Write([]byte(str))
	L.Push(lua.LString(hex.EncodeToString(mac.Sum(nil))))
	return 1
}

func UrlSafeBase64(L *lua.LState) int {
	str := L.ToString(1)
	L.Push(lua.LString(base64.URLEncoding.EncodeToString([]byte(str))))
	return 1
}

func Base64(L *lua.LState) int {
	str := L.ToString(1)
	L.Push(lua.LString(base64.StdEncoding.EncodeToString([]byte(str))))
	return 1
}