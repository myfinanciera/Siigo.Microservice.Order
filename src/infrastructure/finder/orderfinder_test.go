package finder

import (
	"context"
	"testing"

	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs"
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"siigo.com/order/it/containers"
	"siigo.com/order/src/domain/order"
	"siigo.com/order/src/infrastructure/repository"
)

func TestGetAll(t *testing.T) {
	//Arrange
	collection, container := containers.BuildMongoInfrastructure(t)
	defer container.Terminate(context.Background())

	finderOrder := NewOrderFinder(context.Background(), collection, logrus.New())
	repositoryOrder := repository.NewOrderRepository(context.Background(), collection, logrus.New())

	orderDocument := order.NewOrder(cqrs.NewUUID())
	orderError := <-repositoryOrder.Save(orderDocument)
	if orderError != nil {
		t.Fatalf("Failed to save in mongo database: %v", orderError)
	}

	// Act  --------
	response := <-finderOrder.GetAll()

	// Assert --------
	assert.Nil(t, response.Error)
	assert.Equal(t, len(response.Orders), 1)
}

func TestGetUnique(t *testing.T) {
	//Arrange
	collection, container := containers.BuildMongoInfrastructure(t)
	defer container.Terminate(context.Background())
	finderOrder := NewOrderFinder(context.Background(), collection, logrus.New())

	id, invalidIdError := uuid.FromString(cqrs.NewUUID())
	if invalidIdError != nil {
		t.Error(invalidIdError)
	}

	// Act  --------
	response := <-finderOrder.Get(id)

	// Assert --------
	assert.NotNil(t, response.Error) //Order NOt Found
	assert.NotNil(t, response.Order)
}
