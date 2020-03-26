package time

import (
	lua "github.com/yuin/gopher-lua"
	"time"
)

func Now(L *lua.LState) int {
	L.Push(lua.LNumber(time.Now().Nanosecond() / 1000000))
	return 1
}
