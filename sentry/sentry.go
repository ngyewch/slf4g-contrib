package sentry

import (
	"fmt"
	slf4g "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/level"
	"github.com/getsentry/sentry-go"
	"reflect"
	"time"
)

func ToSentryEvent(event slf4g.Event, source slf4g.CoreLogger) *sentry.Event {
	sentryEvent := sentry.NewEvent()
	timestampPopulated := false
	loggerPopulated := false
	messagePopulated := false
	errorPopulated := false

	fieldKeysSpec := source.GetProvider().GetFieldKeysSpec()
	extra := make(map[string]interface{})
	_ = event.ForEach(func(key string, value interface{}) error {
		if key == fieldKeysSpec.GetTimestamp() {
			t, ok := value.(time.Time)
			if ok {
				sentryEvent.Timestamp = t
				timestampPopulated = true
			}
		} else if key == fieldKeysSpec.GetLogger() {
			logger, ok := value.(string)
			if ok {
				sentryEvent.Logger = logger
				loggerPopulated = true
			}
		} else if key == fieldKeysSpec.GetMessage() {
			stringer, ok := value.(fmt.Stringer)
			if ok {
				sentryEvent.Message = stringer.String()
				messagePopulated = true
			}
		} else if key == fieldKeysSpec.GetError() {
			err, ok := value.(error)
			if ok {
				exception := sentry.Exception{
					Value:      err.Error(),
					Type:       reflect.TypeOf(err).String(),
					Stacktrace: newStacktrace(),
				}
				sentryEvent.Exception = append(sentryEvent.Exception, exception)
				errorPopulated = true
			}
		} else {
			extra[key] = value
		}
		return nil
	})
	if len(extra) > 0 {
		sentryEvent.Extra = extra
	}
	if !timestampPopulated {
		sentryEvent.Timestamp = time.Now()
	}
	if !loggerPopulated {
		sentryEvent.Logger = source.GetName()
	}
	if !messagePopulated {
		// do nothing
	}
	if !errorPopulated {
		// do nothing
	}

	switch event.GetLevel() {
	case level.Trace:
		sentryEvent.Level = sentry.LevelDebug
		break
	case level.Debug:
		sentryEvent.Level = sentry.LevelDebug
		break
	case level.Info:
		sentryEvent.Level = sentry.LevelInfo
		break
	case level.Warn:
		sentryEvent.Level = sentry.LevelWarning
		break
	case level.Error:
		sentryEvent.Level = sentry.LevelError
		break
	case level.Fatal:
		sentryEvent.Level = sentry.LevelFatal
		break
	}

	return sentryEvent
}

func newStacktrace() *sentry.Stacktrace {
	const (
		currentModule = "github.com/ngyewch/slf4g-contrib/sentry"
		slf4gModule   = "github.com/echocat/slf4g"
	)

	st := sentry.NewStacktrace()

	threshold := len(st.Frames) - 1
	// drop current module frames
	for ; threshold > 0 && st.Frames[threshold].Module == currentModule; threshold-- {
	}

outer:
	// try to drop slf4go module frames after logger call point
	for i := threshold; i > 0; i-- {
		if st.Frames[i].Module == slf4gModule {
			for j := i - 1; j >= 0; j-- {
				if st.Frames[j].Module != slf4gModule {
					threshold = j
					break outer
				}
			}
			break
		}
	}

	st.Frames = st.Frames[:threshold+1]

	return st
}
