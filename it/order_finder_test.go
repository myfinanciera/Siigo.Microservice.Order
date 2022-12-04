//go:build integration
// +build integration

package it

import (
	"context"
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs"
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"siigo.com/order/it/containers"
	"siigo.com/order/src/domain/order"
	"siigo.com/order/src/infrastructure/finder"
	"siigo.com/order/src/infrastructure/repository"
	"testing"
)

func TestGetAllOrdersFinder(t *testing.T) {

	t.Parallel()

	// Arrange  --------
	collection, container := containers.BuildMongoInfrastructure(t)
	defer container.Terminate(context.Background())

	finderOrder := finder.NewOrderFinder(context.Background(), collection, logrus.New())

	// Act  --------
	response := <-finderOrder.GetAll()

	// Assert --------
	assert.Nil(t, response.Error)
	assert.Equal(t, len(response.Orders), 0)
}

func TestGetAllOrdersWithElementsFinder(t *testing.T) {

	t.Parallel()

	// Arrange  ----------
	collection, container := containers.BuildMongoInfrastructure(t)
	defer container.Terminate(context.Background())

	finderOrder := finder.NewOrderFinder(context.Background(), collection, logrus.New())
	repository := repository.NewOrderRepository(
		context.Background(), collection, logrus.New(),
	)

	id := cqrs.NewUUID()
	orderDocument := order.NewOrder(id)
	orderError := <-repository.Save(orderDocument)
	if orderError != nil {
		t.Fatalf("Failed to save in mongo database: %v", orderError)
	}

	// Act  --------
	response := <-finderOrder.GetAll()

	// Assert --------
	assert.Nil(t, response.Error)
	assert.Equal(t, len(response.Orders), 1)
	assert.Equal(t, response.Orders[0].Id, id)
	assert.Equal(t, response.Orders[0].OccurredAt, orderDocument.OccurredAt)
}

func TestGetOrderNotFoundFinder(t *testing.T) {

	t.Parallel()

	// Arrange  --------
	collection, container := containers.BuildMongoInfrastructure(t)
	defer container.Terminate(context.Background())

	finderOrder := finder.NewOrderFinder(context.Background(), collection, logrus.New())

	id, invalidIdError := uuid.FromString(cqrs.NewUUID())
	if invalidIdError != nil {
		t.Error(invalidIdError)
	}

	// Act  --------
	response := <-finderOrder.Get(id)

	// Assert --------
	assert.NotNil(t, response.Error)
	assert.Equal(t, response.Error.Error(), "order not found")
}

func TestGetOrderFinder(t *testing.T) {

	t.Parallel()

	// Arrange
	collection, container := containers.BuildMongoInfrastructure(t)
	defer container.Terminate(context.Background())

	finderOrder := finder.NewOrderFinder(context.Background(), collection, logrus.New())
	repository := repository.NewOrderRepository(
		context.Background(), collection, logrus.New(),
	)

	id := cqrs.NewUUID()
	orderDocument := order.NewOrder(id)
	orderError := <-repository.Save(orderDocument)
	if orderError != nil {
		t.Fatalf("Failed to save in mongo database: %v", orderError)
	}

	// Act  --------
	uid, invalidIdError := uuid.FromString(id)
	if invalidIdError != nil {
		t.Error(invalidIdError)
	}
	response := <-finderOrder.Get(uid)

	// Assert --------
	assert.Nil(t, response.Error)
	assert.Equal(t, response.Order.Id, id)
	assert.Equal(t, response.Order.OccurredAt, orderDocument.OccurredAt)
}
