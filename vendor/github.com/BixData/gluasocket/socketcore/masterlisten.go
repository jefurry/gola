package gluasocket_socketcore

import (
	"github.com/yuin/gopher-lua"
)

func masterListenMethod(L *lua.LState) int {
	L.RaiseError("master:listen() not implemented yet")
	return 0
}
