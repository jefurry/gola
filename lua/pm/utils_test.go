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

func TestGetMaxNum(t *testing.T) {
	if !assert.Equal(t, DefaultMaxNum, getMaxNum(-1), "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 0, getMaxNum(0), "maxNum mismatching") {
		return
	}
}

func TestGetStartNum(t *testing.T) {
	if !assert.Equal(t, 3, getStartNum(5, -1), "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 1, getStartNum(5, 1), "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 1, getStartNum(0, 1), "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 1, getStartNum(0, -1), "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 26, getStartNum(50, -1), "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 19, getStartNum(50, 19), "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 26, getStartNum(51, 60), "startNum mismatching") {
		return
	}
}

func TestGetMaxRequest(t *testing.T) {
	if !assert.Equal(t, 0, getMaxRequest(0), "getMaxRequest mismatching") {
		return
	}

	if !assert.Equal(t, 0, getMaxRequest(-1), "getMaxRequest mismatching") {
		return
	}

	if !assert.Equal(t, 10, getMaxRequest(10), "getMaxRequest mismatching") {
		return
	}
}

func TestGetIdleTimeout(t *testing.T) {
	var n int
	var s string
	var err error

	n, s, err = getIdleTimeout("5d")
	if !assert.NoError(t, err, `getIdleTimeout should succeed`) {
		return
	}
	if !assert.Equal(t, 5, n, "getIdleTimeout mismatching") {
		return
	}
	if !assert.Equal(t, "d", s, "getIdleTimeout mismatching") {
		return
	}

	n, s, err = getIdleTimeout("0d")
	if !assert.NoError(t, err, `getIdleTimeout should succeed`) {
		return
	}
	if !assert.Equal(t, 0, n, "getIdleTimeout mismatching") {
		return
	}
	if !assert.Equal(t, "d", s, "getIdleTimeout mismatching") {
		return
	}

	n, s, err = getIdleTimeout("1h")
	if !assert.NoError(t, err, `getIdleTimeout should succeed`) {
		return
	}
	if !assert.Equal(t, 1, n, "getIdleTimeout mismatching") {
		return
	}
	if !assert.Equal(t, "h", s, "getIdleTimeout mismatching") {
		return
	}

	n, s, err = getIdleTimeout("12h")
	if !assert.NoError(t, err, `getIdleTimeout should succeed`) {
		return
	}
	if !assert.Equal(t, 12, n, "getIdleTimeout mismatching") {
		return
	}
	if !assert.Equal(t, "h", s, "getIdleTimeout mismatching") {
		return
	}

	n, s, err = getIdleTimeout("8m")
	if !assert.NoError(t, err, `getIdleTimeout should succeed`) {
		return
	}
	if !assert.Equal(t, 8, n, "getIdleTimeout mismatching") {
		return
	}
	if !assert.Equal(t, "m", s, "getIdleTimeout mismatching") {
		return
	}

	n, s, err = getIdleTimeout("57m")
	if !assert.NoError(t, err, `getIdleTimeout should succeed`) {
		return
	}
	if !assert.Equal(t, 57, n, "getIdleTimeout mismatching") {
		return
	}
	if !assert.Equal(t, "m", s, "getIdleTimeout mismatching") {
		return
	}

	n, s, err = getIdleTimeout("9s")
	if !assert.NoError(t, err, `getIdleTimeout should succeed`) {
		return
	}
	if !assert.Equal(t, 9, n, "getIdleTimeout mismatching") {
		return
	}
	if !assert.Equal(t, "s", s, "getIdleTimeout mismatching") {
		return
	}

	n, s, err = getIdleTimeout("43s")
	if !assert.NoError(t, err, `getIdleTimeout should succeed`) {
		return
	}
	if !assert.Equal(t, 43, n, "getIdleTimeout mismatching") {
		return
	}
	if !assert.Equal(t, "s", s, "getIdleTimeout mismatching") {
		return
	}
}

func TestGetIdleTimeoutSeconds(t *testing.T) {
	if !assert.Equal(t, 2*24*60*60, getIdleTimeoutSeconds(2, "d"), "getIdleTimeoutSeconds mismatching") {
		return
	}

	if !assert.Equal(t, 5*60*60, getIdleTimeoutSeconds(5, "h"), "getIdleTimeoutSeconds mismatching") {
		return
	}

	if !assert.Equal(t, 65*60, getIdleTimeoutSeconds(65, "m"), "getIdleTimeoutSeconds mismatching") {
		return
	}

	if !assert.Equal(t, 148, getIdleTimeoutSeconds(148, "s"), "getIdleTimeoutSeconds mismatching") {
		return
	}
}
