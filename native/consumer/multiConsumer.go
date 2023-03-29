package consumer

import (
	slf4g "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/native/consumer"
)

type MultiConsumer struct {
	consumers []consumer.Consumer
}

func NewMultiConsumer(consumers ...consumer.Consumer) *MultiConsumer {
	return &MultiConsumer{
		consumers: consumers,
	}
}

func (consumer *MultiConsumer) Consume(event slf4g.Event, source slf4g.CoreLogger) {
	for _, delegate := range consumer.consumers {
		delegate.Consume(event, source)
	}
}
