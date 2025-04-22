package pkg

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"slices"
	"strings"
	"text-database/pkg/utilities"
)

type Table struct {
	Name     string
	Columns  []string
	rawTable string
}

type Db struct {
	name   string
	tables []Table
}

func (d *Db) GetName() string {
	return dbName
}

func (d *Db) GetTables() []Table {
	return getTables()
}

var dbName string

func (d *Db) GetTable(name string) *Table {
	return getTable(name)
}

func CreateDatabase(databaseName string) {

	dbName = databaseName
	if !utilities.IsFileExist(databaseName) {
		initData := utilities.Must(os.ReadFile("internal/layout.txt"))
		utilities.ErrorHandler(os.WriteFile(databaseName, initData, 0666))
	} else {

		log.Println(fmt.Sprintf("Already exists %s", databaseName))
	}

}
func (table *Table) AddValue(column string, value string) {
	tables := getTables()
	for i, t := range tables {
		if strings.Contains(t.Name, table.Name) {
			s := valueBuilder(t, column, value)
			t.rawTable = strings.Replace(t.rawTable, "!*!", s, 1)
			tables[i] = t
			break
		}
	}
	newTable := addTableFrontiers(tables)
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTable), 0666))
}
func (table *Table) AddValues(values []string) {
	tables := getTables()
	for i, t := range tables {
		if strings.Contains(t.Name, table.Name) {
			s := valuesBuilder(table.rawTable, values)
			t.rawTable = strings.Replace(t.rawTable, "!*!", s, 1)
			tables[i] = t
			break
		}
	}
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(addTableFrontiers(tables)), 0666))
}
func getTable(tableName string) *Table {
	tables := getTables()
	tableNameRaw := fmt.Sprintf("-----%s-----", tableName)

	for _, t := range tables {
		if strings.Contains(t.rawTable, tableNameRaw) {
			return &t
		}

	}
	return nil
}
func getTables() []Table {
	data := utilities.Must(os.ReadFile(dbName))
	dataString := string(data)
	dataString = strings.ReplaceAll(dataString, "U+0020", " ")
	s := strings.Split(dataString, "////")
	sif := utilities.RemoveEmptyIndex(s)
	tables := make([]Table, len(sif))
	for i, t := range sif {
		name := getTableName(t)
		tables[i] = Table{name, getColumns(t), t}
	}
	return tables
}
func (d *Db) AddTable(table Table) {
	data := utilities.Must(os.ReadFile(dbName))
	raw := tableBuilder(table)
	data = append(data, []byte(raw)...)
	utilities.ErrorHandler(os.WriteFile(dbName, data, 0666))
}
func tableBuilder(table Table) string {
	columnsRaw := columnsBuilder(table.Columns)
	tableRaw := fmt.Sprintf("\n-----%s-----\n"+
		"[1] id %s"+
		"\n!*!"+
		"\n-----%s_End-----\n////", table.Name, columnsRaw, table.Name)
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
func valueBuilder(table Table, columnName string, value string) string {
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
func addTableFrontiers(tables []Table) string {
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
func (table *Table) UpdateTableName(newName string) {
	tables := getTables()
	formatName := strings.Replace(table.Name, "-----", "", 2)
	formatName = formatName + "_End"
	formatName = fmt.Sprintf("-----%s-----", formatName)
	rawNewName := fmt.Sprintf("-----%s-----", newName)
	rawNewNameEnd := fmt.Sprintf("-----%s-----", newName+"_End")
	for i, t := range tables {
		if strings.Contains(t.Name, table.Name) {
			t.rawTable = strings.Replace(t.rawTable, table.Name, rawNewName, 1)
			t.rawTable = strings.Replace(t.rawTable, formatName, rawNewNameEnd, 1)
			tables[i] = t
			break
		}
	}
	newTable := addTableFrontiers(tables)
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTable), 0666))
}
func (table *Table) UpdateColumnName(oldColumnName string, newColumnName string) {
	tables := getTables()
	for i, t := range tables {
		if strings.Contains(t.Name, table.Name) {
			for j, c := range table.Columns {
				if c == oldColumnName {
					table.Columns[j] = newColumnName
				}
			}
			t.rawTable = strings.Replace(t.rawTable, oldColumnName, newColumnName, 1)
			tables[i] = t
		}
	}

	newTable := addTableFrontiers(tables)
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTable), 0666))

}
func (table *Table) UpdateValue(columnName string, id string, newValue string) {
	tables := getTables()
	for i, t := range tables {
		if strings.Contains(t.Name, table.Name) {
			for n := 3; n < len(table.Columns); n += 2 {
				if table.Columns[n] == columnName {
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
func (table *Table) GetRow(id string) string {
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
func (table *Table) DeleteTable() {
	tables := getTables()
	for i, t := range tables {
		if t.Name == table.Name {
			tables = slices.Delete(tables, i, i+1)
		}
	}

	newTable := addTableFrontiers(tables)
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTable), 0666))
}

func (table *Table) DeleteRow(id string) {
	tables := getTables()
	for i, t := range tables {
		if t.Name == table.Name {
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
func (table *Table) DeleteColumn(columnName string) {
	tables := getTables()
	var rawTable string
	for n, c := range table.Columns {
		if c == columnName {
			table.Columns = slices.Delete(table.Columns, n, n+1)
			position := fmt.Sprintf("[%d]", n-1)
			table.Columns = slices.Replace(table.Columns, n-1, n, position)
			if n < len(table.Columns) {
				table.Columns = slices.Delete(table.Columns, n, n+1)
			} else {
				table.Columns = slices.Delete(table.Columns, n-1, n)
			}
			rawColumn := strings.Join(table.Columns, " ")
			newTable := strings.Split(table.rawTable, "\n")
			newTable[2] = rawColumn
			rawTable = strings.Join(newTable, "\n")
			rawTable = deleteColumnData(rawTable, n)
		}
	}
	for i, t := range tables {
		if t.Name == table.Name {
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
