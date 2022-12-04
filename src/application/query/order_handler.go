package query

import (
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs"
	"errors"
	"reflect"
)

// Handle processes order commands.
func (handler *OrderQueryHandler) Handle(message cqrs.RequestMessage) (interface{}, error) {
	fn, ok := handler.handlersByStructType[reflect.TypeOf(message.Request()).String()]
	if !ok {
		return nil, errors.New("query type font found")
	}
	return fn(handler, message)
}

func LoadAllOrderQueryHandle(handler *OrderQueryHandler, message cqrs.RequestMessage) (interface{}, error) {
	orderResult, err := handler.orderService.LoadAll()
	if err != nil {
		return nil, err
	}
	return orderResult, nil

}

func LoadOrderQueryHandle(handler *OrderQueryHandler, message cqrs.RequestMessage) (interface{}, error) {

	cmd := message.Request().(*LoadOrderQuery)

	orderResult, err := handler.orderService.Get(cmd.Id)
	if err != nil {
		return nil, err
	}

	return orderResult, nil

}
