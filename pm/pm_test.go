// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefault(t *testing.T) {
	pm := Default()
	defer pm.Shutdown()

	if !assert.Equal(t, DefaultMaxNum, pm.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, DefaultStartNum, pm.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, DefaultMaxRequest, pm.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, DefaultIdleTimeout, pm.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 3600, pm.seconds, "seconds mismatching") {
		return
	}
}

func TestNew(t *testing.T) {
	var pm *poolManager
	var err error

	pm, err = New(45, -1, 230, "1d")
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}
	if !assert.Equal(t, 45, pm.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 23, pm.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 230, pm.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1d", pm.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 24*60*60, pm.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()

	pm, err = New(45, -1, 230, "-2d")
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}
	if !assert.Equal(t, 45, pm.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 23, pm.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 230, pm.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1d", pm.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 24*60*60, pm.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()

	pm, err = New(50, 30, 500, "2h")
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}
	if !assert.Equal(t, 50, pm.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 30, pm.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 500, pm.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "2h", pm.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 2*60*60, pm.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()

	pm, err = New(50, 30, 500, "-3h")
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}
	if !assert.Equal(t, 50, pm.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 30, pm.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 500, pm.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1h", pm.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 60*60, pm.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()

	pm, err = New(83, 0, 450, "78m")
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}
	if !assert.Equal(t, 83, pm.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 42, pm.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 450, pm.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "78m", pm.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 78*60, pm.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()

	pm, err = New(83, 0, 450, "-98m")
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}
	if !assert.Equal(t, 83, pm.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 42, pm.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 450, pm.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1m", pm.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 60, pm.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()

	pm, err = New(59, 100, 763, "1583s")
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}
	if !assert.Equal(t, 59, pm.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 30, pm.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 763, pm.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1583s", pm.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 1583, pm.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()

	pm, err = New(59, 100, 763, "-583s")
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}
	if !assert.Equal(t, 59, pm.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 30, pm.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 763, pm.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1s", pm.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 1, pm.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()
}
