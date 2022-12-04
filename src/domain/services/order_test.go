package services

import (
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs"
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs/uuid"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	finder "siigo.com/order/mocks/src/domain/order"
	repository "siigo.com/order/mocks/src/domain/order"
	"siigo.com/order/src/domain/order"
	"testing"
	"time"
)

func TestNewOrderService(t *testing.T) {
	// Arrange
	dispatcher := &cqrs.DomainEventsDispatcher{}
	repository := &repository.IOrderRepository{}
	finder := &finder.IOrderFinder{}

	id := cqrs.NewUUID()
	orderMock := order.NewOrder(id)

	repository.
		On("Save", mock.Anything).
		Return(nil, nil)

	finder.
		On("Get", mock.Anything).
		Return(orderMock, nil)

	//Act
	var inter = NewOrderService(dispatcher, nil, repository, finder)

	//Assert
	assert.NotNil(t, inter)
	assert.NotNil(t, repository)
	assert.NotNil(t, finder)

}

func TestNewOrderDispatcherNil(t *testing.T) {
	// Arrange

	repository := &repository.IOrderRepository{}
	finder := &finder.IOrderFinder{}

	//Act
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("TestUserFail should have panicked!")
			}
		}()
		// This function should cause a panic
		NewOrderService(nil, nil, repository, finder)
	}()
}

func TestNewOrderServiceGet(t *testing.T) {
	// Arrange
	dispatcher := &cqrs.DomainEventsDispatcher{}
	finderMock := &finder.IOrderFinder{}
	responsemock := make(chan *order.OrderResponse)
	defer close(responsemock)

	id := cqrs.NewUUID()
	uid, _ := uuid.FromString(id)
	//orderMock := order.NewOrder(id)

	finderMock.On("Get", mock.Anything).
		Return(responsemock)

	//Act
	var service = NewOrderService(dispatcher, nil, nil, finderMock)
	contracresponse := &order.OrderResponse{Order: order.NewOrder(id), Error: nil}

	go func() {
		time.Sleep(time.Second)
		responsemock <- contracresponse
	}()

	var response, error = service.Get(uid)

	//Assert
	assert.NotNil(t, response)
	assert.Nil(t, error)
	finderMock.AssertNumberOfCalls(t, "Get", 1)
	finderMock.AssertCalled(t, "Get", uid)
	finderMock.AssertExpectations(t)
}

func TestNewOrderServiceLoadAll(t *testing.T) {
	// Arrange
	dispatcher := &cqrs.DomainEventsDispatcher{}
	finderMock := &finder.IOrderFinder{}
	responsemock := make(chan order.OrdersResponse)
	defer close(responsemock)

	//orderMock := order.NewOrder(id)

	finderMock.On("GetAll", mock.Anything).
		Return(responsemock)

	//Act
	var service = NewOrderService(dispatcher, nil, nil, finderMock)
	contracresponse := order.OrdersResponse{}

	go func() {
		time.Sleep(time.Second)
		responsemock <- contracresponse
	}()

	var response, error = service.LoadAll()

	//Assert
	assert.Nil(t, response)
	assert.Nil(t, error)
}

func TestNewOrderServiceLoadAllError(t *testing.T) {
	// Arrange
	dispatcher := &cqrs.DomainEventsDispatcher{}
	finderMock := &finder.IOrderFinder{}
	responsemock := make(chan order.OrdersResponse)
	defer close(responsemock)

	//orderMock := order.NewOrder(id)

	finderMock.On("GetAll", mock.Anything).
		Return(responsemock)

	//Act
	var service = NewOrderService(dispatcher, nil, nil, finderMock)
	contracresponse := order.OrdersResponse{Error: errors.New("order not found")}

	go func() {
		time.Sleep(time.Second)
		responsemock <- contracresponse
	}()

	var response, error = service.LoadAll()

	//Assert
	assert.Nil(t, response)
	assert.NotNil(t, error)
}

func TestNewOrderServiceLoadErrorResponse(t *testing.T) {
	// Arrange
	dispatcher := &cqrs.DomainEventsDispatcher{}
	finderMock := &finder.IOrderFinder{}
	responsemock := make(chan *order.OrderResponse)
	defer close(responsemock)

	id := cqrs.NewUUID()
	uid, _ := uuid.FromString(id)
	//orderMock := order.NewOrder(id)

	finderMock.On("Get", mock.Anything).
		Return(responsemock)

	//Act
	var service = NewOrderService(dispatcher, nil, nil, finderMock)

	go func() {
		time.Sleep(time.Second)
		responsemock <- &order.OrderResponse{Error: errors.New("")}
	}()

	var response, error = service.Get(uid)

	//Assert
	assert.Nil(t, response)
	assert.NotNil(t, error)
	finderMock.AssertNumberOfCalls(t, "Get", 1)
	finderMock.AssertCalled(t, "Get", uid)
	finderMock.AssertExpectations(t)
}

func TestNewOrderServiceSave(t *testing.T) {
	// Arrange
	dispatcher := &cqrs.DomainEventsDispatcher{}
	repositoryMock := &repository.IOrderRepository{}
	responsemock := make(chan error)
	defer close(responsemock)
	ct := &order.Order{}

	repositoryMock.On("Save", mock.Anything).
		Return(responsemock)

	//Act
	var service = NewOrderService(dispatcher, nil, repositoryMock, nil)

	go func() {
		time.Sleep(time.Second)
		responsemock <- nil
	}()

	var result = service.Save(ct, nil)

	//Assert
	assert.Nil(t, result)
}

func TestNewOrderServiceSaveErrorGeneral(t *testing.T) {
	// Arrange
	dispatcher := &cqrs.DomainEventsDispatcher{}
	repositoryMock := &repository.IOrderRepository{}
	responsemock := make(chan error)
	defer close(responsemock)
	ct := &order.Order{}

	repositoryMock.On("Save", mock.Anything).
		Return(responsemock)

	//Act
	var service = NewOrderService(dispatcher, nil, repositoryMock, nil)

	go func() {
		time.Sleep(time.Second)
		responsemock <- errors.New("Fail General")
	}()

	var result = service.Save(ct, nil)

	//Assert
	assert.NotNil(t, result)
}

func TestNewOrderServiceErrorIsWithinBudget(t *testing.T) {
	// Arrange
	dispatcher := &cqrs.DomainEventsDispatcher{}
	repositoryMock := &repository.IOrderRepository{}

	ct := &order.Order{
		//Cost: 20,
	}

	//Act
	var service = NewOrderService(dispatcher, nil, repositoryMock, nil)

	var result = service.Save(ct, nil)

	//Assert
	assert.NotNil(t, result)
	assert.Equal(t, result.Error(), "rpc error: code = InvalidArgument desc = cost 20 exceeds the budget.")
}

func TestNewOrderServiceDelete(t *testing.T) {
	// Arrange
	dispatcher := &cqrs.DomainEventsDispatcher{}
	finderMock := &finder.IOrderFinder{}
	repositoryMock := &repository.IOrderRepository{}
	responsemock := make(chan error)
	defer close(responsemock)

	responsemockLoad := make(chan *order.OrderResponse)
	defer close(responsemockLoad)

	id := cqrs.NewUUID()

	finderMock.On("Get", mock.Anything).
		Return(responsemockLoad)
	contracresponseLoad := &order.OrderResponse{Order: order.NewOrder(id), Error: nil}

	repositoryMock.On("Delete", mock.Anything).
		Return(responsemock)

	//Act
	var service = NewOrderService(dispatcher, nil, repositoryMock, finderMock)

	go func() {
		time.Sleep(time.Second)
		responsemockLoad <- contracresponseLoad
		responsemock <- nil
	}()

	var error = service.Delete(id)

	//Assert
	assert.Nil(t, error)
}

func TestNewOrderServiceDeleteErrorId(t *testing.T) {
	// Arrange
	dispatcher := &cqrs.DomainEventsDispatcher{}
	finderMock := &finder.IOrderFinder{}
	repositoryMock := &repository.IOrderRepository{}

	//Act
	var service = NewOrderService(dispatcher, nil, repositoryMock, finderMock)

	var error = service.Delete("id")

	//Assert
	assert.NotNil(t, error)
}

func TestNewOrderServiceDeleteContracNotFound(t *testing.T) {
	// Arrange
	dispatcher := &cqrs.DomainEventsDispatcher{}
	finderMock := &finder.IOrderFinder{}
	repositoryMock := &repository.IOrderRepository{}
	responsemock := make(chan error)
	defer close(responsemock)

	responsemockLoad := make(chan *order.OrderResponse)
	defer close(responsemockLoad)

	id := cqrs.NewUUID()

	finderMock.On("Get", mock.Anything).
		Return(responsemockLoad)
	contracresponseLoad := &order.OrderResponse{Order: order.NewOrder(id), Error: errors.New("order not found")}

	repositoryMock.On("Delete", mock.Anything).
		Return(responsemock)

	//Act
	var service = NewOrderService(dispatcher, nil, repositoryMock, finderMock)

	go func() {
		time.Sleep(time.Second)
		responsemockLoad <- contracresponseLoad
	}()

	var error = service.Delete(id)

	//Assert
	assert.NotNil(t, error)
}

func TestNewOrderServiceDeleteErrorGeneral(t *testing.T) {
	// Arrange
	dispatcher := &cqrs.DomainEventsDispatcher{}
	finderMock := &finder.IOrderFinder{}
	repositoryMock := &repository.IOrderRepository{}
	responsemock := make(chan error)
	defer close(responsemock)

	responsemockLoad := make(chan *order.OrderResponse)
	defer close(responsemockLoad)

	id := cqrs.NewUUID()

	finderMock.On("Get", mock.Anything).
		Return(responsemockLoad)
	contracresponseLoad := &order.OrderResponse{Order: order.NewOrder(id), Error: nil}

	repositoryMock.On("Delete", mock.Anything).
		Return(responsemock)

	//Act
	var service = NewOrderService(dispatcher, nil, repositoryMock, finderMock)

	go func() {
		time.Sleep(time.Second)
		responsemockLoad <- contracresponseLoad
		responsemock <- errors.New("Fail General")
	}()

	var error = service.Delete(id)

	//Assert
	assert.NotNil(t, error)
}

func TestNewOrderServiceUpdate(t *testing.T) {
	// Arrange
	dispatcher := &cqrs.DomainEventsDispatcher{}
	finderMock := &finder.IOrderFinder{}
	repositoryMock := &repository.IOrderRepository{}
	responsemock := make(chan error)
	defer close(responsemock)
	ct := &order.Order{}
	id := cqrs.NewUUID()
	ct.Id = id

	responsemockLoad := make(chan *order.OrderResponse)
	defer close(responsemockLoad)

	repositoryMock.On("Update", mock.Anything).
		Return(responsemock)

	finderMock.On("Get", mock.Anything).
		Return(responsemockLoad)
	contracresponseLoad := &order.OrderResponse{Order: order.NewOrder(id), Error: nil}

	//Act
	var service = NewOrderService(dispatcher, nil, repositoryMock, finderMock)

	go func() {
		time.Sleep(time.Second)
		responsemockLoad <- contracresponseLoad
		responsemock <- nil
	}()

	var error = service.Update(ct, nil)

	//Assert
	assert.Nil(t, error)
}

func TestNewOrderServiceUpdateIsWithinBudget(t *testing.T) {
	// Arrange
	dispatcher := &cqrs.DomainEventsDispatcher{}
	finderMock := &finder.IOrderFinder{}
	repositoryMock := &repository.IOrderRepository{}
	responsemock := make(chan error)
	defer close(responsemock)
	ct := &order.Order{}
	id := cqrs.NewUUID()
	ct.Id = id
	//ct.Cost = 20

	//Act
	var service = NewOrderService(dispatcher, nil, repositoryMock, finderMock)

	var error = service.Update(ct, nil)

	//Assert
	assert.NotNil(t, error)
}

func TestNewOrderServiceUpdateOrderNotFound(t *testing.T) {
	// Arrange
	dispatcher := &cqrs.DomainEventsDispatcher{}
	finderMock := &finder.IOrderFinder{}
	repositoryMock := &repository.IOrderRepository{}
	responsemock := make(chan error)
	defer close(responsemock)
	ct := &order.Order{}
	id := cqrs.NewUUID()
	ct.Id = id

	responsemockLoad := make(chan *order.OrderResponse)
	defer close(responsemockLoad)

	repositoryMock.On("Update", mock.Anything).
		Return(responsemock)

	finderMock.On("Get", mock.Anything).
		Return(responsemockLoad)
	contracresponseLoad := &order.OrderResponse{Order: order.NewOrder(id), Error: errors.New("order not found")}

	//Act
	var service = NewOrderService(dispatcher, nil, repositoryMock, finderMock)

	go func() {
		time.Sleep(time.Second)
		responsemockLoad <- contracresponseLoad
	}()

	var error = service.Update(ct, nil)

	//Assert
	assert.NotNil(t, error)
}

func TestNewOrderServiceUpdateOrderErrorTransaction(t *testing.T) {
	// Arrange
	dispatcher := &cqrs.DomainEventsDispatcher{}
	finderMock := &finder.IOrderFinder{}
	repositoryMock := &repository.IOrderRepository{}
	responsemock := make(chan error)
	defer close(responsemock)
	ct := &order.Order{}
	id := cqrs.NewUUID()
	ct.Id = id

	responsemockLoad := make(chan *order.OrderResponse)
	defer close(responsemockLoad)

	repositoryMock.On("Update", mock.Anything).
		Return(responsemock)

	finderMock.On("Get", mock.Anything).
		Return(responsemockLoad)
	contracresponseLoad := &order.OrderResponse{Order: order.NewOrder(id), Error: nil}

	//Act
	var service = NewOrderService(dispatcher, nil, repositoryMock, finderMock)

	go func() {
		time.Sleep(time.Second)
		responsemockLoad <- contracresponseLoad
		responsemock <- errors.New("Fail General")
	}()

	var error = service.Update(ct, nil)

	//Assert
	assert.NotNil(t, error)
}

func TestNewOrderServiceUpdateErrorId(t *testing.T) {
	// Arrange
	dispatcher := &cqrs.DomainEventsDispatcher{}
	finderMock := &finder.IOrderFinder{}
	repositoryMock := &repository.IOrderRepository{}
	responsemock := make(chan error)
	defer close(responsemock)
	ct := &order.Order{}
	id := cqrs.NewUUID()

	responsemockLoad := make(chan *order.OrderResponse)
	defer close(responsemockLoad)

	repositoryMock.On("Update", mock.Anything).
		Return(responsemock)

	finderMock.On("Get", mock.Anything).
		Return(responsemockLoad)
	contracresponseLoad := &order.OrderResponse{Order: order.NewOrder(id), Error: nil}

	//Act
	var service = NewOrderService(dispatcher, nil, repositoryMock, finderMock)

	go func() {
		time.Sleep(time.Second)
		responsemockLoad <- contracresponseLoad
		responsemock <- nil
	}()

	var error = service.Update(ct, nil)

	//Assert
	assert.NotNil(t, error)
}
