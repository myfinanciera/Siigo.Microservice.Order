package command_test

import (
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mockServices "siigo.com/order/mocks/src/domain/services"
	"siigo.com/order/src/application/command"
	"siigo.com/order/src/domain/order"
	"testing"
)

func TestCreateOrderCommandHandlerWithBudgedError_Handle(t *testing.T) {

	// Arrange
	domainOrderService := &mockServices.IOrderService{}
	domainError := &order.ExceedsCostOrderError{}

	domainOrderService.
		On("Save", mock.Anything, mock.Anything).
		Return(domainError)

	id := cqrs.NewUUID()
	ct := order.NewOrder(id)
	commandMessage := cqrs.NewCommandMessage(id, &command.CreateOrderCommand{
		Order: ct,
	})

	// Act
	handler := command.NewOrderCommandHandler(domainOrderService)
	resp, err := handler.Handle(commandMessage)

	// Assert
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, domainError))
	assert.EqualError(t, err, domainError.Error())
	assert.Nil(t, resp)

	domainOrderService.AssertNumberOfCalls(t, "Save", 1)
	domainOrderService.AssertCalled(t, "Save", ct, cqrs.Int(ct.OriginalVersion()))
	domainOrderService.AssertExpectations(t)

}

func TestCreateOrderCommandHandlerSuccess_Handle(t *testing.T) {

	// Arrange
	domainOrderService := &mockServices.IOrderService{}
	id := cqrs.NewUUID()
	ct := order.NewOrder(id)

	domainOrderService.
		On("Save", mock.Anything, mock.Anything).
		Return(nil)

	commandMessage := cqrs.NewCommandMessage(id, &command.CreateOrderCommand{
		Order: ct,
	})

	// Act
	handler := command.NewOrderCommandHandler(domainOrderService)
	response, err := handler.Handle(commandMessage)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, id, response.(*order.Order).Id)

	domainOrderService.AssertNumberOfCalls(t, "Save", 1)
	domainOrderService.AssertCalled(t, "Save", ct, cqrs.Int(ct.OriginalVersion()))
	domainOrderService.AssertExpectations(t)

}

func TestDeleteOrderCommandHandlerSuccess_Handle(t *testing.T) {

	// Arrange
	domainOrderService := &mockServices.IOrderService{}
	id := cqrs.NewUUID()

	domainOrderService.
		On("Delete", mock.Anything, mock.Anything).
		Return(nil)

	commandMessage := cqrs.NewCommandMessage(id, &command.DeleteOrderCommand{
		Id: id,
	})

	// Act
	handler := command.NewOrderCommandHandler(domainOrderService)
	_, err := handler.Handle(commandMessage)

	// Assert
	assert.Nil(t, err)

	domainOrderService.AssertNumberOfCalls(t, "Delete", 1)
	domainOrderService.AssertExpectations(t)

}

func TestDeleteOrderCommandHandlerSuccess_ErrorHandle(t *testing.T) {

	// Arrange
	domainOrderService := &mockServices.IOrderService{}
	id := cqrs.NewUUID()
	domainError := errors.New("order problems")

	domainOrderService.
		On("Delete", mock.Anything, mock.Anything).
		Return(domainError)

	commandMessage := cqrs.NewCommandMessage(id, &command.DeleteOrderCommand{
		Id: id,
	})

	// Act
	handler := command.NewOrderCommandHandler(domainOrderService)
	_, err := handler.Handle(commandMessage)

	// Assert
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, domainError))
	assert.EqualError(t, err, domainError.Error())

	domainOrderService.AssertNumberOfCalls(t, "Delete", 1)
	domainOrderService.AssertExpectations(t)

}

func TestUpdateOrderCommandHandlerSuccess_Handle(t *testing.T) {

	// Arrange
	domainOrderService := &mockServices.IOrderService{}
	id := cqrs.NewUUID()
	ct := order.NewOrder(id)

	domainOrderService.
		On("Update", mock.Anything, mock.Anything).
		Return(nil)

	commandMessage := cqrs.NewCommandMessage(id, &command.UpdateOrderCommand{
		Order: ct,
	})

	// Act
	handler := command.NewOrderCommandHandler(domainOrderService)
	response, err := handler.Handle(commandMessage)

	// Assert
	assert.Nil(t, err)
	assert.Nil(t, response)

	domainOrderService.AssertNumberOfCalls(t, "Update", 1)
	domainOrderService.AssertCalled(t, "Update", ct, cqrs.Int(ct.OriginalVersion()))
	domainOrderService.AssertExpectations(t)

}

func TestUpdateOrderCommandHandlerSuccess_ErrorHandle(t *testing.T) {

	// Arrange
	domainOrderService := &mockServices.IOrderService{}
	id := cqrs.NewUUID()
	ct := order.NewOrder(id)
	domainError := errors.New("order problems")

	domainOrderService.
		On("Update", mock.Anything, mock.Anything).
		Return(domainError)

	commandMessage := cqrs.NewCommandMessage(id, &command.UpdateOrderCommand{
		Order: ct,
	})

	// Act
	handler := command.NewOrderCommandHandler(domainOrderService)
	_, err := handler.Handle(commandMessage)

	// Assert
	assert.NotNil(t, err)

	domainOrderService.AssertNumberOfCalls(t, "Update", 1)
	domainOrderService.AssertCalled(t, "Update", ct, cqrs.Int(ct.OriginalVersion()))
	domainOrderService.AssertExpectations(t)

}
