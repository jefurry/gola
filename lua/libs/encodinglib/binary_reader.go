// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package encodinglib

import (
	"bytes"
	"encoding/binary"
	"github.com/yuin/gopher-lua"
)

const (
	encodingBinaryReaderTypeName = EncodingBinaryLibName + ".BINARY_READER*"
)

type (
	encodingBinaryReader struct {
		r *bytes.Buffer
	}
)

func encodingBinaryReaderNew(L *lua.LState) int {
	binstr := L.CheckString(1)

	r := bytes.NewBuffer([]byte(binstr))
	br := &encodingBinaryReader{r: r}

	ud := newBinaryReader(L, br)

	L.Push(ud)

	return 1

}

func readBinaryNumber(L *lua.LState, r *bytes.Buffer, order binary.ByteOrder, v interface{}) error {
	if err := binary.Read(r, order, v); err != nil {
		return err
	}

	return nil
}

func encodingBinaryReaderLen(L *lua.LState) int {
	br := checkBinaryReader(L, 1)

	L.Push(lua.LNumber(br.r.Len()))

	return 1
}

func encodingBinaryReaderReset(L *lua.LState) int {
	br := checkBinaryReader(L, 1)

	br.r.Reset()

	return 0
}

func encodingBinaryReaderString(L *lua.LState) int {
	br := checkBinaryReader(L, 1)

	s := br.r.String()

	L.Push(lua.LString(s))

	return 1
}

func encodingBinaryReaderRead(L *lua.LState) int {
	br := checkBinaryReader(L, 1)
	dtype := DataType(L.CheckInt(2))
	opt := L.OptInt(3, 0) // string

	bo := ByteOrder(opt)

	var order binary.ByteOrder
	if bo == LittleEndian {
		order = binary.LittleEndian
	} else {
		order = binary.BigEndian
	}

	switch dtype {
	case Int8:
		var v int8
		if err := readBinaryNumber(L, br.r, order, &v); err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))

			return 2
		}

		L.Push(lua.LNumber(v))

		return 1
	case Int16:
		var v int16
		if err := readBinaryNumber(L, br.r, order, &v); err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))

			return 2
		}

		L.Push(lua.LNumber(v))

		return 1
	case Int32:
		var v int32
		if err := readBinaryNumber(L, br.r, order, &v); err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))

			return 2
		}

		L.Push(lua.LNumber(v))

		return 1
	case Int64:
		var v int64
		if err := readBinaryNumber(L, br.r, order, &v); err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))

			return 2
		}

		L.Push(lua.LNumber(v))

		return 1
	case Int:
		var v int
		if err := readBinaryNumber(L, br.r, order, &v); err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))

			return 2
		}

		L.Push(lua.LNumber(v))

		return 1
	case Uint8:
		var v uint8
		if err := readBinaryNumber(L, br.r, order, &v); err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))

			return 2
		}

		L.Push(lua.LNumber(v))

		return 1
	case Uint16:
		var v uint16
		if err := readBinaryNumber(L, br.r, order, &v); err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))

			return 2
		}

		L.Push(lua.LNumber(v))

		return 1
	case Uint32:
		var v uint32
		if err := readBinaryNumber(L, br.r, order, &v); err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))

			return 2
		}

		L.Push(lua.LNumber(v))

		return 1
	case Uint64:
		var v uint64
		if err := readBinaryNumber(L, br.r, order, &v); err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))

			return 2
		}

		L.Push(lua.LNumber(v))

		return 1
	case Uint:
		var v uint
		if err := readBinaryNumber(L, br.r, order, &v); err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))

			return 2
		}

		L.Push(lua.LNumber(v))

		return 1
	case Float32:
		var v float32
		if err := readBinaryNumber(L, br.r, order, &v); err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))

			return 2
		}

		L.Push(lua.LNumber(v))

		return 1
	case Float64:
		var v float64
		if err := readBinaryNumber(L, br.r, order, &v); err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))

			return 2
		}

		L.Push(lua.LNumber(v))

		return 1
	case Byte:
		var v byte
		if err := readBinaryNumber(L, br.r, order, &v); err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))

			return 2
		}

		L.Push(lua.LNumber(v))

		return 1
	case String:
		v := make([]byte, 0, opt)
		if err := binary.Read(br.r, order, &v); err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))

			return 2
		}

		L.Push(lua.LString(string(v)))

		return 1
	}

	L.Push(lua.LFalse)
	L.Push(lua.LString("invalid data type"))

	return 2
}

func encodingRegisterBinaryReaderMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(encodingBinaryReaderTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), encodingBinaryReaderFuncs))
	L.SetField(mt, "__tostring", L.NewFunction(encodingBinaryReaderString))
}

var encodingBinaryReaderFuncs = map[string]lua.LGFunction{
	"read":  encodingBinaryReaderRead,
	"len":   encodingBinaryReaderLen,
	"reset": encodingBinaryReaderReset,
}

var encodingBinaryReaderFields = map[string]lua.LValue{}
