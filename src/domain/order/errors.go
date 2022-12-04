package order

import "fmt"

// Create domain custom errors

type ExceedsCostOrderError struct {
	Cost int32
}

func (e *ExceedsCostOrderError) Error() string {
	return fmt.Sprintf("cost %d exceeds the budget.", e.Cost)
}
