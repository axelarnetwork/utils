package testutils

import (
	"sync"

	"github.com/rabbitmq/amqp091-go"
	"github.com/wagslane/go-rabbitmq"
	"golang.org/x/exp/slices"

	slices2 "github.com/axelarnetwork/utils/slices"
)

// RMQFake is a fake RMQ broker that acts both as publisher and consumer
type RMQFake struct {
	bindings map[string][]string
	queues   map[string]chan rabbitmq.Delivery
	done     bool
	stopLock *sync.RWMutex
}

// NewRMQFake returns a new RMQFake
func NewRMQFake() RMQFake {
	return RMQFake{
		bindings: make(map[string][]string),
		queues:   make(map[string]chan rabbitmq.Delivery),
		done:     false,
		stopLock: &sync.RWMutex{},
	}
}

// StartConsuming creates a new goroutine where the handler consumes msgs from the queue. The queue is created if it doesn't exist
func (fake *RMQFake) StartConsuming(handler rabbitmq.Handler, queue string, routingKeys []string, _ ...func(options *rabbitmq.ConsumeOptions)) {
	for _, key := range routingKeys {
		queues := fake.bindings[key]
		if !slices.Contains(queues, queue) {
			fake.bindings[key] = append(queues, queue)
		}
	}

	if _, ok := fake.queues[queue]; !ok {
		fake.queues[queue] = make(chan rabbitmq.Delivery, 10000)
	}

	go func() {
		deliveries := fake.queues[queue]
		for delivery := range deliveries {
			switch handler(delivery) {
			case rabbitmq.Ack:
				continue
			case rabbitmq.NackDiscard:
				continue
			case rabbitmq.NackRequeue:
				deliveries <- delivery
			}
		}
	}()
}

// Publish publishes the given data to the queues bound to the routing keys. Msg is lost of no queue is bound to the keys
func (fake *RMQFake) Publish(data []byte, routingKeys []string, optionFuncs ...func(options *rabbitmq.PublishOptions)) {
	fake.stopLock.RLock()
	defer fake.stopLock.RUnlock()

	if fake.done {
		return
	}

	queues := slices2.Distinct(
		slices2.Reduce(routingKeys, nil,
			func(queues []string, key string) []string {
				return append(queues, fake.bindings[key]...)
			}))

	options := &rabbitmq.PublishOptions{}
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}
	if options.DeliveryMode == 0 {
		options.DeliveryMode = rabbitmq.Transient
	}

	delivery := rabbitmq.Delivery{
		Delivery: amqp091.Delivery{
			ContentType:     options.ContentType,
			DeliveryMode:    options.DeliveryMode,
			Body:            data,
			Headers:         amqp091.Table(options.Headers),
			Expiration:      options.Expiration,
			ContentEncoding: options.ContentEncoding,
			Priority:        options.Priority,
			CorrelationId:   options.CorrelationID,
			ReplyTo:         options.ReplyTo,
			MessageId:       options.MessageID,
			Timestamp:       options.Timestamp,
			Type:            options.Type,
			UserId:          options.UserID,
			AppId:           options.AppID,
		}}

	for _, queue := range queues {
		fake.queues[queue] <- delivery
	}
}

// StopConsuming stops msg deliveries to queues
func (fake *RMQFake) StopConsuming() {
	fake.stopLock.Lock()
	defer fake.stopLock.Unlock()
	for _, deliveries := range fake.queues {
		close(deliveries)
	}
	fake.done = true
}
