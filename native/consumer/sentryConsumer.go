package consumer

import (
	"fmt"
	slf4g "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/level"
	sentry2 "github.com/ngyewch/slf4g-contrib/sentry"
)

type SentryConsumer struct {
	level level.Level
}

func NewSentryConsumer(level level.Level) *SentryConsumer {
	return &SentryConsumer{
		level: level,
	}
}

func (consumer *SentryConsumer) Consume(event slf4g.Event, source slf4g.CoreLogger) {
	if event.GetLevel().CompareTo(consumer.level) < 0 {
		return
	}

	sentryEvent := sentry2.ToSentryEvent(event, source)
	fmt.Printf("sentryEvent=%+v\n", sentryEvent)
	//sentry.CaptureEvent(sentryEvent)
}
