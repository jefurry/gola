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
	encodingBinaryWriterTypeName = EncodingBinaryLibName + ".BINARY_WRITER*"
)

type (
	encodingBinaryWriter struct {
		w *bytes.Buffer
	}
)

func encodingBinaryWriterNew(L *lua.LState) int {
	binstr := L.OptString(1, "")

	w := bytes.NewBuffer([]byte(binstr))
	wr := &encodingBinaryWriter{w: w}

	ud := newBinaryWriter(L, wr)

	L.Push(ud)

	return 1

}

func writeBinaryNumber(L *lua.LState, w *bytes.Buffer, order binary.ByteOrder, v interface{}) error {
	if err := binary.Write(w, order, v); err != nil {
		return err
	}

	return nil
}

func encodingBinaryWriterLen(L *lua.LState) int {
	bw := checkBinaryWriter(L, 1)

	L.Push(lua.LNumber(bw.w.Len()))

	return 1
}

func encodingBinaryWriterReset(L *lua.LState) int {
	bw := checkBinaryWriter(L, 1)

	bw.w.Reset()

	return 0
}

func encodingBinaryWriterString(L *lua.LState) int {
	bw := checkBinaryWriter(L, 1)

	s := bw.w.String()

	L.Push(lua.LString(s))

	return 1
}

func encodingBinaryWriterWrite(L *lua.LState) int {
	bw := checkBinaryWriter(L, 1)
	L.CheckTypes(2, lua.LTNumber, lua.LTString, lua.LTBool)
	value := L.CheckAny(2)
	dtype := DataType(L.CheckInt(3))
	opt := L.OptInt(4, 0) // string length

	bo := ByteOrder(opt)

	var order binary.ByteOrder
	if bo == LittleEndian {
		order = binary.LittleEndian
	} else {
		order = binary.BigEndian
	}

	typ := value.Type()
	switch typ {
	case lua.LTNumber:
		val, _ := value.(lua.LNumber)

		switch dtype {
		case Int8:
			v := int8(val)
			if err := writeBinaryNumber(L, bw.w, order, &v); err != nil {
				L.Push(lua.LFalse)
				L.Push(lua.LString(err.Error()))

				return 2
			}

			L.Push(lua.LTrue)

			return 1
		case Int16:
			v := int16(val)
			if err := writeBinaryNumber(L, bw.w, order, &v); err != nil {
				L.Push(lua.LFalse)
				L.Push(lua.LString(err.Error()))

				return 2
			}

			L.Push(lua.LTrue)

			return 1
		case Int32:
			v := int32(val)
			if err := writeBinaryNumber(L, bw.w, order, &v); err != nil {
				L.Push(lua.LFalse)
				L.Push(lua.LString(err.Error()))

				return 2
			}

			L.Push(lua.LTrue)

			return 1
		case Int64:
			v := int64(val)
			if err := writeBinaryNumber(L, bw.w, order, &v); err != nil {
				L.Push(lua.LFalse)
				L.Push(lua.LString(err.Error()))

				return 2
			}

			L.Push(lua.LTrue)

			return 1
		case Int:
			v := int(val)
			if err := writeBinaryNumber(L, bw.w, order, &v); err != nil {
				L.Push(lua.LFalse)
				L.Push(lua.LString(err.Error()))

				return 2
			}

			L.Push(lua.LTrue)

			return 1
		case Uint8:
			v := uint8(val)
			if err := writeBinaryNumber(L, bw.w, order, &v); err != nil {
				L.Push(lua.LFalse)
				L.Push(lua.LString(err.Error()))

				return 2
			}

			L.Push(lua.LTrue)

			return 1
		case Uint16:
			v := uint16(val)
			if err := writeBinaryNumber(L, bw.w, order, &v); err != nil {
				L.Push(lua.LFalse)
				L.Push(lua.LString(err.Error()))

				return 2
			}

			L.Push(lua.LTrue)

			return 1
		case Uint32:
			v := uint32(val)
			if err := writeBinaryNumber(L, bw.w, order, &v); err != nil {
				L.Push(lua.LFalse)
				L.Push(lua.LString(err.Error()))

				return 2
			}

			L.Push(lua.LTrue)

			return 1
		case Uint64:
			v := uint64(val)
			if err := writeBinaryNumber(L, bw.w, order, &v); err != nil {
				L.Push(lua.LFalse)
				L.Push(lua.LString(err.Error()))

				return 2
			}

			L.Push(lua.LTrue)

			return 1
		case Uint:
			v := uint(val)
			if err := writeBinaryNumber(L, bw.w, order, &v); err != nil {
				L.Push(lua.LFalse)
				L.Push(lua.LString(err.Error()))

				return 2
			}

			L.Push(lua.LTrue)

			return 1
		case Float32:
			v := float32(val)
			if err := writeBinaryNumber(L, bw.w, order, &v); err != nil {
				L.Push(lua.LFalse)
				L.Push(lua.LString(err.Error()))

				return 2
			}

			L.Push(lua.LTrue)

			return 1
		case Float64:
			v := float64(val)
			if err := writeBinaryNumber(L, bw.w, order, &v); err != nil {
				L.Push(lua.LFalse)
				L.Push(lua.LString(err.Error()))

				return 2
			}

			L.Push(lua.LTrue)

			return 1
		case Byte:
			v := byte(val)
			if err := writeBinaryNumber(L, bw.w, order, &v); err != nil {
				L.Push(lua.LFalse)
				L.Push(lua.LString(err.Error()))

				return 2
			}

			L.Push(lua.LTrue)

			return 1
		default:
			L.Push(lua.LFalse)
			L.Push(lua.LString("invalid data type"))

			return 2
		}
	case lua.LTString:
		val, _ := value.(lua.LString)
		vv := []byte(val)
		if dtype == Byte {
			if err := writeBinaryNumber(L, bw.w, order, &vv[0]); err != nil {
				L.Push(lua.LFalse)
				L.Push(lua.LString(err.Error()))

				return 2
			}

			L.Push(lua.LTrue)

			return 1
		} else {
			if err := writeBinaryNumber(L, bw.w, order, &vv); err != nil {
				L.Push(lua.LFalse)
				L.Push(lua.LString(err.Error()))

				return 2
			}

			L.Push(lua.LTrue)

			return 1
		}

	case lua.LTBool:
		val, _ := value.(lua.LBool)
		v := bool(val)

		if err := writeBinaryNumber(L, bw.w, order, &v); err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))

			return 2
		}

		L.Push(lua.LTrue)

		return 1
	}

	L.Push(lua.LFalse)
	L.Push(lua.LString("invalid data type"))

	return 2
}

func encodingRegisterBinaryWriterMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(encodingBinaryWriterTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), encodingBinaryWriterFuncs))
	L.SetField(mt, "__tostring", L.NewFunction(encodingBinaryWriterString))
}

var encodingBinaryWriterFuncs = map[string]lua.LGFunction{
	"write": encodingBinaryWriterWrite,
	"len":   encodingBinaryWriterLen,
	"reset": encodingBinaryWriterReset,
}

var encodingBinaryWriterFields = map[string]lua.LValue{}
