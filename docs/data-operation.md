## Table of content

<!-- ts -->
  * [Data Operations](#data-operations)
    * [Adding Data to Tables](#adding-data-to-tables)
    * [Querying Data](#querying-data)
    * [Updating Data](#updating-data)
    * [Deleting Data](#deleting-data)
  * [Row Operations](#row-operations)
    * [Working with Rows](#working-with-rows)
    * [Sorting Result](#sorting-result)
<!-- te -->
## Data Operations

Database provides multiple ways to manipulate data in tables. Below are the main operations available.

### Adding Data to Tables

This section demonstrates how to insert single or multiple values into database tables. You can add individual values
using AddValue() or multiple values at once with AddValues().

```go
// Add single values
err := userTable.AddValue("name", "John Doe")
if err != nil {
    fmt.Println("Error adding value:", err)
}

// Add multiple values at once
userTable.AddValues("John", "john@example.com", "30")
userTable.AddValues("Jane", "jane@example.com", "25")
```
### Querying Data

The database supports various query operations to retrieve data. You can fetch all rows, get specific rows by ID, or
search for records matching certain criteria.

```go
// Get all rows
rows := userTable.GetRows()
for _, row := range rows {
fmt.Println(row.String())
}

// Get row by ID
user, err := userTable.GetRowById("1")
if err != nil {
fmt.Println("User not found:", err)
} else {
fmt.Println("User:", user.String())
}

// Search for specific values
user, err := userTable.SearchOne("name", "John")
if err != nil {
fmt.Println("User not found:", err)
} else {
fmt.Println("Found user:", user.String())
}

// Search all matching records
users := userTable.SearchAll("age", "30")
for _, user := range users {
fmt.Println("User:", user.String())
}
```

### Updating Data

Update operations allow you to modify existing data in the database. You can update specific values in rows using the
row ID and column name.

```go
// Update a specific value
err := userTable.UpdateValue("name", "1", "John Smith")
if err != nil {
    fmt.Println("Error updating value:", err)
}
```
### Deleting Data

Delete operations support removing rows and columns from tables. You can delete individual rows with or without cascade
deletion for related records or remove entire columns from tables.

```go
// Delete a row by ID
err := userTable.DeleteRow("1", false) // false = no cascade delete
if err != nil {
    fmt.Println("Error deleting row:", err)
}

// Delete with cascade (removes related foreign key records)
err := userTable.DeleteRow("1", true)
if err != nil {
    fmt.Println("Error deleting row:", err)
}

// Delete a column
err := userTable.DeleteColumn("age")
if err != nil {
    fmt.Println("Error deleting column:", err)
}
```
## Row Operations

The database provides various operations for working with individual rows and row collections, allowing you to
manipulate and analyze your data effectively.

### Working with Rows

These operations allow you to retrieve specific values from rows and convert row data to string format for display or
processing purposes.

```go
// Get specific value from row
row, err := userTable.GetRowById("1")
if err != nil {
    fmt.Println("Row not found:", err)
} else {
    name := row.SearchValue("name")
    fmt.Println("User name:", name)
}

// Convert row to string
fmt.Println("Row data:", row.String())
```
### Sorting Result

The database supports sorting operations on row collections, allowing you to order your data based on specific columns
in either ascending or descending order.

```go
rows := userTable.GetRows()

// Sort ascending
err := rows.OrderByAscend("name")
if err != nil {
    fmt.Println("Sort error:", err)
}

// Sort descending
err = rows.OrderByDescend("age")
if err != nil {
    fmt.Println("Sort error:", err)
}
```
