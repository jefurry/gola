// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package event implements event manager.
package event

type (
	Event struct {
		Data    interface{}
		Context interface{}
	}
)

func NewEvent(data interface{}, cxts ...interface{}) *Event {
	evt := &Event{
		Data: data,
	}

	if len(cxts) > 0 {
		evt.Context = cxts[0]
	}

	return evt
}
