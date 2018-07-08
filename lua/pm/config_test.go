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

func TestConfig(t *testing.T) {
	var c *Config
	var err error

	c, err = NewConfig(30, -1, -1, 120, "1h")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	if !assert.Equal(t, 30, c.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 16, c.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 0, c.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1h", c.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 3600, c.seconds, "seconds mismatching") {
		return
	}

	c, err = NewConfig(30, 20, 4, 120, "1d")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	if !assert.Equal(t, 30, c.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 20, c.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 4, c.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1d", c.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 24*3600, c.seconds, "seconds mismatching") {
		return
	}

	options := c.Options()
	if !assert.Equal(t, false, options.SkipOpenLibs, "SkipOpenLibs should be false") {
		return
	}

	options.SkipOpenLibs = true
	c.SetOptions(options)
	if !assert.Equal(t, true, options.SkipOpenLibs, "SkipOpenLibs should be false") {
		return
	}
}
