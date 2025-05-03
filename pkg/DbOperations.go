package pkg

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"slices"
	"strings"
	"text-database/pkg/utilities"
)

type Table interface {
	AddValue(column string, value string) (Table, error)
	AddValues(values []string) Table
	UpdateTableName(newName string) Table
	UpdateColumnName(oldColumnName string, newColumnName string) (Table, error)
	UpdateValue(columnName string, id string, newValue string) (Table, error)
	DeleteRow(id string) (Table, error)
	DeleteColumn(columnName string) (Table, error)
	GetRowById(id string) (string, error)
	GetRows() []string
	GetColumns() []string
	PrintTable()
	GetName() string
}
type Db interface {
	GetName() string
	GetTables() []Table
	GetTableByName(name string) (Table, error)
	PrintTables()
	NewTable(name string, columns []string) Table
	DeleteTable(tableName string)
}
type table struct {
	name     string
	columns  []string
	data     string
	rawTable string
}
type db struct {
	name   string
	tables []table
}

func (d db) GetName() string {
	return dbName
}
func (d db) GetTables() []Table {
	tables := getTables()
	iTables := make([]Table, len(tables))
	for i, t := range tables {
		iTables[i] = &t
	}
	return iTables
}
func (d db) PrintTables() {
	tables := getTables()
	for _, t := range tables {
		fmt.Println(t.rawTable)
	}
}
func (d db) NewTable(name string, columns []string) Table {
	t := &table{name, columns, "", ""}
	tb := d.addTable(*t)
	return tb
}

var dbName string

func (d db) GetTableByName(name string) (Table, error) {
	return getTableByName(name)
}
func CreateDatabase(databaseName string) Db {

	dbName = databaseName
	if !utilities.IsFileExist(databaseName) {
		initData, err := os.ReadFile("internal/layout.txt")
		if err != nil {
			initData = utilities.Must(os.ReadFile("../internal/layout.txt"))
		}
		utilities.ErrorHandler(os.WriteFile(databaseName, initData, 0666))
	}
	return db{name: databaseName, tables: getTables()}

}
func (table table) AddValue(column string, value string) (Table, error) {
	tables := getTables()
	s, err := valueBuilder(table, column, value)
	if err != nil {
		return nil, err
	}

	table.rawTable = strings.Replace(table.rawTable, "!*!", s, 1)
	for i, t := range tables {
		if strings.Contains(t.name, table.name) {
			tables[i] = table
			break
		}
	}
	newTable := addTableFrontiers(tables)
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTable), 0666))
	return table, nil
}
func (table table) PrintTable() {
	fmt.Println(table.rawTable)
}
func (table table) GetName() string {
	return table.name
}
func (table table) AddValues(values []string) Table {
	s := valuesBuilder(table.rawTable, values)
	table.rawTable = strings.Replace(table.rawTable, "!*!", s, 1)
	tables := getTables()
	for i, t := range tables {
		if strings.Contains(t.name, table.name) {
			tables[i] = table
			break
		}
	}
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(addTableFrontiers(tables)), 0666))
	return table
}
func (table table) GetColumns() []string {
	return getColumns(table.rawTable)
}
func getTableByName(tableName string) (table, error) {
	tables := getTables()
	tableNameRaw := fmt.Sprintf("-----%s-----", tableName)

	for _, t := range tables {
		if strings.Contains(t.rawTable, tableNameRaw) {
			return t, nil
		}

	}
	return table{}, &NotFoundError{itemName: "Table"}
}
func getTables() []table {
	data := utilities.Must(os.ReadFile(dbName))
	dataString := string(data)
	dataString = strings.ReplaceAll(dataString, "\r", "")
	dataString = strings.ReplaceAll(dataString, "U+0020", " ")
	s := strings.Split(dataString, "////")
	sif := utilities.RemoveEmptyIndex(s)
	tables := make([]table, len(sif))
	for i, t := range sif {
		name := getTableName(t)
		rowData := getData(t)
		tables[i] = table{name, getColumns(t), rowData, t}
	}
	return tables
}
func (d db) addTable(table table) Table {
	data := utilities.Must(os.ReadFile(dbName))
	raw := tableBuilder(table)
	data = append(data, []byte(raw)...)
	utilities.ErrorHandler(os.WriteFile(dbName, data, 0666))
	return utilities.Must(d.GetTableByName(table.name))

}
func tableBuilder(table table) string {
	columnsRaw := columnsBuilder(table.columns)
	tableRaw := fmt.Sprintf("\n-----%s-----\n"+
		"[1] id %s"+
		"\n!*!"+
		"\n-----%s_End-----\n////", table.name, columnsRaw, table.name)
	return tableRaw
}
func columnsBuilder(columns []string) string {

	if len(columns) == 0 {
		return ""
	}

	var stringBuilder strings.Builder
	count := len(columns)
	count = count + 2
	for i := 2; i < count; i++ {
		stringBuilder.WriteString(fmt.Sprintf("[%d] %s ", i, columns[i-2]))
	}
	return stringBuilder.String()
}
func valueBuilder(table table, columnName string, value string) (string, error) {
	co := getColumns(table.rawTable)
	co = utilities.RemoveEmptyIndex(co)
	if !slices.Contains(co, columnName) {
		return "", &NotFoundError{itemName: "Column"}
	}
	count := len(co)
	co[0] = "[1]"
	co[1] = " " + uuid.New().String() + " "

	for i := 0; i < count; i += 2 {
		co[i] = strings.ReplaceAll(co[i], "[", "|")
		co[i] = strings.ReplaceAll(co[i], "]", "|")
	}

	for i := 3; i < count; i += 2 {
		if co[i] == columnName {
			value = strings.ReplaceAll(value, " ", "U+0020")
			co[i] = " " + value + " "
		} else {
			co[i] = " null "
		}
	}
	co = append(co, "\n!*!")
	return strings.Join(co, ""), nil
}
func valuesBuilder(table string, values []string) string {
	co := getColumns(table)
	co = utilities.RemoveEmptyIndex(co)
	count := len(co)
	co[0] = "[1]"
	co[1] = uuid.New().String()

	for i := 0; i < count; i += 2 {
		co[i] = strings.ReplaceAll(co[i], "[", "|")
		co[i] = strings.ReplaceAll(co[i], "]", "|")
	}

	n := 3
	for _, v := range values {

		if n > count {
			break
		}
		co[n] = strings.ReplaceAll(v, " ", "U+0020")
		n = n + 2

	}
	co = append(co, "\n!*!")
	return strings.Join(co, " ")
}
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
func (table table) UpdateTableName(newName string) Table {
	formatName := strings.Replace(table.name, "-----", "", 2)
	formatName = formatName + "_End"
	formatName = fmt.Sprintf("-----%s-----", formatName)
	rawNewName := fmt.Sprintf("-----%s-----", newName)
	rawNewNameEnd := fmt.Sprintf("-----%s-----", newName+"_End")

	table.name = rawNewName
	table.rawTable = strings.Replace(table.rawTable, table.name, rawNewName, 1)
	table.rawTable = strings.Replace(table.rawTable, formatName, rawNewNameEnd, 1)
	tables := getTables()
	for i, t := range tables {
		if t.name == table.name {
			tables[i].name = table.name
			tables[i] = table
			break
		}
	}
	newTable := addTableFrontiers(tables)
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTable), 0666))
	return table
}
func (table table) UpdateColumnName(oldColumnName string, newColumnName string) (Table, error) {
	if !slices.Contains(table.columns, oldColumnName) {
		return nil, &NotFoundError{itemName: "Column"}
	}
	tables := getTables()
	for i, c := range table.columns {
		if c == oldColumnName {
			table.columns[i] = newColumnName
			break
		}
	}

	table.rawTable = strings.Replace(table.rawTable, oldColumnName, newColumnName, 1)
	for i, t := range tables {
		if t.name == table.name {
			tables[i] = table
			break
		}
	}

	newTable := addTableFrontiers(tables)
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTable), 0666))
	return table, nil

}
func (table table) UpdateValue(columnName string, id string, newValue string) (Table, error) {

	if !slices.Contains(table.columns, columnName) {
		return nil, &NotFoundError{itemName: "Column"}
	}
	for i := 3; i < len(table.columns); i += 2 {
		if strings.TrimSpace(table.columns[i]) == columnName {
			row, rowErr := table.GetRowById(id)
			if rowErr != nil {
				return nil, rowErr
			}
			rowSlice := strings.Split(row, "|")
			rowSlice[i+1] = " " + newValue + " "
			row = strings.Join(rowSlice, "|")
			updateTable, err := updateRow(table.rawTable, id, row)
			if err != nil {
				return nil, err
			}
			table.rawTable = updateTable
			break
		}
	}

	tables := getTables()
	for i, t := range tables {
		if t.name == table.name {
			tables[i].rawTable = table.rawTable
			break
		}
	}
	newTable := addTableFrontiers(tables)
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTable), 0666))
	return table, nil
}
func (table table) GetRows() []string {
	rows := strings.Split(table.rawTable, "\n")
	sRows := make([]string, len(rows)-3)

	for i := 3; i < len(rows)-3; i++ {
		sRows[i-3] = rows[i]
	}
	return utilities.RemoveEmptyIndex(sRows)
}
func (table table) GetRowById(id string) (string, error) {
	row := strings.Split(table.rawTable, "\n")
	for i := 3; i < len(row)-3; i++ {
		s := strings.Split(row[i], "|")
		if strings.TrimSpace(s[2]) == id {
			return row[i], nil
		}
	}
	return "", &NotFoundError{itemName: "Row"}
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
		return "", &NotFoundError{itemName: "Id"}
	}
	return strings.Join(row, "\n"), nil
}
func (d db) DeleteTable(tableName string) {
	tables := getTables()
	tableNameRaw := fmt.Sprintf("-----%s-----", tableName)
	for i, t := range tables {
		if t.name == tableNameRaw {
			tables = slices.Delete(tables, i, i+1)
			break
		}
	}

	newTable := addTableFrontiers(tables)
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTable), 0666))
}
func (table table) DeleteRow(id string) (Table, error) {

	row, err := table.GetRowById(id)
	if err != nil {
		return nil, err
	}
	rowSlice := strings.Split(table.rawTable, "\n")
	index := slices.Index(rowSlice, row)
	newRow := slices.Replace(rowSlice, index, index+1, "")
	newRow = utilities.RemoveEmptyIndex(newRow)
	rowString := strings.Join(newRow, "\n")
	table.rawTable = "\n" + rowString + "\n"

	tables := getTables()
	for i, t := range tables {
		if t.name == table.name {
			tables[i] = table
			break
		}
	}

	newTable := addTableFrontiers(tables)
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTable), 0666))
	return table, nil
}
func (table table) DeleteColumn(columnName string) (Table, error) {
	if !slices.Contains(table.columns, columnName) {
		return nil, &NotFoundError{itemName: "Column"}
	}
	for n, c := range table.columns {
		if c == columnName {
			table.columns = slices.Delete(table.columns, n, n+1)
			position := fmt.Sprintf("[%d]", n-1)
			table.columns = slices.Replace(table.columns, n-1, n, position)
			if n < len(table.columns) {
				table.columns = slices.Delete(table.columns, n, n+1)
			} else {
				table.columns = slices.Delete(table.columns, n-1, n)
			}
			rawColumn := strings.Join(table.columns, " ")
			newTable := strings.Split(table.rawTable, "\n")
			newTable[2] = rawColumn
			table.rawTable = strings.Join(newTable, "\n")
			table.rawTable = deleteColumnData(table.rawTable, n)
			break
		}
	}
	tables := getTables()
	for i, t := range tables {
		if t.name == table.name {
			tables[i] = table
			break
		}
	}
	tablesWithFrontier := addTableFrontiers(tables)
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(tablesWithFrontier), 0666))
	return table, nil
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
func getTableIndex(table string) int {
	tables := getTables()
	for i, t := range tables {
		if t.rawTable == table {
			return i
		}
	}
	return -1
}
func getData(table string) string {
	row := strings.Split(table, "\n")
	newRow := make([]string, len(row))
	for i := 3; i < len(row); i++ {
		newRow[i-3] = row[i]
	}
	return strings.Join(newRow, "\n")
}
