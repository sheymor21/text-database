package tdb

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
)

// Db interface defines the contract for database operations including table management and SQL queries
type Db interface {

	// GetName returns the name of the database file
	//
	// Example:
	//  dbName := db.GetName()
	//  fmt.Printf("Working with database: %s\n", dbName)
	GetName() string

	// GetTables returns all tables currently in the database
	//
	// Example:
	//  tables := db.GetTables()
	//  for _, table := range tables {
	//      fmt.Printf("Found table: %s\n", table.GetName())
	//  }
	GetTables() []Table

	// GetTableByName retrieves a specific table by its name
	// Returns error if table doesn't exist
	//
	// Example:
	//  table, err := db.GetTableByName("users")
	//  if err != nil {
	//      log.Fatal("Table not found:", err)
	//  }
	GetTableByName(name string) (Table, error)

	// PrintTables displays all tables and their contents to standard output
	//
	// Example:
	//  db.PrintTables()
	PrintTables()

	// NewTable creates a new table with specified name and columns
	// Returns the newly created table
	//
	// Example:
	//  columns := []string{"id", "name", "email"}
	//  usersTable := db.NewTable("users", columns)
	NewTable(name string, columns []string) Table

	// DeleteTable removes a table from the database
	// Returns error if table doesn't exist
	//
	// Example:
	//  err := db.DeleteTable("temporary_table")
	//  if err != nil {
	//      log.Fatal("Failed to delete table:", err)
	//  }
	DeleteTable(tableName string) error

	// AddForeignKey creates a foreign key relationship between two tables
	// Returns error if tables or columns don't exist
	//
	// Example:
	//  fk := ForeignKey{
	//      TableName: "orders",
	//      ColumnName: "user_id",
	//      ForeignTableName: "users",
	//      ForeignColumnName: "id",
	//  }
	//  err := db.AddForeignKey(fk)
	AddForeignKey(key ForeignKey) error

	// AddForeignKeys adds multiple foreign key relationships at once
	// Returns error if any operation fails
	//
	// Example:
	//  fks := []ForeignKey{
	//      {TableName: "orders", ColumnName: "user_id", ForeignTableName: "users", ForeignColumnName: "id"},
	//      {TableName: "orders", ColumnName: "product_id", ForeignTableName: "products", ForeignColumnName: "id"},
	//  }
	//  err := db.AddForeignKeys(fks)
	AddForeignKeys(keys []ForeignKey) error

	// FromSql executes an SQL query and returns the results
	// Returns error if query is invalid or execution fails
	//
	// Example:
	//  rows, err := db.FromSql("SELECT * FROM users WHERE age > 18")
	//  if err != nil {
	//      log.Fatal("Query failed:", err)
	//  }
	FromSql(sql string) (SqlRows, error)
}
type db struct {
	name   string
	tables []table
}

// DataConfig defines the structure for configuring table data with columns and values
type DataConfig struct {
	TableName string   // Name of the table
	Columns   []string // List of column names
	Values    []Values // List of row values
}

// Values represents a row of data as string values
type Values []string

// DbConfig defines the configuration for creating a new database
type DbConfig struct {
	EncryptionKey string       // Optional encryption key for database content
	DatabaseName  string       // Name of the database file
	DataConfig    []DataConfig // Initial data configuration for tables
}

// ForeignKey defines a relationship between two tables through their columns
type ForeignKey struct {
	TableName         string // Name of the source table
	ColumnName        string // Name of the source column
	ForeignTableName  string // Name of the referenced table
	ForeignColumnName string // Name of the referenced column
}

var encryptionKeyExist bool
var dbName string

// CreateDatabase creates a new database instance with the specified configuration
// Returns a database interface and any error encountered during creation
//
// Example:
//
//	config := DbConfig{
//		DatabaseName: "mydb.txt",
//		EncryptionKey: "secret",
//	}
//	db, err := config.CreateDatabase()
//	if err != nil {
//		log.Fatal(err)
//	}
func (c DbConfig) CreateDatabase() (Db, error) {
	checkDataErr := checkDataConfig(c.DataConfig)
	if checkDataErr != nil {
		return nil, checkDataErr
	}
	validationErr := validateDatabaseName(c.DatabaseName)
	if validationErr != nil {
		return nil, validationErr
	}

	dbName = c.DatabaseName
	if strings.TrimSpace(c.EncryptionKey) != "" {
		globalEncoderKey = *newSecureTextEncoder(c.EncryptionKey)
		encryptionKeyExist = true
	}
	if !isFileExist(c.DatabaseName) {
		errorHandler(os.WriteFile(c.DatabaseName, []byte{}, 0644))
		if c.DataConfig == nil {
			setDefaultData(c)
		}

	} else {
		data := string(must(os.ReadFile(c.DatabaseName)))
		if !isEncode(data) && encryptionKeyExist {
			encodeAndSave(data)
		}
	}

	if c.DataConfig != nil {
		newDb := setDatabaseData(c)
		return &newDb, nil
	}

	return &db{name: c.DatabaseName, tables: getTables(true)}, nil
}

// RemoveEncryption removes encryption from an encrypted database using the provided encryption key
// Returns an error if the database is not found or encryption key is invalid
//
// Example:
//
//	config := DbConfig{
//		DatabaseName: "mydb.txt",
//		EncryptionKey: "secret",
//	}
//	err := config.RemoveEncryption()
//	if err != nil {
//		log.Fatal(err)
//	}
func (c DbConfig) RemoveEncryption() error {
	if dbName == "" {
		return &NotFoundError{itemName: "Database"}
	}
	if strings.TrimSpace(c.EncryptionKey) != "" {
		data := string(must(os.ReadFile(c.DatabaseName)))
		if isEncode(data) {
			decodeAndSave(data)
			return nil
		}
	}
	return &NotFoundError{itemName: "EncryptionKey"}
}

// GetName returns the name of the database
//
// Example:
//
//	name := db.GetName()
//	fmt.Println("Database name:", name)
func (d *db) GetName() string {
	return dbName
}

// GetTables returns a list of all tables in the database
//
// Example:
//
//	tables := db.GetTables()
//	for _, table := range tables {
//		fmt.Println("Table:", table.GetName())
//	}
func (d *db) GetTables() []Table {
	tables := getTables(true)
	iTables := make([]Table, len(tables))
	for i, t := range tables {
		iTables[i] = &t
	}
	return iTables
}

// PrintTables prints all tables in the database to standard output
//
// Example:
//
//	db.PrintTables()
func (d *db) PrintTables() {
	tables := getTables(true)
	for _, t := range tables {
		fmt.Println(t.rawTable)
	}
}

// NewTable creates a new table with the specified name and columns
// Returns the created table interface
//
// Example:
//
//	columns := []string{"id", "name", "age"}
//	table := db.NewTable("users", columns)
func (d *db) NewTable(name string, columns []string) Table {
	t := &table{name, columns, nil, ""}
	tb := d.addTable(*t)
	return tb
}

// GetTableByName retrieves a table by its name
// Returns the table and an error if the table is not found
//
// Example:
//
//	table, err := db.GetTableByName("users")
//	if err != nil {
//		log.Fatal(err)
//	}
func (d *db) GetTableByName(name string) (Table, error) {
	tb, err := getTableByName(name, true)
	return &tb, err
}

// AddForeignKey adds a foreign key relationship between two tables
// Returns an error if the tables or columns don't exist
//
// Example:
//
//	key := ForeignKey{
//		TableName: "orders",
//		ColumnName: "user_id",
//		ForeignTableName: "users",
//		ForeignColumnName: "id",
//	}
//	err := db.AddForeignKey(key)
//	if err != nil {
//		log.Fatal(err)
//	}
func (d *db) AddForeignKey(key ForeignKey) error {
	tb, errTb := getTableByName(key.TableName, false)
	if errTb != nil {
		return &NotFoundError{itemName: "Table: " + key.TableName}
	}
	tbf, errTbf := getTableByName(key.ForeignTableName, false)
	if errTbf != nil {
		return &NotFoundError{itemName: "Table: " + key.ForeignTableName}
	}

	tbRows := getRows(tb.rawTable)
	tbfRows := getRows(tbf.rawTable)

	if !slices.Contains(tbRows[0].columns, key.ColumnName) {
		msg := fmt.Sprintf("Column: %s does not exist in table: %s", key.ColumnName, key.TableName)
		return &NotFoundError{itemName: msg}
	}
	if !slices.Contains(tbfRows[0].columns, key.ForeignColumnName) {
		msg := fmt.Sprintf("Column: %s does not exist in table: %s", key.ForeignColumnName, key.ForeignTableName)
		return &NotFoundError{itemName: msg}
	}
	if !isTableInDatabase("Links") {
		data := string(must(os.ReadFile(dbName)))
		linkAdded := string(linkTableLayout()) + data
		errorHandler(os.WriteFile(dbName, []byte(linkAdded), 0666))
	}

	linkTb, _ := getTableByName("Links", false)
	err := validateForeignKey(linkTb, key)
	if err != nil {
		return err
	}
	linkTb.AddValues(key.TableName, key.ColumnName, key.ForeignTableName, key.ForeignColumnName)
	return nil
}

// AddForeignKeys adds multiple foreign key relationships
// Returns an error if any of the foreign key operations fail
//
// Example:
//
//	keys := []ForeignKey{
//		{TableName: "orders", ColumnName: "user_id", ForeignTableName: "users", ForeignColumnName: "id"},
//		{TableName: "items", ColumnName: "order_id", ForeignTableName: "orders", ForeignColumnName: "id"},
//	}
//	err := db.AddForeignKeys(keys)
//	if err != nil {
//		log.Fatal(err)
//	}
func (d *db) AddForeignKeys(keys []ForeignKey) error {
	for _, key := range keys {
		err := d.AddForeignKey(key)
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *db) addTable(table table) Table {
	data := globalEncoderKey.readAndDecode(dbName)
	dataByte := []byte(data)
	raw := tableBuilder(table)
	dataByte = append(dataByte, []byte(raw)...)
	if encryptionKeyExist {
		dataEncode := must(globalEncoderKey.Encode(string(dataByte)))
		errorHandler(os.WriteFile(dbName, []byte(dataEncode), 0666))
		return must(d.GetTableByName(table.nameRaw))
	}
	errorHandler(os.WriteFile(dbName, dataByte, 0666))
	return must(d.GetTableByName(table.nameRaw))

}

// DeleteTable removes a table from the database by its name
// Returns an error if the table is not found
//
// Example:
//
//	err := db.DeleteTable("users")
//	if err != nil {
//		log.Fatal(err)
//	}
func (d *db) DeleteTable(tableName string) error {
	tables := getTables(true)
	tableNameRaw := fmt.Sprintf("-----%s-----", tableName)
	deleted := false
	for i, t := range tables {
		if t.nameRaw == tableNameRaw {
			tables = slices.Delete(tables, i, i+1)
			deleted = true
			break
		}
	}
	if !deleted {
		return &NotFoundError{itemName: tableName}
	}
	saveTables(tables)
	return nil
}

// FromSql executes an SQL query and returns the results
// Returns the query results and any error encountered during execution
//
// Example:
//
//	rows, err := db.FromSql("SELECT * FROM users WHERE age = 18")
//	if err != nil {
//		log.Fatal(err)
//	}
func (d *db) FromSql(sql string) (SqlRows, error) {
	return validateSql(*d, sql)
}

// getTableByName retrieves a table by its name from the database
// tableName: name of the table to retrieve
// strConv: flag to indicate if string conversion should be applied
// Returns the found table and any error encountered
func getTableByName(tableName string, strConv bool) (table, error) {
	tables := getTables(strConv)
	tableNameRaw := fmt.Sprintf("-----%s-----", tableName)

	for _, t := range tables {
		if t.nameRaw == tableNameRaw {
			return t, nil
		}

	}
	return table{}, &NotFoundError{itemName: "Table"}
}

// getTables retrieves all tables from the database
// strConv: flag to indicate if string conversion should be applied
// Returns a slice of all tables in the database
func getTables(strConv bool) []table {
	data := globalEncoderKey.readAndDecode(dbName)
	data = strings.ReplaceAll(data, "\r", "")
	if strConv {
		data = strings.ReplaceAll(data, "U+0020", " ")
	}
	s := strings.Split(data, "////")
	sif := removeEmptyIndex(s)
	tables := make([]table, len(sif))
	for i, t := range sif {
		name := getTableName(t)
		values := getRows(t)
		tables[i] = table{name, getColumns(t), values, t}
	}
	return tables
}

// tableBuilder constructs a string representation of a table
// table: the table structure to build
// Returns the string representation of the table
func tableBuilder(table table) string {
	if table.columns[0] != "id" {
		slices.Reverse(table.columns)
		table.columns = append(table.columns, "id")
		slices.Reverse(table.columns)
	}
	columnsRaw := columnsBuilder(table.columns)
	var builder strings.Builder
	name := fmt.Sprintf("\n-----%s-----\n", table.nameRaw)
	column := fmt.Sprintf(columnsRaw)
	end := fmt.Sprintf("\n!*!\n-----%s_End-----\n////", table.nameRaw)
	builder.WriteString(name)
	builder.WriteString(column)
	builder.WriteString(end)
	tableRaw := builder.String()
	if table.values != nil {
		values := valuesBuilder(tableRaw, table.values, true)
		tableRaw = strings.Replace(tableRaw, "!*!", values, 1)
	}
	return tableRaw
}

// columnsBuilder creates a formatted string of column definitions
// columns: slice of column names
// Returns a formatted string representing the columns
func columnsBuilder(columns []string) string {

	if len(columns) == 0 {
		return ""
	}

	var stringBuilder strings.Builder

	for i := 0; i < len(columns); i++ {
		stringBuilder.WriteString(fmt.Sprintf("[%d] %s ", i+1, columns[i]))
	}

	return strings.TrimSpace(stringBuilder.String())
}

// addTableFrontiers adds boundary markers between tables
// tables: slice of tables to process
// Returns a string with table boundaries added
func addTableFrontiers(tables []table) string {
	rawTables := make([]string, len(tables))
	rawTables[0] = "////"
	for _, t := range tables {
		rawTables = append(rawTables, t.rawTable)
		rawTables = append(rawTables, "////")
	}
	return strings.Join(rawTables, "")
}
func getTableName(rawTable string) string {
	tableName := strings.Split(rawTable, "-----")[1]
	tableNameRaw := fmt.Sprintf("-----%s-----", tableName)
	return tableNameRaw
}
func getColumns(rawTable string) []string {
	columns := strings.Split(rawTable, "\n")[2]
	columnsSlice := strings.Split(columns, " ")
	return columnsSlice
}
func getRows(table string) []Row {
	row := strings.Split(table, "\n")
	newRow := make([]Row, len(row)-6)
	columns := getColumns(table)
	n := 0
	for i := 3; i < len(row)-3; i++ {
		newRow[n].columns = columns
		newRow[n].value = row[i]
		n++
	}
	return newRow
}

// validateDatabaseName checks if the database name is valid
// name: database name to validate
// Returns an error if the name is invalid
func validateDatabaseName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("database name is required")
	} else {
		split := strings.Split(name, ".")
		if len(split) != 2 || split[1] != "txt" {
			return errors.New("database name must be a .txt file")
		}
	}
	return nil
}
func checkDataConfig(d []DataConfig) error {
	var wg sync.WaitGroup
	var builder strings.Builder
	errs := &[]error{}
	ch := make(chan error)
	for _, v := range d {
		wg.Add(1)
		go func(val DataConfig) {
			defer wg.Done()
			err := validateDataRequirement(val)
			ch <- err
		}(v)
		if err := <-ch; err != nil {
			*errs = append(*errs, <-ch)
		}
	}
	close(ch)
	wg.Wait()

	for _, err := range *errs {
		if err != nil {
			builder.WriteString(err.Error() + "\n")
		}
	}
	if len(*errs) > 0 {
		return errors.New(builder.String())
	}
	return nil
}

// validateDataRequirement checks if the data configuration meets requirements
// d: data configuration to validate
// Returns an error if requirements are not met
func validateDataRequirement(d DataConfig) error {
	for _, v := range d.Values {

		if d.TableName == "" {
			return errors.New("table name is required")
		}
		if len(d.Columns) == 0 {
			return errors.New("columns are required")
		}

		l := len(d.Columns)
		if strings.ToLower(d.Columns[0]) != "id" {
			l++
		}
		if l != len(v) {
			message := fmt.Sprintf("columns and values must have the same length, colums:%d != values:%d, table: %s", l, len(v), d.TableName)
			return errors.New(message)
		}
	}
	return nil
}
func addData(db db, d []DataConfig) {
	for _, v := range d {
		if !isTableInDatabase(v.TableName) {
			generateStaticData(db, v)
		} else {
			addStaticData(db, v)
		}
	}
}

// addStaticData adds predefined data to an existing table
// db: database instance
// v: data configuration containing the values to add
func addStaticData(db db, v DataConfig) {
	tb, _ := db.GetTableByName(v.TableName)
	for _, iv := range v.Values {
		if !areValuesInDatabase(v.TableName, iv[0]) {
			tb.addValuesIdGenerationOff(iv)
		}
	}
}

// generateStaticData creates a new table with predefined data
// db: database instance
// v: data configuration for table creation and data
func generateStaticData(db db, v DataConfig) {
	tb := db.NewTable(v.TableName, v.Columns)
	if v.Values != nil || len(v.Values) != 0 {
		for _, iv := range v.Values {
			tb.addValuesIdGenerationOff(iv)
		}
	}
}

// setDefaultData initializes the database with default data structure
// c: database configuration
func setDefaultData(c DbConfig) {
	if encryptionKeyExist {
		encodeAndSave(string(getLayout()))
	} else {
		errorHandler(os.WriteFile(c.DatabaseName, getLayout(), 0644))
	}
}

// isTableInDatabase checks if a table exists in the database
// tableName: name of the table to check
// Returns true if table exists, false otherwise
func isTableInDatabase(tableName string) bool {
	_, err := getTableByName(tableName, false)
	if err != nil {
		return false
	}
	return true
}

// areValuesInDatabase checks if specific values exist in a table
// tableName: name of the table to check
// value: value to search for
// Returns true if values exist, false otherwise
func areValuesInDatabase(tableName string, value string) bool {

	tb, err := getTableByName(tableName, false)
	if err != nil {
		return false
	}
	_, errR := tb.GetRowById(value)
	if errR != nil {
		return false
	}
	return true
}

// setDatabaseData initializes database with configured data
// c: database configuration containing initial data
// Returns initialized database instance
func setDatabaseData(c DbConfig) db {
	newDb := db{name: c.DatabaseName, tables: getTables(true)}
	addData(newDb, c.DataConfig)
	return newDb
}
func getLayout() []byte {
	layout := `////
-----Users-----
[1] id [2] name [3] age
|1| 1 |2| pedro |3| 32
|1| 2 |2| juan |3| 54
|1| 3 |2| carlos |3| 62
|1| 4 |2| manuel |3| 54
!*!
-----Users_End-----
////`
	return []byte(layout)
}

// validateForeignKey checks if a foreign key relationship is valid
// linkTb: link table containing relationships
// key: foreign key relationship to validate
// Returns error if relationship is invalid or already exists
func validateForeignKey(linkTb table, key ForeignKey) error {
	linkRows := getRows(linkTb.rawTable)
	for i := 1; i < len(linkRows); i++ {
		s := strings.Split(linkRows[i].value, "|")
		s = removeEmptyIndex(s)
		if strings.TrimSpace(s[3]) == key.TableName &&
			strings.TrimSpace(s[5]) == key.ColumnName &&
			strings.TrimSpace(s[7]) == key.ForeignTableName &&
			strings.TrimSpace(s[9]) == key.ForeignColumnName {
			return errors.New("foreign Key already exist")
		}
	}
	return nil
}
func linkTableLayout() []byte {
	layout := `////
-----Links-----
[1] id [2] table1 [3] columnLink1 [3] table2 [4] columnLink2
!*!
-----Links_End-----
////`
	return []byte(layout)
}
