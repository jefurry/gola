// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package timelib implements time for Lua.
package timelib

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"time"
)

const (
	TimeLibName = "time"
)

func Open(L *lua.LState) {
	L.PreloadModule(TimeLibName, Loader)
}

func Loader(L *lua.LState) int {
	dimod := L.SetFuncs(L.NewTable(), timeFuncs)
	L.Push(dimod)

	timeRegisterTimeMetatype(L)
	timeRegisterLocationMetatype(L)

	dimod.RawSetString("UTC", newLocation(L, time.UTC))
	dimod.RawSetString("Local", newLocation(L, time.Local))

	for k, v := range timeFields {
		dimod.RawSetString(k, v)
	}

	return 1
}

var timeFuncs = map[string]lua.LGFunction{
	"parse":                  timeParse,
	"parseInLocation":        timeParseInLocation,
	"parseDuration":          timeParseDuration,
	"fixedZone":              timeFixedZone,
	"loadLocation":           timeLoadLocation,
	"loadLocationFromTZData": timeLoadLocationFromTZData,
	"now":    timeNow,
	"date":   timeDate,
	"since":  timeSince,
	"until":  timeUntil,
	"unix":   timeUnix,
	"sleep":  timeSleep,
	"isLeap": timeIsLeap,
}

var timeFields = map[string]lua.LValue{
	"ANSIC":       lua.LString(time.ANSIC),
	"UnixDate":    lua.LString(time.UnixDate),
	"RubyDate":    lua.LString(time.RubyDate),
	"RFC822":      lua.LString(time.RFC822),
	"RFC822Z":     lua.LString(time.RFC822Z),
	"RFC850":      lua.LString(time.RFC850),
	"RFC1123":     lua.LString(time.RFC1123),
	"RFC1123Z":    lua.LString(time.RFC1123Z),
	"RFC3339":     lua.LString(time.RFC3339),
	"RFC3339Nano": lua.LString(time.RFC3339Nano),
	"Kitchen":     lua.LString(time.Kitchen),
	// Handy time stamps.
	"Stamp":      lua.LString(time.Stamp),
	"StampMilli": lua.LString(time.StampMilli),
	"StampMicro": lua.LString(time.StampMicro),
	"StampNano":  lua.LString(time.StampNano),

	// time.Duration
	"Nanosecond":  lua.LNumber(time.Nanosecond),
	"Microsecond": lua.LNumber(time.Microsecond),
	"Millisecond": lua.LNumber(time.Millisecond),
	"Second":      lua.LNumber(time.Second),
	"Minute":      lua.LNumber(time.Minute),
	"Hour":        lua.LNumber(time.Hour),
}

func timeNow(L *lua.LState) int {
	t := time.Now()

	L.Push(newTime(L, t))

	return 1
}

func timeDate(L *lua.LState) int {
	tb := L.OptTable(1, nil)
	if tb == nil {
		tb = L.CreateTable(0, 0)
	}

	year := tb.RawGetString("year")
	if year == lua.LNil {
		year = lua.LNumber(1970)
	}

	month := tb.RawGetString("month")
	if month == lua.LNil {
		month = lua.LNumber(1)
	}

	day := tb.RawGetString("day")
	if day == lua.LNil {
		day = lua.LNumber(1)
	}

	hour := tb.RawGetString("hour")
	if hour == lua.LNil {
		hour = lua.LNumber(0)
	}

	min := tb.RawGetString("min")
	if min == lua.LNil {
		min = lua.LNumber(0)
	}

	sec := tb.RawGetString("sec")
	if sec == lua.LNil {
		sec = lua.LNumber(0)
	}

	nsec := tb.RawGetString("nsec")
	if nsec == lua.LNil {
		nsec = lua.LNumber(0)
	}

	loc := tb.RawGetString("loc")
	if loc == lua.LNil {
		loc = newLocation(L, time.UTC)
	}

	if year.Type() != lua.LTNumber {
		L.ArgError(1, fmt.Sprintf("%s expected, got %s", lua.LTNumber, year.Type()))
	}
	y, _ := year.(lua.LNumber)

	if month.Type() != lua.LTNumber {
		L.ArgError(1, fmt.Sprintf("%s expected, got %s", lua.LTNumber, month.Type()))
	}
	mon, _ := month.(lua.LNumber)

	if day.Type() != lua.LTNumber {
		L.ArgError(1, fmt.Sprintf("%s expected, got %s", lua.LTNumber, day.Type()))
	}
	d, _ := day.(lua.LNumber)

	if hour.Type() != lua.LTNumber {
		L.ArgError(1, fmt.Sprintf("%s expected, got %s", lua.LTNumber, hour.Type()))
	}
	h, _ := hour.(lua.LNumber)

	if min.Type() != lua.LTNumber {
		L.ArgError(1, fmt.Sprintf("%s expected, got %s", lua.LTNumber, min.Type()))
	}
	m, _ := min.(lua.LNumber)

	if sec.Type() != lua.LTNumber {
		L.ArgError(1, fmt.Sprintf("%s expected, got %s", lua.LTNumber, sec.Type()))
	}
	s, _ := sec.(lua.LNumber)

	if nsec.Type() != lua.LTNumber {
		L.ArgError(1, fmt.Sprintf("%s expected, got %s", lua.LTNumber, nsec.Type()))
	}
	ns, _ := nsec.(lua.LNumber)

	if loc.Type() != lua.LTUserData {
		L.ArgError(1, fmt.Sprintf("%s expected, got %s", lua.LTNumber, loc.Type()))
	}

	ud, _ := loc.(*lua.LUserData)
	l, ok := ud.Value.(*time.Location)
	if !ok {
		L.ArgError(1, "invalid time zone")
	}

	t := time.Date(int(y), time.Month(int(mon)), int(d), int(h), int(m), int(s), int(ns), l)

	L.Push(newTime(L, t))

	return 1
}

func timeSince(L *lua.LState) int {
	t := checkTime(L, 1)

	d := time.Since(t)

	L.Push(lua.LNumber(d))

	return 1
}

func timeUntil(L *lua.LState) int {
	t := checkTime(L, 1)

	d := time.Until(t)

	L.Push(lua.LNumber(d))

	return 1
}

func timeUnix(L *lua.LState) int {
	sec := L.CheckInt(1)
	nsec := L.CheckInt(2)

	u := time.Unix(int64(sec), int64(nsec))

	L.Push(newTime(L, u))

	return 1
}

func timeSleep(L *lua.LState) int {
	d := L.CheckInt(1)

	time.Sleep(time.Duration(d))

	return 0
}

func timeIsLeap(L *lua.LState) int {
	year := L.CheckInt(1)

	L.Push(lua.LBool(isLeap(year)))

	return 1
}
