package controller

import (
	"context"
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs"
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs/uuid"
	"encoding/json"
	"google.golang.org/protobuf/types/known/emptypb"
	"siigo.com/order/src/api/logger"
	"siigo.com/order/src/api/mapper"
	"siigo.com/order/src/api/proto/order/v1"
	"siigo.com/order/src/application/command"
	"siigo.com/order/src/application/query"
	"siigo.com/order/src/domain/order"
)

func (controller *Controller) FindOrders(ctx context.Context, empty *emptypb.Empty) (*orderv1.FindOrdersResponse, error) {
	controller.Logger.Info("Order Find all.")

	// create query message
	em := cqrs.NewQueryMessage(&query.LoadAllOrderQuery{})	
	orders, err := controller.Bus.Send(em)
	if err != nil {
		return nil, err
	}

	var cts []*orderv1.Order
	bytes, _ := json.Marshal(orders)
	_ = json.Unmarshal(bytes, &cts)

	return &orderv1.FindOrdersResponse{Orders: cts}, nil
}

func (controller *Controller) AddOrder(ctx context.Context, grpcOrder *orderv1.AddOrderRequest) (*orderv1.AddOrderResponse, error) {

	controller.Logger.Infof("Create Order %+v", grpcOrder.GetOrder())

	// crear el agregado con el uuid
	orderAggregate := order.NewOrderWithUuid()
	//paso el agregado a modelo grpc grpcOrder
	grpcOrder.Order.Id = orderAggregate.Id
	grpcOrder.Order.Occurred = orderAggregate.Occurred

	// map grpc to aggregate
	err := mapper.Adap(orderAggregate, grpcOrder.Order)
	if err != nil {
		return nil, err
	}

	// create command message
	em := cqrs.NewCommandMessage(orderAggregate.Id, &command.CreateOrderCommand{
		Order: orderAggregate,
	})
     
	//go channels para commandBus
	// send command to handler(s)
	_, err = controller.Bus.Send(em) //Bus - > Handler

	if err != nil {
		return nil, err
	}

	// track business metric
	controller.
		Logger.
		WithFields(
			logger.
				NewBusinessLogger().
				WithApiLayer().
				WithMetricType("OrderCreated").
				ToFields(),
		).Info(grpcOrder.Order.GetId())

	return &orderv1.AddOrderResponse{Order: grpcOrder.Order}, nil
}

func (controller *Controller) GetOrder(ctx context.Context, queryOrder *orderv1.GetOrderRequest) (*orderv1.GetOrderResponse, error) {

	controller.Logger.Debugf("Get Order %+v", queryOrder)

	id, invalidIdError := uuid.FromString(queryOrder.Id)
	if invalidIdError != nil {
		controller.Logger.Error(invalidIdError)
		return nil, invalidIdError
	}

	// create query message
	em := cqrs.NewQueryMessage(&query.LoadOrderQuery{Id: id})

	ct, err := controller.Bus.Send(em)
	if err != nil {
		return nil, err
	}

	// map aggregate to grpc
	response := &orderv1.Order{}
	err = mapper.Adap(response, ct)
	if err != nil {
		return nil, err
	}

	return &orderv1.GetOrderResponse{Order: response}, err

}

func (controller *Controller) DeleteOrder(ctx context.Context, request *orderv1.DeleteOrderRequest) (*orderv1.DeleteOrderResponse, error) {

	controller.Logger.Debugf("Delete order %+v", request)

	//Delete ELement
	em := cqrs.NewCommandMessage(request.Id, &command.DeleteOrderCommand{
		Id: request.Id,
	})

	// send command to handler(s)
	controller.Bus.Send(em)
	return &orderv1.DeleteOrderResponse{Order: nil}, nil
}

func (controller *Controller) UpdateOrder(ctx context.Context, request *orderv1.UpdateOrderRequest) (*orderv1.UpdateOrderResponse, error) {

	controller.Logger.Debugf("Update order %+v", request)

	// map grpc to aggregate
	orderAggregate := order.NewOrder(request.Order.Id)
	err := mapper.Adap(orderAggregate, request.Order)
	if err != nil {
		return nil, err
	}

	// create command message
	em := cqrs.NewCommandMessage(request.Order.Id, &command.UpdateOrderCommand{
		Order: orderAggregate,
	})

	// send command to handler(s)
	_, err = controller.Bus.Send(em)

	if err != nil {
		return nil, err
	}

	controller.
		Logger.
		WithFields(
			logger.
				NewBusinessLogger().
				WithApiLayer().
				WithMetricType("OrderUpdated").
				ToFields(),
		).
		Info(orderAggregate.Id)

	return &orderv1.UpdateOrderResponse{Order: request.Order}, nil

}
