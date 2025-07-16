package tdb

import (
	"fmt"
)

type NotFoundError struct {
	itemName string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.itemName)
}

type SqlSyntaxError struct {
	itemName string
}

func (e *SqlSyntaxError) Error() string {
	return fmt.Sprintf("Sql Syntax Error, not found %s", e.itemName)
}
