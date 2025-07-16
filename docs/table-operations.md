## Table of content

<!-- ts -->
  * [Table Operations](#table-operations)
    * [Creating Tables](#creating-tables)
    * [Getting Tables](#getting-tables)
    * [Deleting Tables](#deleting-tables)
    * [Updating Data](#updating-data)
<!-- te -->
## Table Operations

Database provides multiple ways to manipulate tables. Below are the main operations available.

### Creating Tables

Creates a new table in the database with specified name and column definitions. This operation allows you to define the
structure of your table by providing column names.

```go
// Create a new table with columns
table := db.NewTable("Users", []string{"name", "email", "age"})
```

### Getting Tables

Retrieves table information from the database. You can either get a list of all available tables or fetch a specific
table by its name.

```go
// Get all tables
tables := db.GetTables()
for _, table := range tables {
fmt.Println("Table:", table.GetName())
}

// Get specific table by name
userTable, err := db.GetTableByName("Users")
if err != nil {
fmt.Println("Table not found:", err)
return
}
```
### Deleting Tables

Removes an existing table from the database permanently. This operation will delete both the table structure and all its
data.

```go
// Delete a table
err := db.DeleteTable("Users")
if err != nil {
fmt.Println("Error deleting table:", err)
}
```
### Updating Data

Modifies existing table properties such as table name or column names. These operations allow you to restructure your
tables without losing the data.

```go
// Update table name
userTable.UpdateTableName("Customers")

// Update column name
err := userTable.UpdateColumnName("email", "email_address")
if err != nil {
    fmt.Println("Error updating column:", err)
}
```