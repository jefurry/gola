package gluasocket_socketcore

import (
	"github.com/yuin/gopher-lua"
)

func masterBindMethod(L *lua.LState) int {
	L.RaiseError("master:bind() not implemented yet")
	return 0
}
