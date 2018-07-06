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
	NewFunc func(l *lua.LState) error
)

type (
	// lua state pool manager.
	PoolManager struct {
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

		providedNum       int
		servingNum        int
		totalRequestedNum int

		whenNew NewFunc
	}
)

func Default() *PoolManager {
	return &PoolManager{
		maxNum:      DefaultMaxNum,
		startNum:    DefaultStartNum,
		maxRequest:  DefaultMaxRequest,
		idleTimeout: DefaultIdleTimeout,
		seconds:     defaultIdleTimeoutNum * 3600,
	}
}

func New(maxNum, startNum, maxRequest int, idleTimeout string) (*PoolManager, error) {
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
func (pm *PoolManager) Start(whenNew NewFunc) error {
	pm.whenNew = whenNew
	pm.lss = make([]*lState, 0, pm.maxNum)

	for i := 0; i < pm.startNum; i++ {
		ls, err := pm.provide()
		if err != nil {
			pm.Shutdown()

			return err
		}

		pm.Put(ls)
	}

	return nil
}

// Put put lua state into pool.
func (pm *PoolManager) Put(ls *lState) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if ls.isExpire() {
		pm.Close(ls)
	} else {
		if ls.isServing() {
			pm.servingNum -= 1
		}

		ls.setServing(false)
		pm.lss = append(pm.lss, ls)
	}
}

// Get get lua state from pool.
func (pm *PoolManager) Get() (*lState, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	n := len(pm.lss)
	if n == 0 {
		if pm.providedNum >= pm.maxNum {
			return nil, ErrLuaStatePoolFull
		}

		ls, err := pm.provide()
		if err != nil {
			return nil, err
		}

		ls.setServing(true)
		ls.incRequestNum()
		pm.servingNum += 1
		pm.totalRequestedNum += 1

		return ls, nil
	} else {
		ls := pm.lss[n-1]
		pm.lss = pm.lss[0 : n-1]
		ls.setServing(true)
		ls.incRequestNum()
		pm.servingNum += 1
		pm.totalRequestedNum += 1

		return ls, nil
	}
}

func (pm *PoolManager) Size() int {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	return pm.providedNum
}

func (pm *PoolManager) Cap() int {
	return pm.maxNum
}

func (pm *PoolManager) ServingNum() int {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	return pm.servingNum
}

func (pm *PoolManager) TotalRequestedNum() int {
	return pm.totalRequestedNum
}

// TODO: with context
func (pm *PoolManager) Close(ls *lState) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if ls == nil {
		return
	}

	if ls.isServing() {
		pm.servingNum -= 1
	}

	ls.setServing(false)
	pm.providedNum -= 1
	ls.close()
}

// TODO: with context
func (pm *PoolManager) Shutdown() {
	for _, ls := range pm.lss {
		pm.Close(ls)
	}

	pm.servingNum = 0
	pm.providedNum = 0
	pm.lss = nil
}

func (pm *PoolManager) provide() (*lState, error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	ls, err := newLState(pm.maxRequest, pm.idleTimeout, pm.seconds, pm.whenNew)
	if err != nil {
		return nil, err
	}

	pm.providedNum += 1

	return ls, nil
}
