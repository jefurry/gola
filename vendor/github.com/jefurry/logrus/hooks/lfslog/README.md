# Local Filesystem Hook for Logrus

[![GoDoc](https://godoc.org/github.com/jefurry/logrus/hooks/lfslog?status.svg)](http://godoc.org/github.com/jefurry/logrus/hooks/lfslog)

Sometimes developers like to write directly to a file on the filesystem. This is a hook for [`logrus`](https://github.com/jefurry/logrus) which designed to allow users to do that. The log levels are dynamic at instantiation of the hook, so it is capable of logging at some or all levels.

## Example
```go
import (
	"github.com/jefurry/logrus/hooks/lfslog"
	"github.com/jefurry/logrus"
)

var Log *logrus.Logger

func NewLogger() *logrus.Logger {
	if Log != nil {
		return Log
	}

	pathMap := lfslog.PathMap{
		logrus.InfoLevel:  "/var/log/info.log",
		logrus.ErrorLevel: "/var/log/error.log",
	}

	Log = logrus.New()
	Log.Hooks.Add(lfslog.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))
	return Log
}
```

### Formatters
`lfslog` will strip colors from any `TextFormatter` type formatters when writing to local file, because the color codes don't look great in file.

If no formatter is provided via `lfslog.NewHook`, a default text formatter will be used.

### Log rotation
In order to enable automatic log rotation it's possible to provide an io.Writer instead of the path string of a log file.
In combination with packages like [file-rotatelogs](https://github.com/lestrrat-go/file-rotatelogs) log rotation can easily be achieved.

```go
package main

import (
	"github.com/jefurry/logrus/hooks/lfslog"
	"github.com/jefurry/logrus"
	rotatelog "github.com/jefurry/logrus/hooks/rotatelog"
)

var Log *logrus.Logger

func NewLogger() (*logrus.Logger, error) {
	if Log != nil {
		return Log
	}

	path := "/var/log/go.log"
	writer, err := rotatelog.New(
		path+".%Y%m%d%H%M",
		rotatelog.WithLinkName(path),
		rotatelog.WithMaxAge(time.Duration(86400)*time.Second),
		rotatelog.WithRotationTime(time.Duration(604800)*time.Second),
	)
	if err != nil {
		return nil, err
	}

	Log1 = logrus.New()
	Log1.Hooks.Add(lfslog.NewHook(
		lfslog.WriterMap{
			logrus.InfoLevel:  writer,
			logrus.ErrorLevel: writer,
		},
		&logrus.JSONFormatter{},
	))

	Log = logrus.New()
	Log.Hooks.Add(lfslog.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))

	return Log, nil
}
```

### Note:
User who run the go application must have read/write permissions to the selected log files. If the files do not exists yet, then user must have permission to the target directory.


## Other

Thanks [rifflock](https://github.com/rifflock/lfshook).

