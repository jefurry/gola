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

func TestNewLState(t *testing.T) {
	var ls *lState
	var err error

	for _, v := range []struct {
		maxRequest  int
		idleTimeout string
		seconds     int
		whenNew     newFunc
	}{
		{20, "1d", 24 * 60 * 60, nil},
		{100, "3h", 3600, nil},
		{45, "1m", 60 * 60, nil},
		{77, "10s", 10, nil},
	} {
		ls, err = newLState(v.maxRequest, v.idleTimeout, v.seconds, v.whenNew)
		if !assert.NoError(t, err, `newLState should succeed`) {
			return
		}
		for i := 0; i < v.maxRequest; i++ {
			if !assert.Equal(t, false, ls.isExpire(), "value mismatching") {
				return
			}
			ls.incRequestNum()
		}

		if !assert.Equal(t, true, ls.isExpire(), "value mismatching") {
			return
		}
	}
}
