# Text Database - Go

---
A simple, text-based database system written in Go that stores data in plain text files with optional encryption
support.
## Features

- **Text-based storage**: Data is stored in human-readable text files
- **Table operations**: Create, read, update, and delete tables
- **Row operations**: Insert, update, delete, and query rows
- **SQL-like queries**: Basic SELECT, INSERT, UPDATE, DELETE operations
- **Foreign key support**: Define relationships between tables
- **Encryption**: Optional encryption for data at rest
- **Migration support**: Database schema versioning (Beta)

# Table of contents

- [Installation](#features)
- [Quick Start](#quick-start)
- [File Structure](#file-structure)
- [Table Operations](docs/table-operations.md) 
- [Data Operations](docs/data-operation.md)
- [ForeignKey Operations](docs/foreignkey-operation.md)
- [Initial Data](docs/initial-data.md)
- [Encryption](docs/encryption.md)

## Installation

    go get github.com/sheymor21/text-database/tdb

## Quick Start

Getting started with the Text Database is straightforward. Follow these simple steps to create and initialize your database:
```go
package main

import (
"fmt"
"github.com/sheymor21/text-database/tdb"
)

func DoSomething() {
	
//Create a config
config := &tdb.DbConfig{}

// Create a database
db, err := config.CreateDatabase()
if err != nil {
fmt.Println(err)
return
}

// Your database is ready!
fmt.Println("Database created successfully!")
db.PrintTables()
}

```

## File Structure

The Text Database uses a custom, human-readable format to store data in plain text files. This structure is designed to be both easily parseable by the application and readable by humans for debugging or manual inspection.
```
////
-----Users-----
[1] id [2] name [3] email [4] age
|1| 1 |2| John Doe |3| john@example.com |4| 30
|1| 2 |2| Jane Smith |3| jane@example.com |4| 25
!*!
-----Users_End-----
////
```
## Contributing

Feel free to contribute to this project by submitting issues or pull requests.

## License

This project is open source and available under the MIT License.

## Author

- [sheymor21](https://github.com/sheymor21)