// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package binary

import (
	"bytes"
	bbin "encoding/binary"
	"github.com/yuin/gopher-lua"
)

const (
	binaryWriterTypeName = BinaryLibName + ".WRITER*"
)

type (
	binaryWriter struct {
		w *bytes.Buffer
	}
)

func binaryWriterNew(L *lua.LState) int {
	binstr := L.OptString(1, "")

	w := bytes.NewBuffer([]byte(binstr))
	wr := &binaryWriter{w: w}

	ud := newBinaryWriter(L, wr)

	L.Push(ud)

	return 1

}

func writeBinaryNumber(L *lua.LState, w *bytes.Buffer, order bbin.ByteOrder, v interface{}) error {
	if err := bbin.Write(w, order, v); err != nil {
		return err
	}

	return nil
}

func binaryWriterLen(L *lua.LState) int {
	bw := checkBinaryWriter(L, 1)

	L.Push(lua.LNumber(bw.w.Len()))

	return 1
}

func binaryWriterReset(L *lua.LState) int {
	bw := checkBinaryWriter(L, 1)

	bw.w.Reset()

	return 0
}

func binaryWriterString(L *lua.LState) int {
	bw := checkBinaryWriter(L, 1)

	s := bw.w.String()

	L.Push(lua.LString(s))

	return 1
}

func binaryWriterWrite(L *lua.LState) int {
	bw := checkBinaryWriter(L, 1)
	L.CheckTypes(2, lua.LTNumber, lua.LTString, lua.LTBool)
	value := L.CheckAny(2)
	dtype := DataType(L.CheckInt(3))
	opt := L.OptInt(4, 0) // string length

	bo := ByteOrder(opt)

	var order bbin.ByteOrder
	if bo == LittleEndian {
		order = bbin.LittleEndian
	} else {
		order = bbin.BigEndian
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

func binaryRegisterWriterMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(binaryWriterTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), binaryWriterFuncs))
	L.SetField(mt, "__tostring", L.NewFunction(binaryWriterString))
}

var binaryWriterFuncs = map[string]lua.LGFunction{
	"write": binaryWriterWrite,
	"len":   binaryWriterLen,
	"reset": binaryWriterReset,
}

var binaryWriterFields = map[string]lua.LValue{}
