## Encryption Support

The database supports optional encryption of data at rest using AES-256 encryption. When enabled, all data stored in the
database file is encrypted, ensuring data security while maintaining the same functionality as an unencrypted database.

### How It Works

- Data is encrypted using AES-256 in CBC mode with PKCS7 padding
- Each write operation encrypts data before storing it to the file
- Each read operation decrypts data after reading from the file
- The encryption key is never stored in the database file
- The database structure remains the same, only the content is encrypted

### Usage

To create an encrypted database:

```go
// Create database with encryption
config := tdb.DbConfig{
    EncryptionKey: "your-secret-encryption-key",
    DatabaseName:  "encrypted_database.txt",
    DataConfig:    nil,
}

db, err := config.CreateDatabase()
if err != nil {
    fmt.Println("Error:", err)
    return
}

// Remove encryption from existing database
err = config.RemoveEncryption()
if err != nil {
    fmt.Println("Error removing encryption:", err)
}
```