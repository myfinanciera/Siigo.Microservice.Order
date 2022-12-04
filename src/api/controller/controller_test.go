package controller

import (
	"context"
	"testing"

	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs"
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
	orderv1 "siigo.com/order/src/api/proto/order/v1"
)

func TestNewController(t *testing.T) {
	// Arrange
	cqrsDispatcher := *new(cqrs.Dispatcher)

	//Act
	controller := NewController(cqrsDispatcher, new(logrus.Logger))

	//Assert
	assert.NotNil(t, controller)
}

func TestNewControllerAddOrderDontHandler(t *testing.T) {
	// Arrange
	dispatcher := cqrs.NewInMemoryDispatcher()
	ctx := context.Background()
	// map aggregate to grpc
	order := &orderv1.Order{}
	orderRequest := &orderv1.AddOrderRequest{}
	orderRequest.Order = order

	// Act
	ctl := NewController(dispatcher, logrus.New())
	var response, error = ctl.AddOrder(ctx, orderRequest)

	//Assert
	assert.Nil(t, response)
	assert.NotNil(t, error)
}

func TestNewControllerLoadOrderDontHandler(t *testing.T) {
	// Arrange
	dispatcher := cqrs.NewInMemoryDispatcher()
	ctx := context.Background()

	orderRequest := &orderv1.GetOrderRequest{}
	orderRequest.Id = uuid.NewUUID()

	// Act
	ctl := NewController(dispatcher, logrus.New())
	var response, error = ctl.GetOrder(ctx, orderRequest)

	//Assert
	assert.Nil(t, response)
	assert.NotNil(t, error)
}

func TestNewControllerLoadIdError(t *testing.T) {
	// Arrange
	dispatcher := cqrs.NewInMemoryDispatcher()
	ctx := context.Background()

	orderRequest := &orderv1.GetOrderRequest{}

	// Act
	ctl := NewController(dispatcher, logrus.New())
	var response, error = ctl.GetOrder(ctx, orderRequest)

	//Assert
	assert.Nil(t, response)
	assert.NotNil(t, error)
}

func TestNewControllerUpdateOrder(t *testing.T) {
	// Arrange
	dispatcher := cqrs.NewInMemoryDispatcher()
	ctx := context.Background()
	// map aggregate to grpc
	order := &orderv1.Order{}
	orderRequest := &orderv1.UpdateOrderRequest{}
	orderRequest.Order = order

	// Act
	ctl := NewController(dispatcher, logrus.New())
	var response, error = ctl.UpdateOrder(ctx, orderRequest)

	//Assert
	assert.Nil(t, response)
	assert.NotNil(t, error)
}

func TestNewControllerDeleteOrder(t *testing.T) {
	// Arrange
	dispatcher := cqrs.NewInMemoryDispatcher()
	ctx := context.Background()

	orderRequest := &orderv1.DeleteOrderRequest{}
	orderRequest.Id = uuid.NewUUID()

	// Act
	ctl := NewController(dispatcher, logrus.New())
	var response, error = ctl.DeleteOrder(ctx, orderRequest)

	//Assert
	assert.NotNil(t, response)
	assert.Nil(t, error)
}

func TestNewControllerLoadAllOrder(t *testing.T) {
	// Arrange
	dispatcher := cqrs.NewInMemoryDispatcher()
	ctx := context.Background()
	var protoReq *emptypb.Empty

	orderRequest := &orderv1.GetOrderRequest{}
	orderRequest.Id = uuid.NewUUID()

	// Act
	ctl := NewController(dispatcher, logrus.New())
	var response, error = ctl.FindOrders(ctx, protoReq)

	//Assert
	assert.Nil(t, response)
	assert.NotNil(t, error)
}
