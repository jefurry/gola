# Rotatelog Hooks for Logrus <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>

## Usage

```go
import (
  "time"

  "github.com/jefurry/logrus"
  rlog "github.com/jefurry/logrus/hooks/rotatelog"
)

func main() {
  log := logrus.New()
  hook, err := rlog.NewHook("./access_log.%Y%m%d",
    //rotatelog.WithLinkName("./access_log"),
    rotatelog.WithMaxAge(24 * time.Hour),
    rotatelog.WithRotationTime(time.Hour),
    rotatelog.WithClock(rotatelog.UTC))

  if err != nil {
    log.Hooks.Add(hook)
  }
}
```

# DESCRIPTION

When you integrate this to to you app, it automatically write to logs that
are rotated from within the app: No more disk-full alerts because you forgot
to setup logrotate!

To install, simply issue a `go get`:

```
go get github.com/jefurry/logrus/hooks/rotatelog
```

OPTIONS
====

## Pattern (Required)

The pattern used to generate actual log file names. You should use patterns
using the strftime (3) format. For example:

```go
  rotatelog.New("/var/log/myapp/log.%Y%m%d")
```

## Clock (default: rotatelog.Local)

You may specify an object that implements the roatatelogs.Clock interface.
When this option is supplied, it's used to determine the current time to
base all of the calculations on. For example, if you want to base your
calculations in UTC, you may specify rotatelog.UTC

```go
  rotatelog.New(
    "/var/log/myapp/log.%Y%m%d",
    rotatelog.WithClock(rotatelog.UTC),
  )
```

## Location

This is an alternative to the `WithClock` option. Instead of providing an
explicit clock, you can provide a location for you times. We will create
a Clock object that produces times in your specified location, and configure
the rotatelog to respect it.

## LinkName (default: "")

Path where a symlink for the actual log file is placed. This allows you to 
always check at the same location for log files even if the logs were rotated

```go
  rotatelog.New(
    "/var/log/myapp/log.%Y%m%d",
    rotatelog.WithLinkName("/var/log/myapp/current"),
  )
```

```
  // Else where
  $ tail -f /var/log/myapp/current
```

If not provided, no link will be written.

## RotationTime (default: 86400 sec)

Interval between file rotation. By default logs are rotated every 86400 seconds.
Note: Remember to use time.Duration values.

```go
  // Rotate every hour
  rotatelog.New(
    "/var/log/myapp/log.%Y%m%d",
    rotatelog.WithRotationTime(time.Hour),
  )
```

## MaxAge (default: 7 days)

Time to wait until old logs are purged. By default no logs are purged, which
certainly isn't what you want.
Note: Remember to use time.Duration values.

```go
  // Purge logs older than 1 hour
  rotatelog.New(
    "/var/log/myapp/log.%Y%m%d",
    rotatelog.WithMaxAge(time.Hour),
  )
```

## RotationCount (default: -1)

The number of files should be kept. By default, this option is disabled.

Note: MaxAge should be disabled by specifing `WithMaxAge(-1)` explicitly.

```go
  // Purge logs except latest 7 files
  rotatelog.New(
    "/var/log/myapp/log.%Y%m%d",
    rotatelog.WithMaxAge(-1),
    rotatelog.WithRotationCount(7),
  )
```

## Handler (default: nil)

Sets the event handler to receive event notifications from the RotateLog
object. Currently only supported event type is FiledRotated

```go
  rotatelog.New(
    "/var/log/myapp/log.%Y%m%d",
    rotatelog.Handler(rotatelog.HandlerFunc(func(e Event) {
      if e.Type() != rotatelog.FileRotatedEventType {
        return
      }

      // Do what you want with the data. This is just an idea:
      storeLogFileToRemoteStorage(e.(*FileRotatedEvent).PreviousFile())
    })),
  )
```

# Rotating files forcefully

If you want to rotate files forcefully before the actual rotation time has reached,
you may use the `Rotate()` method. This method forcefully rotates the logs, but
if the generated file name clashes, then a numeric suffix is added so that
the new file will forcefully appear on disk.

For example, suppose you had a pattern of '%Y.log' with a rotation time of
`86400` so that it only gets rotated every year, but for whatever reason you
wanted to rotate the logs now, you could install a signal handler to
trigger this rotation:

```go
rl := rotatelog.New(...)

signal.Notify(ch, syscall.SIGHUP)

go func(ch chan os.Signal) {
  <-ch
  rl.Rotate()
}()
```

And you will get a log file name in like `2018.log.1`, `2018.log.2`, etc.


## Other

Thanks [lestrrat-go](https://github.com/lestrrat-go/file-rotatelogs).
