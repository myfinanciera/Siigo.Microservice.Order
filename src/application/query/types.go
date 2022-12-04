package query

import (
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs"
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs/uuid"
	"reflect"
	"siigo.com/order/src/domain/services"
)

//LoadAllOrder
type LoadAllOrderQuery struct {
}

// LoadOrderQuery create a new order item
type LoadOrderQuery struct {
	Id uuid.UUID
}

type OrderQueryHandler struct {
	orderService      services.IOrderService
	handlersByStructType map[string]func(handler *OrderQueryHandler, message cqrs.RequestMessage) (interface{}, error)
}

// NewOrderQueryHandler a new Order query Handler
func NewOrderQueryHandler(orderService services.IOrderService) *OrderQueryHandler {

	// register handler functions
	handlers := map[string]func(handler *OrderQueryHandler, message cqrs.RequestMessage) (interface{}, error){}
	handlers[reflect.TypeOf(&LoadOrderQuery{}).String()] = LoadOrderQueryHandle
	handlers[reflect.TypeOf(&LoadAllOrderQuery{}).String()] = LoadAllOrderQueryHandle

	return &OrderQueryHandler{
		orderService:      orderService,
		handlersByStructType: handlers,
	}
}
