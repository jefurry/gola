// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pm

import (
	"github.com/yuin/gopher-lua"
	"time"
)

type (
	lState struct {
		L                  *lua.LState
		requestedNum       int
		startTime          int64
		maxRequest         int
		idleTimeout        string
		idleTimeoutSeconds int
	}
)

func newLState(maxRequest int, idleTimeout string, seconds int, whenNew newFunc) (*lState, error) {
	l := lua.NewState()
	if whenNew != nil {
		if err := whenNew(l); err != nil {
			l.Close()

			return nil, err
		}
	}

	ls := &lState{
		L:                  l,
		requestedNum:       0,
		startTime:          time.Now().Unix(),
		maxRequest:         maxRequest,
		idleTimeout:        idleTimeout,
		idleTimeoutSeconds: seconds,
	}

	return ls, nil
}

func (ls *lState) incRequestNum() {
	ls.requestedNum += 1
}

func (ls *lState) isExpire() bool {
	if ls.maxRequest > 0 {
		if ls.requestedNum >= ls.maxRequest {
			return true
		}

		return false
	}

	if ls.idleTimeoutSeconds > 0 {
		if (time.Now().Unix() - ls.startTime) > int64(ls.idleTimeoutSeconds) {
			return true
		}
	}

	return false
}

func (ls *lState) close() {
	ls.L.Close()
}
