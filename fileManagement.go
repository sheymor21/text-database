package main

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"strings"
	"text-database/utilities"
)

type table struct {
	name    string
	columns []string
}

var database struct {
	name string
}

func createFile(databaseName string) {
	database.name = databaseName
	if utilities.IsFileExist(databaseName) {
		log.Println(fmt.Sprintf("Already exists %s", databaseName))
		return
	}

	initData := utilities.Must(os.ReadFile(databaseName))

	utilities.ErrorHandler(os.WriteFile(databaseName, initData, 0666))
}
func addValue(table string, column string, value string) {
	tableName := getTableName(table)
	tables := getTables()
	for i, t := range tables {
		if strings.Contains(t, tableName) {
			s := valueBuilder(table, column, value)
			t = strings.Replace(t, "!*!", s, 1)
			tables[i] = t
			break
		}
	}
	newTable := addTableFrontiers(tables)
	utilities.ErrorHandler(os.WriteFile(database.name, []byte(newTable), 0666))
}
func addValues(table string, values []string) {
	tableName := getTableName(table)
	tables := getTables()
	for i, t := range tables {
		if strings.Contains(t, tableName) {
			s := valuesBuilder(table, values)
			t = strings.Replace(t, "!*!", s, 1)
			tables[i] = t
			break
		}
	}
	utilities.ErrorHandler(os.WriteFile(database.name, []byte(addTableFrontiers(tables)), 0666))
}
func getTable(tableName string) string {
	tables := getTables()
	tableNameRaw := fmt.Sprintf("-----%s-----", tableName)

	for _, t := range tables {
		if strings.Contains(t, tableNameRaw) {
			return strings.Replace(t, "U+0020", " ", 1)
		}

	}
	return ""
}
func getTables() []string {
	data := utilities.Must(os.ReadFile(database.name))
	dataString := string(data)
	dataString = strings.ReplaceAll(dataString, "U+0020", " ")
	s := strings.Split(dataString, "////")
	return utilities.RemoveEmptyIndex(s)
}
func createTables(table table) {
	data := utilities.Must(os.ReadFile(database.name))
	raw := tableBuilder(table)
	data = append(data, []byte(raw)...)
	utilities.ErrorHandler(os.WriteFile(database.name, data, 0666))
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
			co[n] = strings.ReplaceAll(v, " ", "u+0020")
			co[n2] = strings.ReplaceAll(co[n2], "[", "|")
			co[n2] = strings.ReplaceAll(co[n2], "]", "|")
			n = n + 2
			n2 = n2 + 2

		}
	}
	co = append(co, "\n!*!")
	return strings.Join(co, "")
}
func addTableFrontiers(tables []string) string {
	newTables := make([]string, (len(tables)*2)+1)
	newTables[0] = "////"
	i := 1
	for _, t := range tables {
		newTables[i] = t
		newTables[i+1] = "////"
		i = i + 2
	}
	return strings.Join(newTables, "")
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
