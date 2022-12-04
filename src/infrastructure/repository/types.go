package repository

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"siigo.com/order/src/domain/order"
)

// OrderRepository provides methods for processing mongo Person Connections
type OrderRepository struct {
	collection *mongo.Collection
	context    context.Context
	logger     *logrus.Logger
}

// NewOrderRepository a new NewOrderRepository
func NewOrderRepository(context context.Context, collection *mongo.Collection, logger *logrus.Logger) order.IOrderRepository {
	return OrderRepository{context: context, collection: collection, logger: logger}
}
