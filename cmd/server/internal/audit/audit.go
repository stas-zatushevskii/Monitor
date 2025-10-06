package audit

import (
	"log"

	"github.com/stas-zatushevskii/Monitor/cmd/server/config"
)

type Publisher interface {
	Register(subscriber Subscriber)
	NotifyAll(msg string)
}

type Subscriber interface {
	ReactToPublisherMsg(msg []byte)
}

type LogProducer struct {
	SubscriberList []Subscriber
}

func (s *LogProducer) Register(subscriber Subscriber) {
	s.SubscriberList = append(s.SubscriberList, subscriber)
}

func (s *LogProducer) NotifyAll(msg []byte) {
	for _, subscriber := range s.SubscriberList {
		subscriber.ReactToPublisherMsg(msg)
	}
}

type LogConsumer struct {
	Config config.Config
	Logger log.Logger
}

func NewLogConsumer(config config.Config) *LogConsumer {
	return &LogConsumer{Config: config}
}

func (l *LogConsumer) ReactToPublisherMsg(msg []byte) {
	if l.Config.Audit.URL != "" {
		if err := SendToURL(l.Config.Audit.URL, msg); err != nil {
			l.Logger.Fatal(err)
		}
	}
	if l.Config.Audit.FilePath != "" {
		if err := SaveToFile(l.Config.Audit.FilePath, msg); err != nil {
			l.Logger.Fatal(err)
		}
	}
}
