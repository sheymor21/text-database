package pkg

import "fmt"

type NotFoundError struct {
	itemName string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.itemName)
}
