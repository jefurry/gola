// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package timelib

import (
	"github.com/yuin/gopher-lua"
	"time"
)

const (
	timeTimeTypeName = TimeLibName + ".TIME*"
)

func timeTimeAfter(L *lua.LState) int {
	t := checkTime(L, 1)
	u := checkTime(L, 2)

	L.Push(lua.LBool(t.After(u)))

	return 1
}

func timeTimeBefore(L *lua.LState) int {
	t := checkTime(L, 1)
	u := checkTime(L, 2)

	L.Push(lua.LBool(t.Before(u)))

	return 1
}

func timeTimeEqual(L *lua.LState) int {
	t := checkTime(L, 1)
	u := checkTime(L, 2)

	L.Push(lua.LBool(t.Equal(u)))

	return 1
}

func timeTimeIsZero(L *lua.LState) int {
	t := checkTime(L, 1)

	L.Push(lua.LBool(t.IsZero()))

	return 1
}

func timeTimeDate(L *lua.LState) int {
	t := checkTime(L, 1)

	year, month, day := t.Date()

	L.Push(lua.LNumber(year))
	L.Push(lua.LNumber(month))
	L.Push(lua.LNumber(day))

	return 3
}

func timeTimeYear(L *lua.LState) int {
	t := checkTime(L, 1)

	L.Push(lua.LNumber(t.Year()))

	return 1
}

func timeTimeMonth(L *lua.LState) int {
	t := checkTime(L, 1)

	L.Push(lua.LNumber(t.Month()))

	return 1
}

func timeTimeDay(L *lua.LState) int {
	t := checkTime(L, 1)

	L.Push(lua.LNumber(t.Day()))

	return 1
}

func timeTimeWeekday(L *lua.LState) int {
	t := checkTime(L, 1)

	L.Push(lua.LNumber(t.Weekday()))

	return 1
}

func timeTimeISOWeek(L *lua.LState) int {
	t := checkTime(L, 1)

	year, week := t.ISOWeek()

	L.Push(lua.LNumber(year))
	L.Push(lua.LNumber(week))

	return 2
}

func timeTimeClock(L *lua.LState) int {
	t := checkTime(L, 1)

	hour, minute, second := t.Clock()

	L.Push(lua.LNumber(hour))
	L.Push(lua.LNumber(minute))
	L.Push(lua.LNumber(second))

	return 3
}

func timeTimeHour(L *lua.LState) int {
	t := checkTime(L, 1)

	L.Push(lua.LNumber(t.Hour()))

	return 1
}

func timeTimeMinute(L *lua.LState) int {
	t := checkTime(L, 1)

	L.Push(lua.LNumber(t.Minute()))

	return 1
}

func timeTimeSecond(L *lua.LState) int {
	t := checkTime(L, 1)

	L.Push(lua.LNumber(t.Second()))

	return 1
}

func timeTimeNanosecond(L *lua.LState) int {
	t := checkTime(L, 1)

	L.Push(lua.LNumber(t.Nanosecond()))

	return 1
}

func timeTimeYearDay(L *lua.LState) int {
	t := checkTime(L, 1)

	L.Push(lua.LNumber(t.YearDay()))

	return 1
}

func timeTimeAdd(L *lua.LState) int {
	t := checkTime(L, 1)
	d := L.CheckInt(2)

	L.Push(newTime(L, t.Add(time.Duration(d))))

	return 1
}

func timeTimeSub(L *lua.LState) int {
	t := checkTime(L, 1)
	u := checkTime(L, 2)

	d := t.Sub(u)

	L.Push(lua.LNumber(d))

	return 1
}

func timeTimeAddDate(L *lua.LState) int {
	t := checkTime(L, 1)
	year := L.CheckInt(2)
	month := L.CheckInt(3)
	day := L.CheckInt(4)

	u := t.AddDate(year, month, day)

	L.Push(newTime(L, u))

	return 1
}

func timeTimeUTC(L *lua.LState) int {
	t := checkTime(L, 1)

	u := t.UTC()

	L.Push(newTime(L, u))

	return 1
}

func timeTimeLocal(L *lua.LState) int {
	t := checkTime(L, 1)

	u := t.Local()

	L.Push(newTime(L, u))

	return 1
}

func timeTimeIn(L *lua.LState) int {
	t := checkTime(L, 1)
	loc := checkLocation(L, 2)

	u := t.In(loc)

	L.Push(newTime(L, u))

	return 1
}

func timeTimeLocation(L *lua.LState) int {
	t := checkTime(L, 1)

	loc := t.Location()

	L.Push(newLocation(L, loc))

	return 1
}

func timeTimeZone(L *lua.LState) int {
	t := checkTime(L, 1)

	name, offset := t.Zone()

	L.Push(lua.LString(name))
	L.Push(lua.LNumber(offset))

	return 2
}

func timeTimeUnix(L *lua.LState) int {
	t := checkTime(L, 1)

	sec := t.Unix()

	L.Push(lua.LNumber(sec))

	return 1
}

func timeTimeUnixNano(L *lua.LState) int {
	t := checkTime(L, 1)

	nsec := t.UnixNano()

	L.Push(lua.LNumber(nsec))

	return 1
}

func timeTimeTruncate(L *lua.LState) int {
	t := checkTime(L, 1)
	d := L.CheckInt(2)

	u := t.Truncate(time.Duration(d))

	L.Push(newTime(L, u))

	return 1
}

func timeTimeRound(L *lua.LState) int {
	t := checkTime(L, 1)
	d := L.CheckInt(2)

	u := t.Round(time.Duration(d))

	L.Push(newTime(L, u))

	return 1
}

func timeTimeString(L *lua.LState) int {
	t := checkTime(L, 1)

	L.Push(lua.LString(t.String()))

	return 1
}

func timeRegisterTimeMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(timeTimeTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), timeTimeFuncs))
	L.SetField(mt, "__tostring", L.NewFunction(timeTimeString))
}

var timeTimeFuncs = map[string]lua.LGFunction{
	"after":      timeTimeAfter,
	"before":     timeTimeBefore,
	"equal":      timeTimeEqual,
	"isZero":     timeTimeIsZero,
	"date":       timeTimeDate,
	"year":       timeTimeYear,
	"month":      timeTimeMonth,
	"day":        timeTimeDay,
	"weekday":    timeTimeWeekday,
	"isoWeek":    timeTimeISOWeek,
	"clock":      timeTimeClock,
	"hour":       timeTimeHour,
	"minute":     timeTimeMinute,
	"second":     timeTimeSecond,
	"nanosecond": timeTimeNanosecond,
	"yearDay":    timeTimeYearDay,
	"add":        timeTimeAdd,
	"sub":        timeTimeSub,
	"addDate":    timeTimeAddDate,
	"UTC":        timeTimeUTC,
	"Local":      timeTimeLocal,
	"In":         timeTimeIn,
	"location":   timeTimeLocation,
	"zone":       timeTimeZone,
	"unix":       timeTimeUnix,
	"unixNano":   timeTimeUnixNano,
	"truncate":   timeTimeTruncate,
	"round":      timeTimeRound,
}
