package airbrake // import "github.com/jefurry/logrus/hooks/airbrake"

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/airbrake/gobrake"
	"github.com/jefurry/logrus"
)

// AirbrakeHook to send exceptions to an exception-tracking service compatible
// with the Airbrake API.
type AirbrakeHook struct {
	Airbrake *gobrake.Notifier
}

func NewHook(projectID int64, apiKey, env string) *AirbrakeHook {
	airbrake := gobrake.NewNotifier(projectID, apiKey)
	airbrake.AddFilter(func(notice *gobrake.Notice) *gobrake.Notice {
		if env == "development" {
			return nil
		}
		notice.Context["environment"] = env
		return notice
	})
	hook := &AirbrakeHook{
		Airbrake: airbrake,
	}
	return hook
}

func (hook *AirbrakeHook) Fire(entry *logrus.Entry) error {
	var notifyErr error
	err, ok := entry.Data["error"].(error)
	if ok {
		notifyErr = err
	} else {
		notifyErr = errors.New(entry.Message)
	}
	var req *http.Request
	for k, v := range entry.Data {
		if r, ok := v.(*http.Request); ok {
			req = r
			delete(entry.Data, k)
			break
		}
	}
	notice := hook.Airbrake.Notice(notifyErr, req, 3)
	for k, v := range entry.Data {
		notice.Context[k] = fmt.Sprintf("%s", v)
	}

	hook.sendNotice(notice)
	return nil
}

func (hook *AirbrakeHook) sendNotice(notice *gobrake.Notice) {
	if _, err := hook.Airbrake.SendNotice(notice); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to send error to Airbrake: %v\n", err)
	}
}

func (hook *AirbrakeHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

func (hook *AirbrakeHook) Close() error {
	return nil
}
