package command

import (
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs"
	"reflect"
	"siigo.com/order/src/domain/order"
	"siigo.com/order/src/domain/services"
)

// CreateOrderCommand CreateOrder create a new inventory item
type CreateOrderCommand struct {
	Order *order.Order
}

type DeleteOrderCommand struct {
	Id string
}

type UpdateOrderCommand struct {
	Order *order.Order
}

type OrderCommandHandler struct {
	orderService      services.IOrderService
	handlersByStructType map[string]func(handler *OrderCommandHandler, message cqrs.RequestMessage) (interface{}, error)
}

// NewOrderCommandHandler a new Order Command Handler
func NewOrderCommandHandler(orderService services.IOrderService) *OrderCommandHandler {

	// register handler functions
	handlers := map[string]func(handler *OrderCommandHandler, message cqrs.RequestMessage) (interface{}, error){}
	handlers[reflect.TypeOf(&CreateOrderCommand{}).String()] = CreateOrderCommandHandle
	handlers[reflect.TypeOf(&DeleteOrderCommand{}).String()] = DeleteOrderCommandHandle
	handlers[reflect.TypeOf(&UpdateOrderCommand{}).String()] = UpdateOrderCommandHandle

	return &OrderCommandHandler{
		orderService:      orderService,
		handlersByStructType: handlers,
	}
}

type DeleteOrderHandler struct {
	orderService      services.IOrderService
	handlersByStructType map[string]func(handler *DeleteOrderHandler, message cqrs.RequestMessage) (interface{}, error)
}
