package gluasocket_socketcore

import (
	"time"

	"github.com/yuin/gopher-lua"
)

const (
	MASTER_TYPENAME = "tcp{master}"
)

type Master struct {
	Timeout time.Duration
}

var masterMethods = map[string]lua.LGFunction{
	"bind":       masterBindMethod,
	"close":      masterCloseMethod,
	"connect":    masterConnectMethod,
	"listen":     masterListenMethod,
	"settimeout": masterSetTimeoutMethod,
}

// ----------------------------------------------------------------------------

func checkMaster(L *lua.LState) *Master {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Master); ok {
		return v
	}
	L.ArgError(1, "master expected")
	return nil
}
