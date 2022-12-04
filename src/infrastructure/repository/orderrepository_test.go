package repository

import (
	"context"
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"siigo.com/order/it/containers"
	"siigo.com/order/src/domain/order"
	"testing"
)

func TestSaveSuccess(t *testing.T) {
	//Arrange
	collection, container := containers.BuildMongoInfrastructure(t)
	defer container.Terminate(context.Background())
	repositoryOrder := NewOrderRepository(context.Background(), collection, logrus.New())
	orderDocument := order.NewOrder(cqrs.NewUUID())

	// Act  --------
	response := <-repositoryOrder.Save(orderDocument)

	// Assert --------
	assert.Nil(t, response)
}

func TestDeleteSuccess(t *testing.T) {
	//Arrange
	collection, container := containers.BuildMongoInfrastructure(t)
	defer container.Terminate(context.Background())
	repositoryOrder := NewOrderRepository(context.Background(), collection, logrus.New())
	orderDocument := order.NewOrder(cqrs.NewUUID())

	// Act  --------
	response := <-repositoryOrder.Delete(orderDocument)

	// Assert --------
	assert.Nil(t, response)
}

func TestUpdateSuccess(t *testing.T) {
	//Arrange
	collection, container := containers.BuildMongoInfrastructure(t)
	defer container.Terminate(context.Background())
	repositoryOrder := NewOrderRepository(context.Background(), collection, logrus.New())
	orderDocument := order.NewOrder(cqrs.NewUUID())

	// Act  --------
	response := <-repositoryOrder.Update(orderDocument)

	// Assert --------
	assert.Nil(t, response)
}
