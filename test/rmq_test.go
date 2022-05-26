package testutils_test

import (
	"context"
	mathRand "math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wagslane/go-rabbitmq"

	"github.com/axelarnetwork/utils/slices"
	. "github.com/axelarnetwork/utils/test"
	"github.com/axelarnetwork/utils/test/rand"
)

func TestRMQFake(t *testing.T) {
	var (
		rmq          RMQFake
		handler      rabbitmq.Handler
		routingKeys  = []string{"routing_key_1", "routing_key_2"}
		msgCount     = 20
		expectedMsgs [][]byte
		receivedMsgs chan []byte
	)

	Given("an rmq fake", func() {
		rmq = NewRMQFake()
	}).
		Given("a msg handler", func() {
			receivedMsgs = make(chan []byte)
			handler = func(d rabbitmq.Delivery) (action rabbitmq.Action) {
				receivedMsgs <- d.Body
				return rabbitmq.Ack
			}
		},
		).
		When("subscribing to a queue", func() {
			rmq.StartConsuming(handler, "queue", routingKeys)
		}).
		When("publishing msgs to the queue", func() {
			expectedMsgs = slices.Expand(func(_ int) []byte { return rand.BytesBetween(1, 100) }, msgCount)
			for _, msg := range expectedMsgs {
				rmq.Publish(msg, []string{routingKeys[mathRand.Intn(2)]})
			}
		}).
		Then("consume all data from the queue", func(t *testing.T) {
			timeout, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			for i := 0; i < msgCount; i++ {
				select {
				case msg := <-receivedMsgs:
					assert.Equal(t, expectedMsgs[i], msg)
				case <-timeout.Done():
					assert.Fail(t, "timed out")
				}
			}
		}).
		Then("don't publish after fake is stopped", func(t *testing.T) {
			rmq.StopConsuming()
			rmq.Publish(rand.Bytes(2), routingKeys)

			timeout, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()
			select {
			case <-receivedMsgs:
				assert.Fail(t, "should not deliver another msg")
			case <-timeout.Done():
				return
			}
		}).
		Run(t)
}
