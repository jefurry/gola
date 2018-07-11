// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package loglib

import (
	"github.com/jefurry/logrus"
	"github.com/yuin/gopher-lua"
)

var (
	logDiscard = lua.LString(`\0`)

	// These are the different logging levels. You can set the logging level to log
	// on your instance of logger, obtained with `logger:new()`.
	logFields = map[string]lua.LValue{
		// PanicLevel level, highest level of severity. Logs and then calls panic with the
		// message passed to Debug, Info, ...
		"PanicLevel": lua.LNumber(logrus.PanicLevel),

		// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
		// logging level is set to Panic.)
		"FatalLevel": lua.LNumber(logrus.FatalLevel),

		// ErrorLevel level. Logs. Used for errors that should definitely be noted.
		// Commonly used for hooks to send errors to an error tracking service.
		"ErrorLevel": lua.LNumber(logrus.ErrorLevel),

		// WarnLevel level. Non-critical entries that deserve eyes.
		"WarnLevel": lua.LNumber(logrus.WarnLevel),

		// InfoLevel level. General operational entries about what's going on inside the
		// application.
		"InfoLevel": lua.LNumber(logrus.InfoLevel),

		// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
		"DebugLevel": lua.LNumber(logrus.DebugLevel),

		"Discard": logDiscard,

		"ErrorKey": lua.LString(logrus.ErrorKey),
	}
)
