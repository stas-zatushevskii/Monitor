package inmemorystorage

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"
)

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewProducer(fileName string) (*Producer, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Producer) WriteEvent(event *InMemoryStorage) error {
	return p.encoder.Encode(&event)
}

func NewConsumer(fileName string) (*Consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *Consumer) ReadEvent() (*InMemoryStorage, error) {
	event := &InMemoryStorage{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}

	return event, nil
}

func AutoSaveData(ctx context.Context, storage *InMemoryStorage, reportInterval int, filename string) {
	ticker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	producer, err := NewProducer(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.file.Close()
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			snapshot := storage.Snapshot()
			err = producer.WriteEvent(&snapshot)
			if err != nil {
				log.Fatal(err)
			}
		case <-ctx.Done():
			snapshot := storage.Snapshot()
			err = producer.WriteEvent(&snapshot)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func AutoLoadData(filename string, storage *InMemoryStorage) error {
	consumer, err := NewConsumer(filename)
	if err != nil {
		return err
	}
	defer consumer.file.Close()
	err = consumer.decoder.Decode(&storage)
	if err != nil {
		return err
	}
	return nil
}
