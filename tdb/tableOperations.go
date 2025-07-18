package tdb

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"slices"
	"sort"
	"strings"
)

// Table represents the interface for table operations in the database.
// It provides methods for manipulating and querying table data.
//
// Example usage:
//
//	// Create a new table
//	table := db.CreateTable("users", "id", "name", "email")
//
//	// Add a new row
//	table.AddValue("name", "John Doe")
//	table.AddValue("email", "john@example.com")
//
//	// Update a value
//	table.UpdateValue("name", "user123", "Jane Doe")
//
//	// Search for rows
//	row, _ := table.SearchOne("email", "john@example.com")
type Table interface {
	// AddValue adds a single value to the specified column in the table.
	// Returns an error if the column doesn't exist.
	//
	// Example usage:
	//
	//	err := table.AddValue("email", "user@example.com")
	//	if err != nil {
	//	    fmt.Println("Failed to add value:", err)
	//	}
	AddValue(column string, value string) error

	// AddValues adds multiple values to the table in one operation.
	// Automatically generates IDs for new rows.
	//
	// Example usage:
	//
	//	// Add multiple values at once
	//	table.AddValues("John Doe", "john@example.com", "active")
	AddValues(values ...string)

	// UpdateTableName changes the name of the table.
	//
	// Example usage:
	//
	//	// Rename table from "users" to "customers"
	//	table.UpdateTableName("customers")
	UpdateTableName(newName string)

	// UpdateColumnName changes the name of a column.
	// Returns an error if the column doesn't exist.
	//
	// Example usage:
	//
	//	err := table.UpdateColumnName("phone", "contact_number")
	UpdateColumnName(oldColumnName string, newColumnName string) error

	// UpdateValue updates a value in a specific row and column.
	// Returns an error if the row or column doesn't exist.
	//
	// Example usage:
	//
	//	err := table.UpdateValue("email", "user123", "new@example.com")
	UpdateValue(columnName string, id string, newValue string) error

	// DeleteRow removes a row from the table by its ID.
	// If cascade is true, also deletes related rows in other tables.
	//
	// Example usage:
	//
	//	// Delete row with cascade
	//	err := table.DeleteRow("user123", true)
	DeleteRow(id string, cascade bool) error

	// DeleteColumn removes a column from the table.
	// Returns an error if the column doesn't exist.
	//
	// Example usage:
	//
	//	err := table.DeleteColumn("unused_column")
	DeleteColumn(columnName string) error

	// GetRowById retrieves a specific row from the table using its ID.
	// Returns the row if found, or an error if the row doesn't exist.
	//
	// Example usage:
	//
	//	// Get a row by ID
	//	row, err := table.GetRowById("user_123")
	//	if err != nil {
	//	    fmt.Println("Row not found:", err)
	//	    return
	//	}
	//
	//	// Access row data
	//	name := row.SearchValue("name")
	GetRowById(id string) (Row, error)

	// GetRows returns all rows in the table as a Rows collection.
	//
	// Example usage:
	//
	//	// Get all rows from the table
	//	rows := table.GetRows()
	//
	//	// Print each row
	//	for _, row := range rows {
	//	    fmt.Println(row.String())
	//	}
	GetRows() Rows

	// GetColumns returns a slice containing all column names in the table.
	//
	// Example usage:
	//
	//	// Get all column names
	//	columns := table.GetColumns()
	//
	//	// Print column names
	//	for _, col := range columns {
	//	    fmt.Println(col)
	//	}
	GetColumns() []string

	// PrintTable prints the table contents to standard output.
	//
	// Example usage:
	//
	//	// Display the entire table
	//	table.PrintTable()
	PrintTable()

	// GetName returns the table name.
	//
	// Example usage:
	//
	//	tableName := table.GetName()
	//	fmt.Println("Current table:", tableName)
	GetName() string

	// SearchOne finds the first row where the column matches the value.
	// Returns an error if no match is found.
	//
	// Example usage:
	//
	//	// Find user by email
	//	row, err := table.SearchOne("email", "user@example.com")
	SearchOne(column string, value string) (Row, error)

	// SearchAll finds all rows where the specified column matches the given value.
	// Returns a Rows collection containing all matching rows.
	//
	// Example usage:
	//
	//	// Find all users with status "active"
	//	activeUsers := table.SearchAll("status", "active")
	//
	//	// Find orders from a specific date
	//	orders := table.SearchAll("order_date", "2025-07-18")
	SearchAll(column string, value string) Rows

	// SearchByForeignKey finds all related rows in other tables.
	// Returns an error if no foreign key relationships exist.
	//
	// Example usage:
	//
	//	// Find all orders for a customer
	//	related, err := table.SearchByForeignKey("customer_123")
	SearchByForeignKey(id string) ([]ComplexRow, error)

	// Internal methods used by the package implementation
	getSimpleName() string
	addValuesIdGenerationOff(values []string)
	table() table
	save()
}

// Rows represents a collection of Row objects.
// It provides methods for ordering and manipulating multiple rows.
//
// Example usage:
//
//	// Get all rows and sort by name
//	rows := table.GetRows()
//	rows.OrderByAscend("name")
//
//	// Sort by age in descending order
//	rows.OrderByDescend("age")
type Rows []Row

// Row represents a single row in a table.
// It contains the column names and corresponding values.
//
// Example usage:
//
//	// Get a specific value from a row
//	row, _ := table.GetRowById("123")
//	name := row.SearchValue("name")
//
//	// Print the row contents
//	fmt.Println(row.String())
type Row struct {
	columns []string
	value   string
}

// ComplexRow represents a row with its associated table and related rows.
// It is used for handling foreign key relationships and complex queries.
type ComplexRow struct {
	Table Table
	Rows  Rows
}

// table represents the internal structure of a database table.
// It contains the raw table name, columns, values and the raw table string representation.
type table struct {
	nameRaw  string
	columns  []string
	values   []Row
	rawTable string
}
type foreignKey struct {
	tableName string
	column    string
}

func (t *table) table() table {
	return *t
}
func (t *table) getSimpleName() string {
	name := strings.Trim(t.nameRaw, "-----")
	t.GetName()
	return name
}

// AddValue adds a single value to the specified column in the table and updates the raw table representation.
// It requires the column name and the value to be added as arguments.
// Returns an error if the column does not exist or if there is an issue during the value addition process.
func (t *table) AddValue(column string, value string) error {
	s, err := valueBuilder(*t, column, value)
	if err != nil {
		return err
	}

	t.rawTable = strings.Replace(t.rawTable, "!*!", s, 1)
	t.save()
	return nil
}

// PrintTable prints the raw string representation of the table to the standard output.
func (t *table) PrintTable() {
	fmt.Println(t.rawTable)
}

// GetName returns the raw name of the table as a string.
func (t *table) GetName() string {
	return t.nameRaw
}

// AddValues appends one or more values to the table and updates its internal representation.
// Each string in the `values` parameter represents a new row of data to be added.
func (t *table) AddValues(values ...string) {
	*t = addValues(*t, values, true)
}
func (t *table) addValuesIdGenerationOff(values []string) {
	*t = addValues(*t, values, false)
}
func (t *table) GetColumns() []string {
	return getColumns(t.rawTable)
}

// UpdateTableName changes the name of the table to the specified new name.
// The change is persisted to storage automatically.
//
// Example usage:
//
//	// Rename a table from "users" to "customers"
//	table.UpdateTableName("customers")
func (t *table) UpdateTableName(newName string) {
	formatName := strings.Replace(t.nameRaw, "-----", "", 2)
	formatName = formatName + "_End"
	formatName = fmt.Sprintf("-----%s-----", formatName)
	rawNewName := fmt.Sprintf("-----%s-----", newName)
	rawNewNameEnd := fmt.Sprintf("-----%s-----", newName+"_End")

	t.rawTable = strings.Replace(t.rawTable, t.nameRaw, rawNewName, 1)
	t.rawTable = strings.Replace(t.rawTable, formatName, rawNewNameEnd, 1)
	t.save()
}

// UpdateColumnName changes the name of a column from oldColumnName to newColumnName.
// Returns an error if the column doesn't exist in the table.
//
// Example usage:
//
//	// Rename column "phone" to "contact_number"
//	err := table.UpdateColumnName("phone", "contact_number")
//
//	if err != nil {
//	    fmt.Println("Column not found:", err)
//	}
func (t *table) UpdateColumnName(oldColumnName string, newColumnName string) error {
	index := slices.Index(t.columns, oldColumnName)
	if index == -1 {
		return &NotFoundError{itemName: "Column"}
	}
	t.columns[index] = newColumnName
	t.rawTable = strings.Replace(t.rawTable, oldColumnName, newColumnName, 1)
	t.save()
	return nil

}

// UpdateValue updates a value in a specific row and column of the table.
// Returns an error if either the column or row doesn't exist.
//
// Example usage:
//
//	// Update email address for user with ID "123"
//	err := table.UpdateValue("email", "123", "newemail@example.com")
//
//	// Update status of an order
//	err := table.UpdateValue("status", "order_456", "shipped")
func (t *table) UpdateValue(columnName string, id string, newValue string) error {

	index := slices.Index(t.columns, columnName)
	if index == -1 {
		return &NotFoundError{itemName: "Column"}
	}
	row, rowErr := t.GetRowById(id)
	if rowErr != nil {
		return rowErr
	}
	rowSlice := strings.Split(row.value, "|")
	rowSlice[index+1] = " " + newValue + " "
	row.value = strings.Join(rowSlice, "|")
	row.value = strings.Trim(row.value, " ")
	updateTable, err := updateRow(t.rawTable, id, row.value)
	if err != nil {
		return err
	}
	t.rawTable = updateTable
	t.save()
	return nil
}
func (t *table) GetRows() Rows {
	values := getRows(t.rawTable)
	return values
}
func (t *table) GetRowById(id string) (Row, error) {
	rows := getRows(t.rawTable)
	for i, row := range rows {
		s := strings.Split(row.value, " ")
		if s[1] == id {
			return rows[i], nil
		}
	}
	return Row{}, &NotFoundError{itemName: "Row"}
}

// DeleteRow removes a row from the table by its ID.
// If cascade is true, it also deletes any related rows in other tables that reference this row.
// Returns an error if the row doesn't exist or if there's an issue with cascade deletion.
//
// Example usage:
//
//	// Delete a single row
//	err := table.DeleteRow("123", false)
//
//	// Delete a row and its related records
//	err := table.DeleteRow("456", true)
func (t *table) DeleteRow(id string, cascade bool) error {
	newTable, err := deleteRow(*t, id)
	if err != nil {
		return err
	}
	if cascade {
		foreignKeyErr := deleteByForeignKey(*t, id)
		if foreignKeyErr != nil {
			return foreignKeyErr
		}
	}
	*t = newTable
	return nil
}
func (t *table) DeleteColumn(columnName string) error {
	index := slices.Index(t.columns, columnName)
	if index == -1 {
		return &NotFoundError{itemName: "Column"}
	}
	t.columns = slices.Delete(t.columns, index, index+1)
	position := fmt.Sprintf("[%d]", index-1)
	t.columns = slices.Replace(t.columns, index-1, index, position)
	if index < len(t.columns) {
		t.columns = slices.Delete(t.columns, index, index+1)
	} else {
		t.columns = slices.Delete(t.columns, index-1, index)
	}
	rawColumn := strings.Join(t.columns, " ")
	newTable := strings.Split(t.rawTable, "\n")
	newTable[2] = rawColumn
	t.rawTable = strings.Join(newTable, "\n")
	t.rawTable = deleteColumnData(t.rawTable, index)
	t.save()
	return nil
}
func (t *table) SearchOne(column string, value string) (Row, error) {
	rows := t.GetRows()
	index := slices.Index(rows[0].columns, column)
	for _, r := range rows {
		row := strings.Split(r.value, " ")
		if row[index] == value {
			return r, nil
		}
	}
	return Row{}, &NotFoundError{itemName: value}
}

func (t *table) SearchAll(column string, value string) Rows {
	return searchAll(*t, column, value)
}
func (t *table) SearchByForeignKey(id string) ([]ComplexRow, error) {
	keys, err := getTableForeignKey(*t)
	if err != nil {
		return nil, err
	}
	complexRows := &[]ComplexRow{}
	for _, key := range keys {
		tb, _ := getTableByName(key.tableName, false)
		result := searchAll(tb, key.column, id)
		removeStrConv(result)
		complexRow := &ComplexRow{
			Table: &tb,
			Rows:  result,
		}
		*complexRows = append(*complexRows, *complexRow)
	}

	return *complexRows, nil
}
func (t *table) save() {
	tables := getTables(false)
	for i, v := range tables {
		if v.GetName() == t.GetName() {
			tables[i] = *t
		}
	}
	saveTables(tables)
}
func deleteByForeignKey(tb table, id string) error {
	key, err := tb.SearchByForeignKey(id)
	if err != nil {
		return err
	}
	for i, v := range key {
		if v.Table.GetName() != tb.GetName() {
			newTb := v.Table.table()
			for _, j := range v.Rows {
				rowId := j.SearchValue("id")
				tbWithoutRow, deleteErr := deleteRow(newTb, rowId)
				if deleteErr != nil {
					return deleteErr
				}
				newTb = tbWithoutRow.table()
			}
			key[i].Table = &newTb
		}
	}
	return nil
}

func (r *Rows) String() string {
	s := make([]string, len(*r))
	return strings.Join(s, "\n")
}

// OrderByAscend sorts the rows in ascending order based on the specified column.
func (r *Rows) OrderByAscend(column string) error {
	order, err := orderBy(*r, column, true)
	if err != nil {
		return err
	}
	*r = order
	return err
}

// OrderByDescend sorts the rows in descending order based on the specified column.
func (r *Rows) OrderByDescend(column string) error {
	order, err := orderBy(*r, column, false)
	if err != nil {
		return err
	}
	*r = order
	return nil
}

// SearchValue retrieves the value for the specified column in the row.
func (r *Row) SearchValue(column string) string {
	index := slices.Index(r.columns, column)
	if index == -1 {
		return ""
	}
	s := strings.Split(r.value, " ")
	return s[index]
}
func (r *Row) String() string {
	return r.value
}

// deleteRow removes a row from the table by its ID and returns the updated table.
// Returns an error if the row is not found.
func deleteRow(tb table, id string) (table, error) {
	row, err := tb.GetRowById(id)
	if err != nil {
		return table{}, err
	}
	rowSlice := strings.Split(tb.rawTable, "\n")
	index := slices.Index(rowSlice, row.value)
	newRow := slices.Replace(rowSlice, index, index+1, "")
	newRow = removeEmptyIndex(newRow)
	rowString := strings.Join(newRow, "\n")
	tb.rawTable = "\n" + rowString + "\n"

	tb.save()
	return tb, nil

}

// removeStrConv converts special space characters back to normal spaces in row values.
// This is used when retrieving data that was stored with encoded spaces.
func removeStrConv(r Rows) Rows {
	for i := 0; i < len(r); i++ {
		value := strings.ReplaceAll(r[i].value, "U+0020", " ")
		r[i].value = value
	}
	return r
}

// searchAll finds all rows in the table where the specified column matches the given value.
// Returns an empty Rows collection if no matches are found.
func searchAll(tb table, column string, value string) Rows {
	var rowsResult Rows
	for _, row := range tb.values {
		v := row.SearchValue(column)
		if v == value {
			rowsResult = append(rowsResult, row)
		}
	}
	return rowsResult
}

// getTableForeignKey retrieves all foreign key relationships for the given table.
// Returns an error if no foreign keys are found.
func getTableForeignKey(tb table) ([]foreignKey, error) {
	if !isForeignKeyAvailable(tb.getSimpleName()) {
		return nil, &NotFoundError{itemName: "ForeignKey"}
	}
	var foreignKeys []foreignKey

	link, _ := getTableByName("Links", false)
	tb1 := searchAll(link, "table1", tb.getSimpleName())
	for _, row := range tb1 {
		tbName := row.SearchValue("table2")
		columnLink := row.SearchValue("columnLink2")
		foreignKeys = append(foreignKeys, foreignKey{tableName: tbName, column: columnLink})
	}
	return foreignKeys, nil
}
func isForeignKeyAvailable(tableName string) bool {
	tb, err := getTableByName("Links", false)
	if err != nil {
		return false
	}
	rows := tb.GetRows()
	for _, row := range rows {
		v := row.SearchValue("table1")
		if v == tableName {
			return true
		}
		v = row.SearchValue("table2")
		if v == tableName {
			return true
		}
	}
	return false
}
func orderBy(r Rows, column string, ascend bool) ([]Row, error) {
	newSlice := make([]Row, len(r))
	for i, row := range r {
		newSlice[i] = row
	}
	columns := r[0].columns
	columnIndex := slices.Index(columns, column)

	if columnIndex == -1 {
		return nil, &NotFoundError{itemName: "Column"}
	}

	sort.Slice(newSlice, func(i, j int) bool {
		s := strings.Split(newSlice[i].value, " ")
		s2 := strings.Split(newSlice[j].value, " ")
		if ascend {
			return s[columnIndex] > s2[columnIndex]
		} else {

			return s[columnIndex] < s2[columnIndex]
		}
	})
	return newSlice, nil
}
func deleteColumnData(table string, columnIndex int) string {
	tableSlice := strings.Split(table, "\n")
	length := len(tableSlice)
	for i := 3; i < length-3; i++ {
		columns := strings.Split(tableSlice[i], " ")
		columns = slices.Delete(columns, columnIndex, columnIndex+1)
		position := fmt.Sprintf("|%d|", columnIndex-1)
		columns = slices.Replace(columns, columnIndex-1, columnIndex, position)
		if columnIndex < len(columns) {

			columns = slices.Delete(columns, columnIndex, columnIndex+1)
		} else {

			columns = slices.Delete(columns, columnIndex-1, columnIndex)
		}
		tableSlice[i] = strings.Join(columns, " ")
	}
	return strings.Join(tableSlice, "\n")
}
func updateRow(table string, id string, newRow string) (string, error) {
	row := strings.Split(table, "\n")
	idExist := false
	for i := 3; i < len(row); i++ {
		s := strings.Split(row[i], "|")
		if strings.TrimSpace(s[2]) == id {
			row[i] = newRow
			idExist = true
			break
		}
	}
	if !idExist {
		return "", &NotFoundError{itemName: "id"}
	}
	return strings.Join(row, "\n"), nil
}

// valueBuilder constructs a new row value string for the table.
// It handles ID generation and space encoding for the new value.
func valueBuilder(table table, columnName string, value string) (string, error) {
	co := getColumns(table.rawTable)
	co = removeEmptyIndex(co)
	if !slices.Contains(co, columnName) {
		return "", &NotFoundError{itemName: "Column"}
	}
	count := len(co)
	co[0] = "[1]"
	co[1] = uuid.New().String()

	for i := 0; i < count; i += 2 {
		co[i] = strings.ReplaceAll(co[i], "[", "|")
		co[i] = strings.ReplaceAll(co[i], "]", "|")
	}

	for i := 3; i < count; i += 2 {
		if co[i] == columnName {
			value = strings.ReplaceAll(value, " ", "U+0020")
			co[i] = value
		} else {
			co[i] = "null"
		}
	}
	union := strings.Join(co, " ")
	result := union + "\n!*!"
	return result, nil
}

// valuesBuilder constructs multiple row value strings for the table.
// If idGenerate is true, it will generate new UUIDs for the rows.
func valuesBuilder(table string, values []Row, idGenerate bool) string {
	co := getColumns(table)
	co = removeEmptyIndex(co)
	count := len(co)

	for i := 0; i < count; i += 2 {
		co[i] = strings.ReplaceAll(co[i], "[", "|")
		co[i] = strings.ReplaceAll(co[i], "]", "|")
	}

	n := 1
	for _, r := range values {
		if idGenerate && n == 1 {
			co[1] = uuid.New().String()
			n = n + 2
		}

		if n > count {
			break
		}
		co[n] = strings.ReplaceAll(r.value, " ", "U+0020")
		n = n + 2

	}
	union := strings.Join(co, " ")
	result := union + "\n!*!"
	return result
}
func addValues(table table, values []string, idGenerate bool) table {
	rows := make([]Row, len(values))
	for i, v := range values {
		rows[i] = Row{
			columns: table.columns,
			value:   v,
		}
	}
	s := valuesBuilder(table.rawTable, rows, idGenerate)
	table.rawTable = strings.Replace(table.rawTable, "!*!", s, 1)
	table.save()
	return table
}

// saveTables writes the tables to the database file.
// If encryption is enabled, the data will be encrypted before saving.
func saveTables(tables []table) {
	var newTable string
	if len(tables) != 0 {
		newTable = addTableFrontiers(tables)
	}
	if encryptionKeyExist {
		newTable = must(globalEncoderKey.Encode(newTable))
	}
	errorHandler(os.WriteFile(dbName, []byte(newTable), 0666))
}
