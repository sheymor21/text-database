package pkg

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"slices"
	"strings"
	utilities2 "text-database/pkg/utilities"
)

type table struct {
	name    string
	columns []string
}

var database struct {
	name string
}

func CreateDatabase(databaseName string) {

	database.name = databaseName
	if !utilities2.IsFileExist(databaseName) {
		initData := utilities2.Must(os.ReadFile("internal/layout.txt"))
		utilities2.ErrorHandler(os.WriteFile(databaseName, initData, 0666))
	} else {

		log.Println(fmt.Sprintf("Already exists %s", databaseName))
	}

}
func AddValue(table string, column string, value string) {
	tableName := getTableName(table)
	tables := GetTables()
	for i, t := range tables {
		if strings.Contains(t, tableName) {
			s := valueBuilder(table, column, value)
			t = strings.Replace(t, "!*!", s, 1)
			tables[i] = t
			break
		}
	}
	newTable := addTableFrontiers(tables)
	utilities2.ErrorHandler(os.WriteFile(database.name, []byte(newTable), 0666))
}
func AddValues(table string, values []string) {
	tableName := getTableName(table)
	tables := GetTables()
	for i, t := range tables {
		if strings.Contains(t, tableName) {
			s := valuesBuilder(table, values)
			t = strings.Replace(t, "!*!", s, 1)
			tables[i] = t
			break
		}
	}
	utilities2.ErrorHandler(os.WriteFile(database.name, []byte(addTableFrontiers(tables)), 0666))
}
func GetTable(tableName string) string {
	tables := GetTables()
	tableNameRaw := fmt.Sprintf("-----%s-----", tableName)

	for _, t := range tables {
		if strings.Contains(t, tableNameRaw) {
			return t
		}

	}
	return ""
}
func GetTables() []string {
	data := utilities2.Must(os.ReadFile(database.name))
	dataString := string(data)
	dataString = strings.ReplaceAll(dataString, "U+0020", " ")
	s := strings.Split(dataString, "////")
	return utilities2.RemoveEmptyIndex(s)
}
func CreateTables(table table) {
	data := utilities2.Must(os.ReadFile(database.name))
	raw := tableBuilder(table)
	data = append(data, []byte(raw)...)
	utilities2.ErrorHandler(os.WriteFile(database.name, data, 0666))
}
func tableBuilder(table table) string {
	columnsRaw := columnsBuilder(table.columns)
	tableRaw := fmt.Sprintf("\n-----%s-----\n"+
		"[1] id %s"+
		"\n"+
		"!*!"+
		"-----%s_End-----\n////", table.name, columnsRaw, table.name)
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
func valueBuilder(table string, column string, value string) string {
	co := getColumns(table)
	count := len(co)
	co[0] = "[1]"
	co[1] = uuid.New().String()
	n := 3
	n2 := 0
	for range co {
		if n > count {
			break
		}
		if co[n] == column {
			value = strings.ReplaceAll(value, " ", "U+0020")
			co[n] = value
		} else {
			co[n] = "null"
		}
		co[n2] = strings.ReplaceAll(co[n2], "[", "|")
		co[n2] = strings.ReplaceAll(co[n2], "]", "|")
		n = n + 2
		n2 = n2 + 2
	}
	co = append(co, "\n!*!")
	return strings.Join(co, "")
}
func valuesBuilder(table string, values []string) string {
	co := getColumns(table)
	count := len(co)
	co[count-1] = "\n!*!"
	co[0] = "[1]"
	co[1] = uuid.New().String()
	n := 3
	n2 := 0
	for range co {

		for _, v := range values {
			if n > count {
				break
			}
			co[n] = strings.ReplaceAll(v, " ", "U+0020")
			co[n2] = strings.ReplaceAll(co[n2], "[", "|")
			co[n2] = strings.ReplaceAll(co[n2], "]", "|")
			n = n + 2
			n2 = n2 + 2

		}
	}
	co = append(co, "\n!*!")
	return strings.Join(co, " ")
}
func addTableFrontiers(tables []string) string {
	tables = append(tables, "////")
	slices.Reverse(tables)
	tables = append(tables, "////")
	return strings.Join(tables, "")
}
func getTableName(table string) string {
	tableName := strings.Split(table, "-----")[1]
	tableNameRaw := fmt.Sprintf("-----%s-----", tableName)
	return tableNameRaw
}
func getColumns(table string) []string {
	columns := strings.Split(table, "\n")[2]
	columnsSlice := strings.Split(columns, " ")
	return columnsSlice
}
func UpdateTableName(table string, newName string) {
	tables := GetTables()
	tableName := getTableName(table)
	formatName := strings.Replace(tableName, "-----", "", 2)
	formatName = formatName + "_End"
	formatName = fmt.Sprintf("-----%s-----", formatName)
	rawNewName := fmt.Sprintf("-----%s-----", newName)
	rawNewNameEnd := fmt.Sprintf("-----%s-----", newName+"_End")
	for i, t := range tables {
		if strings.Contains(t, tableName) {
			t = strings.Replace(t, tableName, rawNewName, 1)
			t = strings.Replace(t, formatName, rawNewNameEnd, 1)
			tables[i] = t
			break
		}
	}
	newTable := addTableFrontiers(tables)
	utilities2.ErrorHandler(os.WriteFile(database.name, []byte(newTable), 0666))
}
func UpdateColumnName(table string, oldColumnName string, newColumnName string) {
	tables := GetTables()
	columns := getColumns(table)
	for i, t := range tables {
		if strings.Contains(t, table) {
			for j, c := range columns {
				if c == oldColumnName {
					columns[j] = newColumnName
				}
			}
			t = strings.Replace(t, oldColumnName, newColumnName, 1)
			tables[i] = t
		}
	}

	newTable := addTableFrontiers(tables)
	utilities2.ErrorHandler(os.WriteFile(database.name, []byte(newTable), 0666))

}
func UpdateValue(table string, column string, id string, newValue string) {
	tables := GetTables()
	columns := getColumns(table)
	for i, t := range tables {
		if strings.Contains(t, table) {
			for n, c := range columns {
				if c == column {
					row := GetRow(t, id)
					rowSlice := strings.Split(row, "|")
					rowSlice[n+1] = newValue
					row = strings.Join(rowSlice, "|")
					newTableWithNewRow := UpdateRow(t, id, row)
					tables[i] = newTableWithNewRow
				}
			}
		}
	}
	newTable := addTableFrontiers(tables)
	utilities2.ErrorHandler(os.WriteFile(database.name, []byte(newTable), 0666))
}
func GetRow(table string, id string) string {
	row := strings.Split(table, "\n")
	for i, r := range row {
		if i >= 3 {
			s := strings.Split(r, "|")
			if s[2] == id {
				return r
			}
		}
	}
	return ""
}
func UpdateRow(table string, id string, newRow string) string {
	row := strings.Split(table, "\n")
	for i, r := range row {
		if i >= 3 {
			s := strings.Split(r, "|")
			if s[2] == id {
				row[i] = newRow
				break
			}
		}
	}
	return strings.Join(row, "\n")
}
func DeleteTable(table string) {
	tables := GetTables()
	for i, t := range tables {
		if t == table {
			slices.Delete(tables, i, i+1)
		}
	}
	newTable := addTableFrontiers(tables)
	utilities2.ErrorHandler(os.WriteFile(database.name, []byte(newTable), 0666))
}

func DeleteRow(table string, id string) {
	tables := GetTables()
	for i, t := range tables {
		if t == table {
			row := GetRow(t, id)
			rowSlice := strings.Split(t, "\n")
			index := slices.Index(rowSlice, row)
			newRow := slices.Replace(rowSlice, index, index+1, "")
			newRow = utilities2.RemoveEmptyIndex(newRow)
			rowString := strings.Join(newRow, "\n")
			tables[i] = "\n" + rowString + "\n"
		}
	}

	newTable := addTableFrontiers(tables)
	utilities2.ErrorHandler(os.WriteFile(database.name, []byte(newTable), 0666))
}
func DeleteColumn(table string, column string) {
	tables := GetTables()
	columns := getColumns(table)
	var rawTable string
	for n, c := range columns {
		if c == column {
			columns = slices.Delete(columns, n, n+1)
			position := fmt.Sprintf("[%d]", n-1)
			columns = slices.Replace(columns, n-1, n, position)
			if n < len(columns) {
				columns = slices.Delete(columns, n, n+1)
			} else {
				columns = slices.Delete(columns, n-1, n)
			}
			rawColumn := strings.Join(columns, " ")
			newTable := strings.Split(table, "\n")
			newTable[2] = rawColumn
			rawTable = strings.Join(newTable, "\n")
			rawTable = deleteColumnData(rawTable, n)
		}
	}
	for i, t := range tables {
		if t == table {
			tables[i] = rawTable
			fmt.Println(tables[i])
		}
	}
	tablesWithFrontier := addTableFrontiers(tables)
	utilities2.ErrorHandler(os.WriteFile(database.name, []byte(tablesWithFrontier), 0666))
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
