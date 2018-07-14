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
	gbin "github.com/jefurry/gola/core/binary"
	"github.com/pkg/errors"
	"github.com/yuin/gopher-lua"
)

func binaryPack(L *lua.LState) int {
	format := L.CheckString(1)

	params := make([]lua.LValue, 0, 1)
	params = append(params, checkPackVal(L, 2))

	top := L.GetTop()
	if top > 2 {
		for i := 3; i <= top; i++ {
			params = append(params, checkPackVal(L, i))
		}
	}

	offset := 0
	buf := bytes.NewBuffer(nil)
	err := gbin.ScanToken(format, func(psym, c byte, num int) error {
		order := gbin.GetByteOrder(psym)

		// skip
		if num == 0 {
			offset += 1

			return nil
		}

		n := num
		if c == 's' {
			n = 1
		}

		if len(params) < (n + offset) {
			return errors.New("not enough parameters")
		}

		switch c {
		case '?':
			for i := offset; i < (offset + num); i++ {
				v := params[i]
				lv, ok := v.(lua.LBool)
				if !ok {
					return errors.Errorf("%s expected, got %s", lua.LTBool, v.Type())
				}

				if err := bbin.Write(buf, order, bool(lv)); err != nil {
					return err
				}
			}
		case 'h', 'H':
			for i := offset; i < (offset + num); i++ {
				v := params[i]
				lv, ok := v.(lua.LNumber)
				if !ok {
					return errors.Errorf("%s expected, got %s", lua.LTNumber, v.Type())
				}

				if err := bbin.Write(buf, order, int16(lv)); err != nil {
					return err
				}
			}
		case 'i', 'l':
			for i := offset; i < (offset + num); i++ {
				v := params[i]
				lv, ok := v.(lua.LNumber)
				if !ok {
					return errors.Errorf("%s expected, got %s", lua.LTNumber, v.Type())
				}

				if err := bbin.Write(buf, order, int32(lv)); err != nil {
					return err
				}
			}
		case 'I', 'L':
			for i := offset; i < (offset + num); i++ {
				v := params[i]
				lv, ok := v.(lua.LNumber)
				if !ok {
					return errors.Errorf("%s expected, got %s", lua.LTNumber, v.Type())
				}

				if err := bbin.Write(buf, order, uint32(lv)); err != nil {
					return err
				}
			}
		case 'q':
			for i := offset; i < (offset + num); i++ {
				v := params[i]
				lv, ok := v.(lua.LNumber)
				if !ok {
					return errors.Errorf("%s expected, got %s", lua.LTNumber, v.Type())
				}

				if err := bbin.Write(buf, order, int64(lv)); err != nil {
					return err
				}
			}
		case 'Q':
			for i := offset; i < (offset + num); i++ {
				v := params[i]
				lv, ok := v.(lua.LNumber)
				if !ok {
					return errors.Errorf("%s expected, got %s", lua.LTNumber, v.Type())
				}

				if err := bbin.Write(buf, order, uint64(lv)); err != nil {
					return err
				}
			}
		case 'f':
			for i := offset; i < (offset + num); i++ {
				v := params[i]
				lv, ok := v.(lua.LNumber)
				if !ok {
					return errors.Errorf("%s expected, got %s", lua.LTNumber, v.Type())
				}

				if err := bbin.Write(buf, order, float32(lv)); err != nil {
					return err
				}
			}
		case 'd':
			for i := offset; i < (offset + num); i++ {
				v := params[i]
				lv, ok := v.(lua.LNumber)
				if !ok {
					return errors.Errorf("%s expected, got %s", lua.LTNumber, v.Type())
				}

				if err := bbin.Write(buf, order, float64(lv)); err != nil {
					return err
				}
			}
		case 's':
			v := params[offset]
			lv, ok := v.(lua.LString)
			if !ok {
				return errors.Errorf("%s expected, got %s", lua.LTString, v.Type())
			}

			vs := string(lv)
			size := len(vs)
			if num < size {
				vs = vs[:num]
			} else if num > size {
				left := num - size
				vs = string(append([]byte(vs), make([]byte, left, left)...))
			}

			if err := bbin.Write(buf, order, []byte(vs)); err != nil {
				return err
			}
		}

		offset += n

		return nil
	})

	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(buf.String()))

	return 1
}
