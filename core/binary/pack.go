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
	"github.com/pkg/errors"
	"reflect"
)

func Pack(format string, val interface{}, vals ...interface{}) ([]byte, error) {
	params := make([]interface{}, 0, 1)
	params = append(params, val)
	params = append(params, vals...)

	offset := 0
	buf := bytes.NewBuffer(nil)
	err := ScanToken(format, func(psym, c byte, num int) error {
		order := getByteOrder(psym)

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
		case '?': // bool
			for i := offset; i < (offset + num); i++ {
				v := params[i]
				vb, ok := v.(bool)
				if !ok {
					return errors.Errorf("type of passed value doesn't match to expected '%s=%v' (bool)", reflect.TypeOf(v), v)
				}

				if err := binary.Write(buf, order, vb); err != nil {
					return err
				}
			}
		case 'h': // short
			for i := offset; i < (offset + num); i++ {
				v := params[i]
				vh, ok := v.(int)
				if !ok {
					return errors.Errorf("type of passed value doesn't match to expected '%s=%v' (int, 2 bytes)", reflect.TypeOf(v), v)
				}

				if err := binary.Write(buf, order, int16(vh)); err != nil {
					return err
				}
			}
		case 'H': // unsigned short
			for i := offset; i < (offset + num); i++ {
				v := params[i]
				vuh, ok := v.(int)
				if !ok {
					return errors.Errorf("type of passed value doesn't match to expected '%s=%v' (int, 2 bytes)", reflect.TypeOf(v), v)
				}

				if err := binary.Write(buf, order, uint16(vuh)); err != nil {
					return err
				}
			}
		case 'i': // int
			for i := offset; i < (offset + num); i++ {
				v := params[i]
				vi, ok := v.(int)
				if !ok {
					return errors.Errorf("type of passed value doesn't match to expected '%s=%v' (int, 4 bytes)", reflect.TypeOf(v), v)
				}

				if err := binary.Write(buf, order, int32(vi)); err != nil {
					return err
				}
			}
		case 'I': // unsigned int
			for i := offset; i < (offset + num); i++ {
				v := params[i]
				vui, ok := v.(int)
				if !ok {
					return errors.Errorf("type of passed value doesn't match to expected '%s=%v' (int, 4 bytes)", reflect.TypeOf(v), v)
				}

				if err := binary.Write(buf, order, uint32(vui)); err != nil {
					return err
				}
			}
		case 'l': // long
			for i := offset; i < (offset + num); i++ {
				v := params[i]
				vi, ok := v.(int)
				if !ok {
					return errors.Errorf("type of passed value doesn't match to expected '%s=%v' (int, 4 bytes)", reflect.TypeOf(v), v)
				}

				if err := binary.Write(buf, order, int32(vi)); err != nil {
					return err
				}
			}

		case 'L': // unsigned long
			for i := offset; i < (offset + num); i++ {
				v := params[i]
				vui, ok := v.(int)
				if !ok {
					return errors.Errorf("type of passed value doesn't match to expected '%s=%v' (int, 4 bytes)", reflect.TypeOf(v), v)
				}

				if err := binary.Write(buf, order, uint32(vui)); err != nil {
					return err
				}
			}
		case 'q': // long long
			for i := offset; i < (offset + num); i++ {
				v := params[i]
				vll, ok := v.(int64)
				if !ok {
					return errors.Errorf("type of passed value doesn't match to expected '%s=%v' (int, 8 bytes)", reflect.TypeOf(v), v)
				}

				if err := binary.Write(buf, order, vll); err != nil {
					return err
				}
			}
		case 'Q': // unsigned long long
			for i := offset; i < (offset + num); i++ {
				v := params[i]
				vull, ok := v.(uint64)
				if !ok {
					return errors.Errorf("type of passed value doesn't match to expected '%s=%v' (int, 8 bytes)", reflect.TypeOf(v), v)
				}

				if err := binary.Write(buf, order, vull); err != nil {
					return err
				}
			}
		case 'f': // float
			for i := offset; i < (offset + num); i++ {
				v := params[i]
				vf, ok := v.(float64)
				if !ok {
					return errors.Errorf("type of passed value doesn't match to expected '%s=%v' (float32)", reflect.TypeOf(v), v)
				}

				if err := binary.Write(buf, order, float32(vf)); err != nil {
					return err
				}
			}
		case 'd': // double
			for i := offset; i < (offset + num); i++ {
				v := params[i]
				vd, ok := v.(float64)
				if !ok {
					return errors.Errorf("type of passed value doesn't match to expected '%s=%v' (float64)", reflect.TypeOf(v), v)
				}

				if err := binary.Write(buf, order, vd); err != nil {
					return err
				}
			}
		case 's': // char[]
			v := params[offset]
			vs, ok := v.(string)
			if !ok {
				return errors.Errorf("type of passed value doesn't match to expected '%s=%v' (string)", reflect.TypeOf(v), v)
			}

			size := len(vs)
			if num < size {
				vs = vs[:num]
			} else if num > size {
				left := num - size
				vs = string(append([]byte(vs), make([]byte, left, left)...))
			}

			if err := binary.Write(buf, order, []byte(vs)); err != nil {
				return err
			}
		}

		offset += n

		return nil
	})

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
