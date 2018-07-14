// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package binary

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBinary(t *testing.T) {
	bs, err := Pack("?0L3s10s2ifhHiIlLfdqQ", true, 35, "one", "two fields", 0xff, 56, 9.8, -28, 45, 123, 45, -99, 20, 3.5, 5.4, int64(321), uint64(567))
	if !assert.NoError(t, err, "Pack should succeed") {
		return
	}

	if !assert.Equal(t, 74, len(bs), "bytes length mismatching") {
		return
	}

	vals, err := Unpack("?3s10s2ifhHiIlLfdqQ", bs)
	if !assert.NoError(t, err, "Unpack should succeed") {
		return
	}

	if !assert.Equal(t, 16, len(vals), "vals length mismatching") {
		return
	}

	if !assert.Equal(t, true, vals[0], "value mismatching") {
		return
	}

	if !assert.Equal(t, "one", vals[1], "value mismatching") {
		return
	}

	if !assert.Equal(t, "two fields", vals[2], "value mismatching") {
		return
	}

	if !assert.Equal(t, int32(0xff), vals[3], "value mismatching") {
		return
	}

	if !assert.Equal(t, int32(56), vals[4], "value mismatching") {
		return
	}

	if !assert.Equal(t, float32(9.8), vals[5], "value mismatching") {
		return
	}

	if !assert.Equal(t, int16(-28), vals[6], "value mismatching") {
		return
	}

	if !assert.Equal(t, uint16(45), vals[7], "value mismatching") {
		return
	}

	if !assert.Equal(t, int32(123), vals[8], "value mismatching") {
		return
	}

	if !assert.Equal(t, uint32(45), vals[9], "value mismatching") {
		return
	}

	if !assert.Equal(t, int32(-99), vals[10], "value mismatching") {
		return
	}

	if !assert.Equal(t, uint32(20), vals[11], "value mismatching") {
		return
	}

	if !assert.Equal(t, float32(3.5), vals[12], "value mismatching") {
		return
	}

	if !assert.Equal(t, 5.4, vals[13], "value mismatching") {
		return
	}

	if !assert.Equal(t, int64(321), vals[14], "value mismatching") {
		return
	}

	if !assert.Equal(t, uint64(567), vals[15], "value mismatching") {
		return
	}
}
