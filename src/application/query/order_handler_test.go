package query_test

import (
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs"
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs/uuid"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mockServices "siigo.com/order/mocks/src/domain/services"
	"siigo.com/order/src/application/query"
	"siigo.com/order/src/domain/order"
	"testing"
)

func TestLoadOrderQueryHandlerDomainServiceError_Handle(t *testing.T) {

	// Arrange
	domainOrderService := &mockServices.IOrderService{}
	id, _ := uuid.FromString(cqrs.NewUUID())
	errorMessage := "domain error"

	domainOrderService.
		On("Get", mock.Anything).
		Return(nil, errors.New(errorMessage))

	commandMessage := cqrs.NewQueryMessage(&query.LoadOrderQuery{Id: id})

	// Act
	handler := query.NewOrderQueryHandler(domainOrderService)
	response, err := handler.Handle(commandMessage)

	// Assert
	assert.NotNil(t, err)
	assert.Nil(t, response)

	domainOrderService.AssertNumberOfCalls(t, "Get", 1)
	domainOrderService.AssertCalled(t, "Get", id)
	domainOrderService.AssertExpectations(t)
}

func TestLoadOrderQueryHandler_Handle(t *testing.T) {

	// Arrange
	domainOrderService := &mockServices.IOrderService{}
	id := cqrs.NewUUID()
	uid, _ := uuid.FromString(id)
	orderMock := order.NewOrder(id)

	domainOrderService.
		On("Get", mock.Anything).
		Return(orderMock, nil)

	commandMessage := cqrs.NewQueryMessage(&query.LoadOrderQuery{Id: uid})

	// Act
	handler := query.NewOrderQueryHandler(domainOrderService)
	response, err := handler.Handle(commandMessage)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, response)

	domainOrderService.AssertNumberOfCalls(t, "Get", 1)
	domainOrderService.AssertCalled(t, "Get", uid)
	domainOrderService.AssertExpectations(t)
}

func TestLoadAllOrderQueryHandler_Handle(t *testing.T) {

	// Arrange
	domainOrderService := &mockServices.IOrderService{}

	domainOrderService.
		On("LoadAll", mock.Anything).
		Return([]*order.Order{}, nil)

	commandMessage := cqrs.NewQueryMessage(&query.LoadAllOrderQuery{})

	// Act
	handler := query.NewOrderQueryHandler(domainOrderService)
	response, err := handler.Handle(commandMessage)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, response)

	domainOrderService.AssertNumberOfCalls(t, "LoadAll", 1)
	domainOrderService.AssertExpectations(t)
}

func TestLoadAllOrderQueryHandler_ErrorHandle(t *testing.T) {

	// Arrange
	domainOrderService := &mockServices.IOrderService{}
	errorMessage := "domain error"
	domainOrderService.
		On("LoadAll", mock.Anything).
		Return(nil, errors.New(errorMessage))

	commandMessage := cqrs.NewQueryMessage(&query.LoadAllOrderQuery{})

	// Act
	handler := query.NewOrderQueryHandler(domainOrderService)
	response, err := handler.Handle(commandMessage)

	// Assert
	assert.NotNil(t, err)
	assert.Nil(t, response)

	domainOrderService.AssertNumberOfCalls(t, "LoadAll", 1)
	domainOrderService.AssertExpectations(t)
}
