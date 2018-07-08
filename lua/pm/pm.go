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
	DefaultIdleTimeout = "1h"
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
	LPM struct {
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

func newLPM(ctx context.Context, config *Config, whenNews ...NewFunc) (*LPM, error) {
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

	lpm := &LPM{
		config:   c,
		opStatus: OpReady,
		whenNew:  whenNew,
	}

	err := lpm.start(ctx)
	if err != nil {
		return nil, err
	}

	return lpm, nil
}

func Default(ctx context.Context, whenNews ...NewFunc) (*LPM, error) {
	return newLPM(ctx, nil, whenNews...)
}

func New(ctx context.Context, config *Config, whenNews ...NewFunc) (*LPM, error) {
	if config == nil {
		return Default(ctx, whenNews...)
	}

	return newLPM(ctx, config, whenNews...)
}

func (lpm *LPM) Load(ctx context.Context, reader io.Reader, name string) (*lua.LFunction, error) {
	ls, err := lpm.get(ctx)
	if err != nil {
		return nil, err
	}

	defer lpm.put(ls)

	fn, err := ls.L.Load(reader, name)
	if err != nil {
		lpm.Close(ls)

		return nil, err
	}

	return fn, nil
}

func (lpm *LPM) LoadFile(ctx context.Context, path string) (*lua.LFunction, error) {
	ls, err := lpm.get(ctx)
	if err != nil {
		return nil, err
	}

	defer lpm.put(ls)

	fn, err := ls.L.LoadFile(path)
	if err != nil {
		lpm.Close(ls)

		return nil, err
	}

	return fn, nil
}

func (lpm *LPM) LoadString(ctx context.Context, source string) (*lua.LFunction, error) {
	ls, err := lpm.get(ctx)
	if err != nil {
		return nil, err
	}

	defer lpm.put(ls)

	fn, err := ls.L.LoadString(source)
	if err != nil {
		lpm.Close(ls)

		return nil, err
	}

	return fn, nil
}

func (lpm *LPM) DoFile(ctx context.Context, path string) error {
	ls, err := lpm.get(ctx)
	if err != nil {
		return err
	}

	defer lpm.put(ls)

	if err := ls.L.DoFile(path); err != nil {
		lpm.Close(ls)

		return err
	}

	return nil
}

func (lpm *LPM) DoString(ctx context.Context, source string) error {
	ls, err := lpm.get(ctx)
	if err != nil {
		return err
	}

	defer lpm.put(ls)

	if err := ls.L.DoString(source); err != nil {
		lpm.Close(ls)

		return err
	}

	return nil
}

func (lpm *LPM) Status() OpStatus {
	lpm.lock.Lock()
	defer lpm.lock.Unlock()

	return lpm.opStatus
}

func (lpm *LPM) StatusString() string {
	opStatus := lpm.Status()

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

func (lpm *LPM) ServingNum() int {
	lpm.lock.Lock()
	defer lpm.lock.Unlock()

	return lpm.servingNum
}

func (lpm *LPM) TotalRequestedNum() int {
	lpm.lock.Lock()
	defer lpm.lock.Unlock()

	return lpm.totalRequestedNum
}

func (lpm *LPM) Len() int {
	lpm.lock.Lock()
	defer lpm.lock.Unlock()

	return lpm.length
}

func (lpm *LPM) Cap() int {
	return lpm.config.maxNum
}

func (lpm *LPM) Close(ls *lState) {
	lpm.lock.Lock()
	lpm.lock.Unlock()

	if !ls.serving {
		panic("lua state not running")
	}

	lpm.servingNum -= 1

	if lpm.servingNum <= 0 {
		lpm.cond.Signal()
	}

	lpm.length -= 1
	ls.close()
}

func (lpm *LPM) Shutdown() {
	lpm.readyExit()
	lpm.exit()

	lpm.length = 0
	lpm.servingNum = 0
	lpm.lss = nil
}

func (lpm *LPM) Restart(ctx context.Context) error {
	lpm.Shutdown()

	if err := lpm.start(ctx); err != nil {
		return err
	}

	return nil
}

func (lpm *LPM) start(ctx context.Context) error {
	lpm.lss = make([]*lState, 0, lpm.config.maxNum)
	lpm.lock = new(sync.Mutex)
	lpm.cond = sync.NewCond(lpm.lock)

	for i := 0; i < lpm.config.startNum; i++ {
		ls, err := lpm.gen(ctx)
		if err != nil {
			lpm.Shutdown()

			return err
		}

		lpm.length += 1
		ls.setServing(false)
		lpm.put(ls)
	}

	lpm.opStatus = OpRunning

	return nil
}

// put put lua state into pool.
func (lpm *LPM) put(ls *lState) error {
	if ls.closed {
		return ErrLSClosed
	}

	lpm.lock.Lock()
	defer lpm.lock.Unlock()

	if lpm.opStatus == OpExiting {
		lpm.Close(ls)

		return ErrLSPExiting
	}

	if lpm.opStatus == OpDead {
		lpm.Close(ls)

		return ErrLSPDead
	}

	if ls.mustTerminate() {
		lpm.Close(ls)

		return nil
	}

	if ls.serving {
		lpm.servingNum -= 1
	}

	ls.setServing(false)

	lpm.lss = append(lpm.lss, ls)

	return nil
}

func (lpm *LPM) get(ctx context.Context) (*lState, error) {
	lpm.lock.Lock()
	defer lpm.lock.Unlock()

	if lpm.opStatus == OpExiting {
		return nil, ErrLSPExiting
	}

	if lpm.opStatus == OpDead {
		return nil, ErrLSPDead
	}

	var ls *lState
	n := len(lpm.lss)
	if n == 0 {
		if lpm.config.maxNum != 0 && lpm.length >= lpm.config.maxNum {
			return nil, ErrLSPFulled
		}

		lstate, err := lpm.gen(ctx)
		if err != nil {
			return nil, err
		}

		lpm.length += 1
		ls = lstate
	} else {
		ls = lpm.lss[n-1]
		lpm.lss = lpm.lss[0 : n-1]
	}

	ls.setServing(true)
	lpm.servingNum += 1
	lpm.totalRequestedNum += 1
	ls.incRequestNum()

	return ls, nil
}

func (lpm *LPM) gen(ctx context.Context) (*lState, error) {
	ls, err := newLState(ctx, lpm.config.maxRequest, lpm.config.idleTimeout,
		lpm.config.seconds, lpm.config.options, lpm.whenNew)
	if err != nil {
		return nil, err
	}

	return ls, nil
}

func (lpm *LPM) readyExit() {
	lpm.lock.Lock()
	defer lpm.lock.Unlock()

	lpm.opStatus = OpExiting
}

func (lpm *LPM) exit() {
	lpm.lock.Lock()
	defer lpm.lock.Unlock()

	for lpm.servingNum > 0 {
		lpm.cond.Wait()
	}

	for _, ls := range lpm.lss {
		ls.close()
	}
}
