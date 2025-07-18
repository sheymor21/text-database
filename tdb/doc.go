// Package tdb provides functionality for database operations and management.
//
// This package offers a set of utilities to interact with databases, handle errors,
// manage migrations, and perform various database operations including SQL queries
// and table manipulations.
//
// Key components of the package include:
//
// Error Handling:
//   - Provides specialized error types and handling mechanisms for database operations
//   - Includes error wrapping and categorization for better debugging
//
// Database Operations:
//   - Functions for common database operations
//   - Connection management and transaction support
//
// SQL Operations:
//   - SQL query building and execution
//   - Parameter binding and result processing
//
// Table Operations:
//   - Table creation, modification and deletion
//   - Schema management functions
//
// Encoding:
//   - Data serialization and deserialization utilities
//   - Format conversion for database storage
//
// Utilities:
//   - Helper functions for common database tasks
//   - Convenience wrappers for repetitive operations
//
// Example usage:
//
//		import "github.com/sheymor21/text-database/tdb"
//
//		func main() {
//			// Create a DbConfiguration
//			config := tdb.DbConfig{EncryptionKey: "", DatabaseName: "database.txt"}
//
//
//			// Create a new database
//			db, err := config.CreateDatabase()
//			if err != nil {
//				fmt.Println(err)
//				return
//				}
//
//			// Print the tables
//			db.PrintTables()
//	}
package tdb
