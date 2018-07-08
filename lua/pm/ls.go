// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pm

import (
	"context"
	"github.com/yuin/gopher-lua"
	"time"
)

type (
	lState struct {
		L                  *lua.LState
		cancel             context.CancelFunc
		requestedNum       int
		startTime          int64
		maxRequest         int
		idleTimeout        string
		idleTimeoutSeconds int
		serving            bool
		closed             bool
	}
)

func newLState(ctx context.Context, maxRequest int, idleTimeout string, seconds int, options lua.Options, whenNew NewFunc) (*lState, error) {
	l := lua.NewState(options)
	if whenNew != nil {
		if err := whenNew(l); err != nil {
			l.Close()

			return nil, err
		}
	}

	lctx, cancel := context.WithCancel(ctx)
	l.SetContext(lctx)

	ls := &lState{
		L:                  l,
		cancel:             cancel,
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

func (ls *lState) mustTerminate() bool {
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

func (ls *lState) setServing(serving bool) {
	ls.serving = serving
}

func (ls *lState) close() {
	ls.closed = true
	ls.setServing(false)
	ls.cancel()
	ls.L.Close()
}
