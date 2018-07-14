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

func TestScanToken(t *testing.T) {
	for _, v := range []struct {
		format          string
		length          int
		haserr          bool
		precursorSymbol string
	}{
		{
			format: "2sil04s98L",
			haserr: false,
		},
		{
			format: "!fqnN",
			haserr: true,
		},
		{
			format: ">3i 8s i Q",
			haserr: false,
		},
	} {
		err := scanToken(v.format, func(psym, c byte, num int) error {
			return nil
		})

		if v.haserr {
			if !assert.Error(t, err, "scanToken should be not succeed") {
				return
			}
		} else {
			if !assert.NoError(t, err, "scanToken should succeed") {
				return
			}
		}
	}
}

func TestSplitToken(t *testing.T) {
	for _, v := range []struct {
		format          string
		length          int
		haserr          bool
		precursorSymbol string
	}{
		{
			format:          "2sil04s98L",
			length:          5,
			haserr:          false,
			precursorSymbol: "@",
		},
		{
			format:          "!fqnN",
			length:          4,
			haserr:          true,
			precursorSymbol: "!",
		},
		{
			format:          ">3i 8s i Q",
			length:          4,
			haserr:          false,
			precursorSymbol: ">",
		},
	} {
		tokens, err := splitToken(v.format)
		if v.haserr {
			if !assert.Error(t, err, "splitToken should be not succeed") {
				return
			}
		} else {
			if !assert.NoError(t, err, "splitToken should succeed") {
				return
			}

			if !assert.Equal(t, v.length+1, len(tokens), "token length mismatching") {
				return
			}

			if !assert.Equal(t, v.precursorSymbol, tokens[0], "token precursorSymbol mismatching") {
				return
			}
		}
	}
}
