package command

import (
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs"
	"errors"
	"reflect"
)

// Handle processes order commands.
func (handler *OrderCommandHandler) Handle(message cqrs.RequestMessage) (interface{}, error) {
	fn, ok := handler.handlersByStructType[reflect.TypeOf(message.Request()).String()]
	if !ok {
		return nil, errors.New("command type font found")
	}
	return fn(handler, message)
}

//llega el comando al handler
func CreateOrderCommandHandle(handler *OrderCommandHandler, message cqrs.RequestMessage) (interface{}, error) {   
	cmd := message.Request().(*CreateOrderCommand)
    //invoca el servicio de dominio
	err := handler.orderService.Save(cmd.Order, cqrs.Int(cmd.Order.OriginalVersion()))
	if err != nil {
		return nil, err
	}

	return cmd.Order, nil

}

func DeleteOrderCommandHandle(handler *OrderCommandHandler, message cqrs.RequestMessage) (interface{}, error) {

	cmd := message.Request().(*DeleteOrderCommand)

	err := handler.orderService.Delete(cmd.Id)
	if err != nil {
		return nil, err
	}

	return nil, nil

}

func UpdateOrderCommandHandle(handler *OrderCommandHandler, message cqrs.RequestMessage) (interface{}, error) {

	cmd := message.Request().(*UpdateOrderCommand)

	err := handler.orderService.Update(cmd.Order, cqrs.Int(cmd.Order.OriginalVersion()))
	if err != nil {
		return nil, err
	}

	return nil, nil

}
