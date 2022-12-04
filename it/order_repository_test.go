//go:build integration
// +build integration

package it

import (
	"context"
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"siigo.com/order/it/containers"
	"siigo.com/order/src/domain/order"
	"siigo.com/order/src/infrastructure/finder"
	"siigo.com/order/src/infrastructure/repository"
	"sync"
	"testing"
)

func TestSaveOrderRepository(t *testing.T) {
	t.Parallel()

	// Arrange
	collection, container := containers.BuildMongoInfrastructure(t)
	orderRepository := repository.NewOrderRepository(context.Background(), collection, logrus.New())
	orderFinder := finder.NewOrderFinder(context.Background(), collection, logrus.New())
	defer container.Terminate(context.Background())

	// Act  --------
	totalOrders := 50
	wg := sync.WaitGroup{}
	wg.Add(totalOrders)

	for i := 0; i < totalOrders; i++ {
		go func() {
			defer wg.Done()
			orderDocument := order.NewOrder(cqrs.NewUUID())

			// save orders
			orderError := <-orderRepository.Save(orderDocument)
			if orderError != nil {
				t.Fatalf("Failed to save in mongo database: %v", orderError)
			}
		}()
	}

	wg.Wait()

	// Assert --------

	// Find all
	records := <-orderFinder.GetAll()

	assert.Nil(t, records.Error)
	assert.Equal(t, len(records.Orders), totalOrders)

}
