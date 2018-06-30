// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package pm implements LuaState Manages. It Wraps gopher-lua as runner.
package pm

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/yuin/gopher-lua"
	"sync"
)

const (
	idleTimeoutDay    = "d"
	idleTimeoutHour   = "h"
	idleTimeoutMinute = "m"
	idleTimeoutSecond = "s"
)

const (
	defaultIdleTimeoutNum = 1
)

const (
	DefaultMaxNum      = 0
	DefaultStartNum    = 1
	DefaultMaxRequest  = 0
	DefaultIdleTimeout = string(defaultIdleTimeoutNum) + idleTimeoutHour
)

var (
	ErrIdleTimeoutFormat = errors.New("not a valid idle timeout format")
	ErrLuaStatePoolFull  = errors.New("lua state pool fulled")
)

type (
	newFunc func(l *lua.LState) error
)

type (
	// lua state pool manager.
	poolManager struct {
		// The maximum of processes lua state will start. This has been designed to control
		// the global number of lua state when using a lot of pools.
		// Use it with caution.
		// Note: A value of 0 indicates no limit.
		// Default Value: 0.
		maxNum int

		// The number of lua state created on startup.
		// Note: The number must be smaller than or equal to `maxNum`.
		// Defaulut Value: (maxNum / 2) + 1.
		startNum int

		// The number of requests each lua state should execute before respawning.
		// This can be useful to work around memory leaks in 3rd party libraries.
		// For endless request processing specify '0'.
		// Note: The priority is higher than `idleTimeout`.
		// Default Value: 0.
		maxRequest int

		// The number of seconds after which on idle lua state will be killed.
		// Available Units: s(econds), m(inutes), h(ours), or d(ays)
		// Note: The priority is lower than `maxRequest`,
		//       A value of 0(d, h, m, s) indicates no limit.
		// Default Value: 1h.
		idleTimeout string
		seconds     int

		// lua state container.
		lss []*lState

		mutex sync.Mutex

		suppliedNum int

		whenNew newFunc
	}
)

func Default() *poolManager {
	return &poolManager{
		maxNum:      DefaultMaxNum,
		startNum:    DefaultStartNum,
		maxRequest:  DefaultMaxRequest,
		idleTimeout: DefaultIdleTimeout,
		seconds:     defaultIdleTimeoutNum * 3600,
	}
}

func New(maxNum, startNum, maxRequest int, idleTimeout string) (*poolManager, error) {
	pm := Default()
	pm.maxNum = getMaxNum(maxNum)
	pm.startNum = getStartNum(pm.maxNum, startNum)
	pm.maxRequest = getMaxRequest(maxRequest)
	n, s, err := getIdleTimeout(idleTimeout)
	if err != nil {
		return nil, err
	}
	pm.idleTimeout = fmt.Sprintf("%d%s", n, s)
	pm.seconds = getIdleTimeoutSeconds(n, s)

	return pm, nil
}

// Start start the lua state pool manager.
func (pm *poolManager) Start(whenNew newFunc) error {
	pm.whenNew = whenNew
	pm.lss = make([]*lState, pm.maxNum)

	for i := 0; i < pm.startNum; i++ {
		ls, err := pm.supply()
		if err != nil {
			pm.Shutdown()

			return err
		}

		pm.Put(ls)
	}

	return nil
}

// Put put lua state into pool.
func (pm *poolManager) Put(ls *lState) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if ls.isExpire() {
		pm.suppliedNum -= 1
		ls.close()
	} else {
		pm.lss = append(pm.lss, ls)
	}
}

// Get get lua state from pool.
func (pm *poolManager) Get() (*lState, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	n := len(pm.lss)
	if n == 0 {
		if pm.suppliedNum >= pm.maxNum {
			return nil, ErrLuaStatePoolFull
		}

		return pm.supply()
	}

	ls := pm.lss[n-1]
	pm.lss = pm.lss[0 : n-1]
	ls.incRequestNum()

	return ls, nil
}

func (pm *poolManager) Close(ls *lState) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	ls.close()
	pm.suppliedNum -= 1
}

func (pm *poolManager) Shutdown() {
	for _, ls := range pm.lss {
		pm.Close(ls)
	}
}

func (pm *poolManager) supply() (*lState, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	ls, err := newLState(pm.maxRequest, pm.idleTimeout, pm.seconds, pm.whenNew)
	if err != nil {
		return nil, err
	}

	pm.suppliedNum += 1

	return ls, nil
}
