// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package binary implements binary encoding for Lua.
// implements pack/unpack. like struct.pack/struct.unpack of python2(https://docs.python.org/2/library/struct.html).
package binary

import (
	"github.com/yuin/gopher-lua"
)

const (
	BinaryLibName = "encoding.binary"
)

const (
	LittleEndian ByteOrder = iota
	BigEndian
)

const (
	Int8 DataType = iota
	Int16
	Int32
	Int64
	Int
	Uint8
	Uint16
	Uint32
	Uint64
	Uint
	Float32
	Float64
	String
	Byte
	Bool
)

type (
	ByteOrder int
	DataType  int
)

func Open(L *lua.LState) {
	L.PreloadModule(BinaryLibName, Loader)
}

func Loader(L *lua.LState) int {
	binmod := L.SetFuncs(L.NewTable(), binaryFuncs)
	L.Push(binmod)

	binaryRegisterReaderMetatype(L)
	binaryRegisterWriterMetatype(L)
	binaryRegisterUnpackRetbMetatype(L)

	for k, v := range binaryByteOrderFields {
		binmod.RawSetString(k, lua.LNumber(v))
	}

	for k, v := range binaryDataTypeFields {
		binmod.RawSetString(k, lua.LNumber(v))
	}

	return 1
}

var binaryFuncs = map[string]lua.LGFunction{
	"newReader": binaryReaderNew,
	"newWriter": binaryWriterNew,
	"pack":      binaryPack,
	"unpack":    binaryUnpack,
}

var binaryByteOrderFields = map[string]ByteOrder{
	"LITTLE_ENDIAN": LittleEndian,
	"BIG_ENDIAN":    BigEndian,
}

var binaryDataTypeFields = map[string]DataType{
	"INT8":    Int8,
	"INT16":   Int16,
	"INT32":   Int32,
	"INT64":   Int64,
	"INT":     Int,
	"UINT8":   Uint8,
	"UINT16":  Uint16,
	"UINT32":  Uint32,
	"UINT64":  Uint64,
	"UINT":    Uint,
	"FLOAT32": Float32,
	"FLOAT64": Float64,
	"STRING":  String,
	"BYTE":    Byte,
	"BOOL":    Bool,
}
