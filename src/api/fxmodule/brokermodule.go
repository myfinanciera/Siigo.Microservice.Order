package fxmodule

import (
	kafkaprovider "dev.azure.com/SiigoDevOps/Siigo/_git/go-slim.git/kafka"
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-slim.git/slim"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"siigo.com/order/src/application/consumer"
	"siigo.com/order/src/domain/order"
)

// BrokerModule Create FX Broker Kafka Module
//modulo fx con dependencias
var BrokerModule = fx.Options(
	fx.Provide(
		NewMessageBus,
	),
	fx.Invoke(),
)

func NewMessageBus(v *viper.Viper) *slim.MessageBusBuilder {
	nametopic := v.GetString("domain.topic")

	if len(nametopic) == 0 {
		panic("Topic Name doesnt set")
	}

	return slim.NewMessageBusBuilder().

		// Inicialice kafka provider
		WithProviderKafka(kafkaprovider.NewKafkaConfig(v)).

		// register consumers
		WithConsumer(nametopic, consumer.TestConsumer).

		// register producers
		WithProduce(kafkaprovider.
			NewProducerBuilder().
			WithTopic(nametopic).
			WithType(&order.OrderCreatedDomainEvent{}). //event type
			WithPartitionKey(func(event interface{}) []byte {
				return []byte(event.(*order.OrderCreatedDomainEvent).Order.Id)//garantizar el orden de eventos en kafka
			}).
			Build(),
		).

		// add ssl
		//WithSsl().

		// start up
		Build()
}
