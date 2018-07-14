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
	"github.com/yuin/gopher-lua"
	"io"
	//"github.com/pkg/errors"
)

func binaryUnpack(L *lua.LState) int {
	format := L.CheckString(1)
	binstr := L.CheckString(2)

	vals := newUnpackRetb(L, nil)
	buf := bytes.NewBuffer([]byte(binstr))
	err := gbin.ScanToken(format, func(psym, c byte, num int) error {
		order := gbin.GetByteOrder(psym)

		var skipNum int
		switch c {
		case '?':
			skipNum = 1
		case 'h', 'H':
			skipNum = 2
		case 'i', 'I', 'l', 'L', 'f':
			skipNum = 4
		case 'q', 'Q', 'd':
			skipNum = 8
		}

		if num == 0 && skipNum != 0 {
			buf.Next(skipNum)

			return nil
		}

		switch c {
		case '?':
			for i := 0; i < num; i++ {
				var v bool
				if err := bbin.Read(buf, order, &v); err != nil {
					return err
				}

				vals.Append(lua.LBool(v))
			}
		case 'h', 'H':
			for i := 0; i < num; i++ {
				var v int16
				if err := bbin.Read(buf, order, &v); err != nil {
					return err
				}

				vals.Append(lua.LNumber(v))
			}
		case 'i', 'l':
			for i := 0; i < num; i++ {
				var v int32
				if err := bbin.Read(buf, order, &v); err != nil {
					return err
				}

				vals.Append(lua.LNumber(v))
			}
		case 'I', 'L':
			for i := 0; i < num; i++ {
				var v uint32
				if err := bbin.Read(buf, order, &v); err != nil {
					return err
				}

				vals.Append(lua.LNumber(v))
			}
		case 'q':
			for i := 0; i < num; i++ {
				var v int64
				if err := bbin.Read(buf, order, &v); err != nil {
					return err
				}

				vals.Append(lua.LNumber(v))
			}
		case 'Q':
			for i := 0; i < num; i++ {
				var v uint64
				if err := bbin.Read(buf, order, &v); err != nil {
					return err
				}

				vals.Append(lua.LNumber(v))
			}
		case 'f':
			for i := 0; i < num; i++ {
				var v float32
				if err := bbin.Read(buf, order, &v); err != nil {
					return err
				}

				vals.Append(lua.LNumber(v))
			}
		case 'd':
			for i := 0; i < num; i++ {
				var v float64
				if err := bbin.Read(buf, order, &v); err != nil {
					return err
				}

				vals.Append(lua.LNumber(v))
			}
		case 's':
			b := make([]byte, num)
			if err := bbin.Read(buf, order, &b); err != nil {
				if err != io.EOF {
					return err
				}
			}

			vals.Append(lua.LString(string(b)))
		}

		return nil
	})

	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(vals)

	return 1
}
