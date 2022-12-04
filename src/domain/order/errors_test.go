package order

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorMessage(t *testing.T) {

	//Arrange
	var exceeds = new(ExceedsCostOrderError)
	exceeds.Cost = 1

	//Act
	var msg = exceeds.Error()

	//Assert
	assert.NotNil(t, msg)
}
