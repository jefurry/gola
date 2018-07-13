// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pm

import (
	"bytes"
	"strconv"
)

func getMaxNum(maxNum int) int {
	if maxNum < 0 {
		return DefaultMaxNum
	}

	return maxNum
}

func getStartNum(maxNum, startNum int) int {
	if maxNum <= 0 {
		if startNum <= 0 {
			startNum = DefaultStartNum
		}
	} else {
		if startNum <= 0 || startNum > maxNum {
			startNum = int(maxNum/2) + 1
		}
	}

	return startNum
}

func getMaxRequest(maxRequest int) int {
	if maxRequest < 0 {
		return DefaultMaxRequest
	}

	return maxRequest
}

func getIdleTimeout(idleTimeout string) (int, string, error) {
	it := []byte(idleTimeout)
	l := len(it)
	if l < 2 {
		return 0, "", ErrIdleTimeoutFormat
	}

	last := string(it[l-1])
	if last != idleTimeoutDay && last != idleTimeoutHour && last != idleTimeoutMinute && last != idleTimeoutSecond {
		return 0, "", ErrIdleTimeoutFormat
	}

	it = bytes.TrimRight(it, idleTimeoutDay+idleTimeoutHour+idleTimeoutMinute+idleTimeoutSecond)
	n, err := strconv.Atoi(string(it))
	if err != nil {
		return 0, "", err
	}

	if n < 0 {
		n = defaultIdleTimeoutNum
	}

	return n, last, nil
}

func getIdleTimeoutSeconds(n int, s string) int {
	if n <= 0 {
		return 0
	}

	var multiple int
	if s == idleTimeoutDay {
		multiple = 24 * 60 * 60
	} else if s == idleTimeoutHour {
		multiple = 60 * 60
	} else if s == idleTimeoutMinute {
		multiple = 60
	} else {
		multiple = 1
	}

	return n * multiple
}

func getRequestTerminateTimeout(rtt int) int {
	if rtt < 0 {
		return DefaultRequestTerminateTimeout
	}

	return rtt
}
