## Foreign Key Relationships

---

The database system supports foreign key relationships between tables, allowing you to establish connections and
maintain referential integrity between related data. Below are the main operations available for working with foreign
keys.

### Adding Foreign Keys

You can create foreign key relationships between tables either by adding a single foreign key or multiple foreign keys
at once. This operation ensures referential integrity between tables by defining parent-child relationships.

```go
// Define foreign key relationship
foreignKey := tdb.ForeignKey{
    TableName:         "Orders",
    ColumnName:        "user_id",
    ForeignTableName:  "Users",
    ForeignColumnName: "id",
}

// Add single foreign key
err := db.AddForeignKey(foreignKey)
if err != nil {
    fmt.Println("Error adding foreign key:", err)
}

// Add multiple foreign keys
foreignKeys := []tdb.ForeignKey{
    {
        TableName:         "Orders",
        ColumnName:        "user_id",
        ForeignTableName:  "Users",
        ForeignColumnName: "id",
    },
    {
        TableName:         "OrderItems",
        ColumnName:        "order_id",
        ForeignTableName:  "Orders",
        ForeignColumnName: "id",
    },
}

err := db.AddForeignKeys(foreignKeys)
if err != nil {
    fmt.Println("Error adding foreign keys:", err)
}
```
### Querying with Foreign Keys

The foreign key querying functionality allows you to retrieve related records across tables using their established
relationships. This operation helps in fetching all associated data from connected tables based on foreign key values.

```go
// Search by foreign key relationships
complexRows, err := userTable.SearchByForeignKey("1")
if err != nil {
    fmt.Println("Foreign key search error:", err)
} else {
    for _, complexRow := range complexRows {
        fmt.Printf("Related table: %s\n", complexRow.Table.GetName())
        for _, row := range complexRow.Rows {
            fmt.Println("  Row:", row.String())
        }
    }
}
```