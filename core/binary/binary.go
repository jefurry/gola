// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package binary implements pack/unpack.
// like struct.pack/struct.unpack of python2(https://docs.python.org/2/library/struct.html).
package binary

import (
	"encoding/binary"
)

const (
	AtPrecursorSymbol    byte = '@'
	EqualPrecursorSymbol      = '='
	LTPrecursorSymbol         = '<'
	GTPrecursorSymbol         = '>'
	BangPrecursorSymbol       = '!'
)

const (
	DefaultPrecursorSymbol = AtPrecursorSymbol
)

var (
	AllowFormatSymbols = []byte{
		'?', 'h', 'H', 'i',
		'I', 'l', 'L', 'q',
		'Q', 'f', 'd', 's',
	}
)

var (
	nativeByteOrder binary.ByteOrder = binary.LittleEndian
)

func init() {
	i := int16(0x1234)
	b := int8(i)

	if b == 0x12 {
		nativeByteOrder = binary.BigEndian
	}
}
