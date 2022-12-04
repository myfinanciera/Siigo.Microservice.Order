package fxmodule

import (
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs"
	slim "dev.azure.com/SiigoDevOps/Siigo/_git/go-slim.git/slim"
	"go.uber.org/fx"
	header "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"siigo.com/order/src/api/fxmodule/interceptor"
	"siigo.com/order/src/domain/services"
)
//modulo fx con dependencias
var CQRSDDDModule = fx.Options(
	fx.Provide(
		/*Register Event Dispatcher*/
		cqrs.NewInternalEventBus,
		cqrs.NewEventStoreClient,
		cqrs.NewEventStoreCommonDomainRepo,
		cqrs.NewInMemoryRepo,
		NewDomainEventsDispatcher,
		interceptor.NewInMemoryDispatcherWithInterceptors,

		/*Register domain services */
		services.NewOrderService,
	),
)

func NewDomainEventsDispatcher(eventStore *cqrs.GetEventStoreCommonDomainRepo, slim *slim.MessageBusBuilder) *cqrs.DomainEventsDispatcher {
	return &cqrs.DomainEventsDispatcher{
		EventStore: eventStore,
		EventPublisher: func(event interface{}, headers map[string]string) {
			if event == nil {
				return
			}
			var headersarr = []header.Header{}
			for key, item := range headers {
				head := header.Header{}
				head.Key = key
				head.Value = []byte(item)
				headersarr = append(headersarr, head)
			}
			slim.PublishWithHeader(event, headersarr)
		},
	}
}
