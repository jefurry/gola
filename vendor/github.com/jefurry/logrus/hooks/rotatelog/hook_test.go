package rotatelog_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jefurry/logrus"
	rotatelog "github.com/jefurry/logrus/hooks/rotatelog"
	"github.com/stretchr/testify/assert"
)

func TestLogEntryWritten(t *testing.T) {
	dir, err := ioutil.TempDir("", "file-rotatelog-test")
	if !assert.NoError(t, err, "creating temporary directory should succeed") {
		return
	}
	defer os.RemoveAll(dir)

	log := logrus.New()
	hook, err := rotatelog.NewHook(filepath.Join(dir, "log.%Y%m%d%H%M%S"),
		rotatelog.WithLinkName(filepath.Join(dir, "log")),
		rotatelog.WithMaxAge(24*time.Hour),
		rotatelog.WithRotationTime(time.Hour),
		rotatelog.WithClock(rotatelog.UTC))

	if err != nil {
		t.Errorf(err.Error())
	}

	log.Hooks.Add(hook.SetFormatter(&logrus.JSONFormatter{}))

	for _, level := range hook.Levels() {
		if len(log.Hooks[level]) != 1 {
			t.Errorf("RotatelogHook was not added. The length of log.Hooks[%v]: %v", level, len(log.Hooks[level]))
		}
	}

	log.Info("Congratulations!")
}
