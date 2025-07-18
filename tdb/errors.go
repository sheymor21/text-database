package tdb

import (
	"fmt"
)

// NotFoundError represents an error when a requested item cannot be found in the database.
type NotFoundError struct {
	itemName string
}

// Error returns a formatted error message indicating which item was not found.
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.itemName)
}

// SqlSyntaxError represents an error in SQL syntax with the specified item.
type SqlSyntaxError struct {
	itemName string
}

// Error returns a formatted error message indicating SQL syntax error with the specified item.
func (e *SqlSyntaxError) Error() string {
	return fmt.Sprintf("Sql Syntax Error, not found %s", e.itemName)
}
