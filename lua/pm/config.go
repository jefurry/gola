// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pm

import (
	"fmt"
	"github.com/yuin/gopher-lua"
)

type (
	Config struct {
		options lua.Options

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

		// The timeout (in seconds) for serving a single request after which the worker process will be terminated.
		// Note: A value of 0 indicates no limit.
		// Note: A value of negative indicates `DefaultRequestTerminateTimeout`.
		// Default Value 120s.
		requestTerminateTimeout int
		seconds                 int
	}
)

func NewConfig(maxNum, startNum, maxRequest, requestTerminateTimeout int, idleTimeout string, options ...lua.Options) (*Config, error) {
	c := &Config{}

	if len(options) > 0 {
		c.options = options[0]
	}

	c.maxNum = getMaxNum(maxNum)
	c.startNum = getStartNum(c.maxNum, startNum)
	c.maxRequest = getMaxRequest(maxRequest)

	n, s, err := getIdleTimeout(idleTimeout)
	if err != nil {
		return nil, err
	}

	c.idleTimeout = fmt.Sprintf("%d%s", n, s)
	c.seconds = getIdleTimeoutSeconds(n, s)
	c.requestTerminateTimeout = getRequestTerminateTimeout(requestTerminateTimeout)

	return c, nil
}

func (c *Config) Options() lua.Options {
	return c.options
}

func (c *Config) SetOptions(options lua.Options) {
	c.options = options
}
