package rotatelog

import (
	"fmt"
	"os"

	"github.com/jefurry/logrus"
)

// SyslogHook to send logs via syslog.
type RotatelogHook struct {
	w         *RotateLog
	formatter logrus.Formatter
}

// Creates a hook to be added to an instance of logger. This is called with
// `hook, err := NewHook("./access_log.%Y%m%d",
//									WithLinkName("./access_log"),
// 									WithMaxAge(24 * time.Hour),
//									WithRotationTime(time.Hour),
//									WithClock(UTC))`
// `if err == nil { log.Hooks.Add(hook.SetFormatter(&logrus.JSONFormatter{})) }`
func NewHook(p string, options ...Option) (*RotatelogHook, error) {
	w, err := New(p, options...)
	if err != nil {
		return nil, err
	}

	return &RotatelogHook{w: w}, nil
}

func (hook *RotatelogHook) Fire(entry *logrus.Entry) error {
	var line []byte

	if hook.formatter != nil {
		msg, err := hook.formatter.Format(entry)
		if err != nil {
			return err
		}

		line = msg
	} else {
		msg, err := entry.String()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
			return err
		}

		line = []byte(msg)
	}

	_, err := hook.w.Write(line)

	return err
}

func (hook *RotatelogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *RotatelogHook) SetFormatter(formatter logrus.Formatter) *RotatelogHook {
	hook.formatter = formatter

	return hook
}

func (hook *RotatelogHook) Close() error {
	return hook.w.Close()
}
