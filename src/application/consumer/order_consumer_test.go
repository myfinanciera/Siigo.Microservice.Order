package consumer

import (
	"testing"

	slim "dev.azure.com/SiigoDevOps/Siigo/_git/go-slim.git/abstractions"
)

func TestOrderConsumeEvent(t *testing.T) {

	//Arrange
	var message = new(slim.Message)

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("TestUserFail should have panicked!")
			}
		}()
		// This function should cause a panic
		TestConsumer(*message)
	}()

}
