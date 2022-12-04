package repository

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestNewOrderRespository(t *testing.T) {
	// Arrange
	mongo := new(*mongo.Collection)
	//Act
	finderOrder := NewOrderRepository(context.Background(), *mongo, logrus.New())
	//Assert
	assert.NotNil(t, finderOrder)
}
