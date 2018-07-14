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
	"encoding/binary"
	"io"
)

func Unpack(format string, bs []byte) ([]interface{}, error) {
	vals := make([]interface{}, 0, 0)
	buf := bytes.NewBuffer(bs)
	err := ScanToken(format, func(psym, c byte, num int) error {
		order := GetByteOrder(psym)

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
				if err := binary.Read(buf, order, &v); err != nil {
					return err
				}

				vals = append(vals, v)
			}
		case 'h':
			for i := 0; i < num; i++ {
				var v int16
				if err := binary.Read(buf, order, &v); err != nil {
					return err
				}

				vals = append(vals, v)
			}
		case 'H':
			for i := 0; i < num; i++ {
				var v uint16
				if err := binary.Read(buf, order, &v); err != nil {
					return err
				}

				vals = append(vals, v)
			}
		case 'i', 'l':
			for i := 0; i < num; i++ {
				var v int32
				if err := binary.Read(buf, order, &v); err != nil {
					return err
				}

				vals = append(vals, v)
			}
		case 'I', 'L':
			for i := 0; i < num; i++ {
				var v uint32
				if err := binary.Read(buf, order, &v); err != nil {
					return err
				}

				vals = append(vals, v)
			}
		case 'q':
			for i := 0; i < num; i++ {
				var v int64
				if err := binary.Read(buf, order, &v); err != nil {
					return err
				}

				vals = append(vals, v)
			}
		case 'Q':
			for i := 0; i < num; i++ {
				var v uint64
				if err := binary.Read(buf, order, &v); err != nil {
					return err
				}

				vals = append(vals, v)
			}
		case 'f':
			for i := 0; i < num; i++ {
				var v float32
				if err := binary.Read(buf, order, &v); err != nil {
					return err
				}

				vals = append(vals, v)
			}
		case 'd':
			for i := 0; i < num; i++ {
				var v float64
				if err := binary.Read(buf, order, &v); err != nil {
					return err
				}

				vals = append(vals, v)
			}
		case 's':
			b := make([]byte, num)
			if err := binary.Read(buf, order, &b); err != nil {
				if err != io.EOF {
					return err
				}
			}

			vals = append(vals, string(b))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return vals, nil
}
