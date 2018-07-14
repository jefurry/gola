// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package binary

import (
	"bufio"
	"encoding/binary"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

func scanToken(s string, handler func(psym, c byte, num int) error) error {
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanBytes)

	idx := 0
	buf := make([]byte, 0, 3)
	psym := DefaultPrecursorSymbol
	for scanner.Scan() {
		s := scanner.Text()
		c := byte(s[0])

		// skip space
		if c == ' ' {
			continue
		}

		if idx == 0 {
			idx += 1
			if inPrecursorSymbols(c) {
				psym = c
				continue
			}
		}

		// parse number
		if c >= '0' && c <= '9' {
			buf = append(buf, c)
			continue
		}

		if !inAllowFormatSymbols(c) {
			return errors.Errorf("unknown format character for '%c'", c)
		}

		num := 1
		var err error
		if len(buf) > 0 {
			num, err = strconv.Atoi(string(buf))
			if err != nil {
				return err
			}

			if num < 0 {
				return errors.Errorf("invalid format character width for '%d", num)
			}
		}

		buf = buf[:0]

		if err := handler(psym, c, num); err != nil {
			return err
		}
	}

	return nil
}

func splitToken(s string) ([]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanBytes)

	idx := 0
	tokens := make([]string, 0, 5)
	buf := make([]byte, 0, 5)
	for scanner.Scan() {
		s := scanner.Text()
		c := byte(s[0])

		// skip space
		if c == ' ' {
			continue
		}

		// append precursor symbol
		// '@', '=', '<', '>', '!'
		if idx == 0 {
			idx += 1
			if !inPrecursorSymbols(c) {
				tokens = append(tokens, string(DefaultPrecursorSymbol))
			} else {
				tokens = append(tokens, string(c))
				continue
			}
		}

		// parse number
		if c >= '0' && c <= '9' {
			buf = append(buf, c)
			continue
		}

		if !inAllowFormatSymbols(c) {
			return nil, errors.Errorf("unknown format character for '%c'", c)
		}

		buf = append(buf, c)
		tokens = append(tokens, string(buf))
		buf = buf[:0]
	}

	if len(tokens) < 1 {
		return nil, errors.New("invalid format characters")
	}

	return tokens, nil
}

func inAllowFormatSymbols(c byte) bool {
	for _, r := range AllowFormatSymbols {
		if r == c {
			return true
		}
	}

	return false
}

func inPrecursorSymbols(c byte) bool {
	if c == AtPrecursorSymbol || c == EqualPrecursorSymbol ||
		c == LTPrecursorSymbol || c == GTPrecursorSymbol ||
		c == BangPrecursorSymbol {
		return true
	}

	return false
}

func getByteOrder(psym byte) binary.ByteOrder {
	switch psym {
	case AtPrecursorSymbol, EqualPrecursorSymbol:
		return nativeByteOrder
	case LTPrecursorSymbol:
		return binary.LittleEndian
	}

	return binary.BigEndian
}
