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
	AddValue(column string, value string)
	AddValues(values []string)
	UpdateTableName(newName string)
	UpdateColumnName(oldColumnName string, newColumnName string)
	UpdateValue(columnName string, id string, newValue string)
	DeleteTable()
	DeleteRow(id string)
	DeleteColumn(columnName string)
	GetRow(id string) string
	PrintTable()
}
type Db interface {
	GetName() string
	GetTables() []Table
	GetTableByName(name string) Table
	AddTable(table table)
	PrintTables()
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

var dbName string

func (d db) GetTableByName(name string) Table {
	return getTable(name)
}
func CreateDatabase(databaseName string) Db {

	dbName = databaseName
	if !utilities.IsFileExist(databaseName) {
		initData := utilities.Must(os.ReadFile("internal/layout.txt"))
		utilities.ErrorHandler(os.WriteFile(databaseName, initData, 0666))
	}
	return db{}

}
func (table table) AddValue(column string, value string) {
	tables := getTables()
	for i, t := range tables {
		if strings.Contains(t.name, table.name) {
			s := valueBuilder(t, column, value)
			t.rawTable = strings.Replace(t.rawTable, "!*!", s, 1)
			tables[i] = t
			break
		}
	}
	newTable := addTableFrontiers(tables)
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTable), 0666))
}
func (table table) PrintTable() {
	fmt.Println(table.rawTable)
}
func (table table) AddValues(values []string) {
	tables := getTables()
	for i, t := range tables {
		if strings.Contains(t.name, table.name) {
			s := valuesBuilder(table.rawTable, values)
			t.rawTable = strings.Replace(t.rawTable, "!*!", s, 1)
			tables[i] = t
			break
		}
	}
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(addTableFrontiers(tables)), 0666))
}
func getTable(tableName string) Table {
	tables := getTables()
	tableNameRaw := fmt.Sprintf("-----%s-----", tableName)

	for _, t := range tables {
		if strings.Contains(t.rawTable, tableNameRaw) {
			return t
		}

	}
	return nil
}
func getTables() []table {
	data := utilities.Must(os.ReadFile(dbName))
	dataString := string(data)
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
func (d db) AddTable(table table) {
	data := utilities.Must(os.ReadFile(dbName))
	raw := tableBuilder(table)
	data = append(data, []byte(raw)...)
	utilities.ErrorHandler(os.WriteFile(dbName, data, 0666))
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
func valueBuilder(table table, columnName string, value string) string {
	co := getColumns(table.rawTable)
	co = utilities.RemoveEmptyIndex(co)
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
	return strings.Join(co, "")
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
func (table table) UpdateTableName(newName string) {
	tables := getTables()
	formatName := strings.Replace(table.name, "-----", "", 2)
	formatName = formatName + "_End"
	formatName = fmt.Sprintf("-----%s-----", formatName)
	rawNewName := fmt.Sprintf("-----%s-----", newName)
	rawNewNameEnd := fmt.Sprintf("-----%s-----", newName+"_End")
	for i, t := range tables {
		if strings.Contains(t.name, table.name) {
			t.rawTable = strings.Replace(t.rawTable, table.name, rawNewName, 1)
			t.rawTable = strings.Replace(t.rawTable, formatName, rawNewNameEnd, 1)
			tables[i] = t
			break
		}
	}
	newTable := addTableFrontiers(tables)
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTable), 0666))
}
func (table table) UpdateColumnName(oldColumnName string, newColumnName string) {
	tables := getTables()
	for i, t := range tables {
		if strings.Contains(t.name, table.name) {
			for j, c := range table.columns {
				if c == oldColumnName {
					table.columns[j] = newColumnName
				}
			}
			t.rawTable = strings.Replace(t.rawTable, oldColumnName, newColumnName, 1)
			tables[i] = t
		}
	}

	newTable := addTableFrontiers(tables)
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTable), 0666))

}
func (table table) UpdateValue(columnName string, id string, newValue string) {
	tables := getTables()
	for i, t := range tables {
		if strings.Contains(t.name, table.name) {
			for n := 3; n < len(table.columns); n += 2 {
				if table.columns[n] == columnName {
					row := table.GetRow(id)
					rowSlice := strings.Split(row, "|")
					rowSlice[n+1] = " " + newValue + " "
					row = strings.Join(rowSlice, "|")
					newTableWithNewRow := updateRow(t.rawTable, id, row)
					tables[i].rawTable = newTableWithNewRow
				}
			}
		}
	}
	newTable := addTableFrontiers(tables)
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTable), 0666))
}
func (table table) GetRow(id string) string {
	row := strings.Split(table.rawTable, "\n")
	for i, r := range row {
		if i >= 3 {
			s := strings.Split(r, "|")

			if strings.TrimSpace(s[2]) == id {
				return r
			}
		}
	}
	return ""
}
func updateRow(table string, id string, newRow string) string {
	row := strings.Split(table, "\n")
	for i := 3; i < len(row); i++ {
		s := strings.Split(row[i], "|")
		if strings.TrimSpace(s[2]) == id {
			row[i] = newRow
			break
		}
	}
	return strings.Join(row, "\n")
}
func (table table) DeleteTable() {
	tables := getTables()
	for i, t := range tables {
		if t.name == table.name {
			tables = slices.Delete(tables, i, i+1)
		}
	}

	newTable := addTableFrontiers(tables)
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTable), 0666))
}
func (table table) DeleteRow(id string) {
	tables := getTables()
	for i, t := range tables {
		if t.name == table.name {
			row := table.GetRow(id)
			rowSlice := strings.Split(t.rawTable, "\n")
			index := slices.Index(rowSlice, row)
			newRow := slices.Replace(rowSlice, index, index+1, "")
			newRow = utilities.RemoveEmptyIndex(newRow)
			rowString := strings.Join(newRow, "\n")
			tables[i].rawTable = "\n" + rowString + "\n"
		}
	}

	newTable := addTableFrontiers(tables)
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTable), 0666))
}
func (table table) DeleteColumn(columnName string) {
	tables := getTables()
	var rawTable string
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
			rawTable = strings.Join(newTable, "\n")
			rawTable = deleteColumnData(rawTable, n)
		}
	}
	for i, t := range tables {
		if t.name == table.name {
			tables[i].rawTable = rawTable
		}
	}
	tablesWithFrontier := addTableFrontiers(tables)
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(tablesWithFrontier), 0666))
}
func deleteColumnData(table string, columnIndex int) string {
	tableSlice := strings.Split(table, "\n")
	length := len(tableSlice)
	for i, _ := range tableSlice {
		if i > 2 && i < length-3 {
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
