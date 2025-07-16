## Database Configuration with Initial Data

---

The database configuration allows you to initialize your database with predefined data structures and values. You can
specify table names, column definitions, and initial row values during the database creation process. The configuration
supports optional encryption for data security.

Below is an example of how to configure and create a database with initial data:

```go
// Define initial data structure
dataConfig := []tdb.DataConfig{
    {
        TableName: "Users",
        Columns:   []string{"name", "email", "age"},
        Values: []tdb.Values{
            {"1", "John Doe", "john@example.com", "30"},
            {"2", "Jane Smith", "jane@example.com", "25"},
        },
    },
    {
        TableName: "Orders",
        Columns:   []string{"user_id", "product", "amount"},
        Values: []tdb.Values{
            {"1", "1", "Laptop", "999.99"},
            {"2", "2", "Mouse", "29.99"},
        },
    },
}

// Create database with initial data
config := tdb.DbConfig{
    EncryptionKey: "your-secret-key", // Optional
    DatabaseName:  "database.txt",
    DataConfig:    dataConfig,
}

db, err := config.CreateDatabase()
if err != nil {
    fmt.Println("Database creation error:", err)
    return
}
```
