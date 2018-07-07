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
	"context"
	"github.com/pkg/errors"
	"github.com/yuin/gopher-lua"
	"io"
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

const (
	DefaultRequestTerminateTimeout = 120
)

const (
	OpReady OpStatus = iota
	OpRunning
	OpExiting
	OpDead
)

var (
	ErrIdleTimeoutFormat = errors.New("not a valid idle timeout format")
	ErrLSPFulled         = errors.New("lua state pool fulled")
	ErrLSPExiting        = errors.New("lua state pool exiting")
	ErrLSPDead           = errors.New("lua state pool dead")
	ErrLSClosed          = errors.New("lua status has closed")
)

type (
	// Operating status
	OpStatus int
	NewFunc  func(l *lua.LState) error
)

type (
	// lua state pool manager.
	PoolManager struct {
		config *Config

		opStatus OpStatus

		// lua state container.
		lss []*lState

		lock *sync.Mutex
		cond *sync.Cond

		length            int
		servingNum        int
		totalRequestedNum int

		whenNew NewFunc
	}
)

func newPM(ctx context.Context, config *Config, whenNews ...NewFunc) (*PoolManager, error) {
	var c *Config
	if config == nil {
		c = &Config{
			maxNum:                  DefaultMaxNum,
			startNum:                DefaultStartNum,
			maxRequest:              DefaultMaxRequest,
			idleTimeout:             DefaultIdleTimeout,
			seconds:                 defaultIdleTimeoutNum * 3600,
			requestTerminateTimeout: DefaultRequestTerminateTimeout,
		}
	} else {
		c = config
	}

	var whenNew NewFunc = nil
	if len(whenNews) > 0 {
		whenNew = whenNews[0]
	}

	pm := &PoolManager{
		config:   c,
		opStatus: OpReady,
		whenNew:  whenNew,
	}

	err := pm.start(ctx)
	if err != nil {
		return nil, err
	}

	return pm, nil
}

func Default(ctx context.Context, whenNews ...NewFunc) (*PoolManager, error) {
	return newPM(ctx, nil, whenNews...)
}

func New(ctx context.Context, config *Config, whenNews ...NewFunc) (*PoolManager, error) {
	if config == nil {
		return Default(ctx, whenNews...)
	}

	return newPM(ctx, config, whenNews...)
}

func (pm *PoolManager) Load(ctx context.Context, reader io.Reader, name string) (*lua.LFunction, error) {
	ls, err := pm.get(ctx)
	if err != nil {
		return nil, err
	}

	defer pm.put(ls)

	return ls.L.Load(reader, name)
}

func (pm *PoolManager) LoadFile(ctx context.Context, path string) (*lua.LFunction, error) {
	ls, err := pm.get(ctx)
	if err != nil {
		return nil, err
	}

	defer pm.put(ls)

	return ls.L.LoadFile(path)
}

func (pm *PoolManager) LoadString(ctx context.Context, source string) (*lua.LFunction, error) {
	ls, err := pm.get(ctx)
	if err != nil {
		return nil, err
	}

	defer pm.put(ls)

	return ls.L.LoadString(source)
}

func (pm *PoolManager) DoFile(ctx context.Context, path string) error {
	ls, err := pm.get(ctx)
	if err != nil {
		return err
	}

	defer pm.put(ls)

	return ls.L.DoFile(path)
}

func (pm *PoolManager) DoString(ctx context.Context, source string) error {
	ls, err := pm.get(ctx)
	if err != nil {
		return err
	}

	defer pm.put(ls)

	return ls.L.DoString(source)
}

func (pm *PoolManager) Status() OpStatus {
	pm.lock.Lock()
	defer pm.lock.Unlock()

	return pm.opStatus
}

func (pm *PoolManager) StatusString() string {
	opStatus := pm.Status()

	if opStatus == OpReady {
		return "Ready"
	} else if opStatus == OpRunning {
		return "Running"
	} else if opStatus == OpExiting {
		return "Exiting"
	} else if opStatus == OpDead {
		return "Dead"
	}

	return "Unknown"
}

func (pm *PoolManager) ServingNum() int {
	pm.lock.Lock()
	defer pm.lock.Unlock()

	return pm.servingNum
}

func (pm *PoolManager) TotalRequestedNum() int {
	pm.lock.Lock()
	defer pm.lock.Unlock()

	return pm.totalRequestedNum
}

func (pm *PoolManager) Len() int {
	pm.lock.Lock()
	defer pm.lock.Unlock()

	return pm.length
}

func (pm *PoolManager) Cap() int {
	return pm.config.maxNum
}

func (pm *PoolManager) Close(ls *lState) {
	pm.lock.Lock()
	pm.lock.Unlock()

	if !ls.serving {
		panic("lua state not running")
	}

	pm.servingNum -= 1

	if pm.servingNum <= 0 {
		pm.cond.Signal()
	}

	pm.length -= 1
	ls.close()
}

func (pm *PoolManager) Shutdown() {
	pm.readyExit()
	pm.exit()

	pm.length = 0
	pm.servingNum = 0
	pm.lss = nil
}

func (pm *PoolManager) Restart(ctx context.Context) error {
	pm.Shutdown()

	if err := pm.start(ctx); err != nil {
		return err
	}

	return nil
}

func (pm *PoolManager) start(ctx context.Context) error {
	pm.lss = make([]*lState, 0, pm.config.maxNum)
	pm.lock = new(sync.Mutex)
	pm.cond = sync.NewCond(pm.lock)

	for i := 0; i < pm.config.startNum; i++ {
		ls, err := pm.gen(ctx)
		if err != nil {
			pm.Shutdown()

			return err
		}

		pm.length += 1
		ls.setServing(false)
		pm.put(ls)
	}

	pm.opStatus = OpRunning

	return nil
}

// put put lua state into pool.
func (pm *PoolManager) put(ls *lState) error {
	if ls.closed {
		return ErrLSClosed
	}

	pm.lock.Lock()
	defer pm.lock.Unlock()

	if pm.opStatus == OpExiting {
		pm.Close(ls)

		return ErrLSPExiting
	}

	if pm.opStatus == OpDead {
		pm.Close(ls)

		return ErrLSPDead
	}

	if ls.mustTerminate() {
		pm.Close(ls)

		return nil
	}

	if ls.serving {
		pm.servingNum -= 1
	}

	ls.setServing(false)

	pm.lss = append(pm.lss, ls)

	return nil
}

func (pm *PoolManager) get(ctx context.Context) (*lState, error) {
	pm.lock.Lock()
	defer pm.lock.Unlock()

	if pm.opStatus == OpExiting {
		return nil, ErrLSPExiting
	}

	if pm.opStatus == OpDead {
		return nil, ErrLSPDead
	}

	var ls *lState
	n := len(pm.lss)
	if n == 0 {
		if pm.config.maxNum != 0 && pm.length >= pm.config.maxNum {
			return nil, ErrLSPFulled
		}

		lstate, err := pm.gen(ctx)
		if err != nil {
			return nil, err
		}

		pm.length += 1
		ls = lstate
	} else {
		ls = pm.lss[n-1]
		pm.lss = pm.lss[0 : n-1]
	}

	ls.setServing(true)
	pm.servingNum += 1
	pm.totalRequestedNum += 1
	ls.incRequestNum()

	return ls, nil
}

func (pm *PoolManager) gen(ctx context.Context) (*lState, error) {
	ls, err := newLState(ctx, pm.config.maxRequest, pm.config.idleTimeout, pm.config.seconds, pm.whenNew)
	if err != nil {
		return nil, err
	}

	return ls, nil
}

func (pm *PoolManager) readyExit() {
	pm.lock.Lock()
	defer pm.lock.Unlock()

	pm.opStatus = OpExiting
}

func (pm *PoolManager) exit() {
	pm.lock.Lock()
	defer pm.lock.Unlock()

	for pm.servingNum > 0 {
		pm.cond.Wait()
	}

	for _, ls := range pm.lss {
		ls.close()
	}
}
