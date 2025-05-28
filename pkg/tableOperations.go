package pkg

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"slices"
	"sort"
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
	GetRowById(id string) (Row, error)
	GetRows() Rows
	GetColumns() []string
	PrintTable()
	GetName() string
}

type Rows []Row
type Row struct {
	Value string
}
type table struct {
	name     string
	columns  []string
	values   []Row
	rawTable string
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
	newTableEncode := utilities.Must(globalEncoderKey.Encode(newTable))
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTableEncode), 0666))
	return table, nil
}
func (table table) PrintTable() {
	fmt.Println(table.rawTable)
}
func (table table) GetName() string {
	return table.name
}
func (table table) AddValues(values []string) Table {
	rows := make([]Row, len(values))
	for i, v := range values {
		rows[i] = Row{Value: v}
	}
	s := valuesBuilder(table.rawTable, rows)
	table.rawTable = strings.Replace(table.rawTable, "!*!", s, 1)
	tables := getTables()
	for i, t := range tables {
		if strings.Contains(t.name, table.name) {
			tables[i] = table
			break
		}
	}
	newTable := addTableFrontiers(tables)
	newTableEncode := utilities.Must(globalEncoderKey.Encode(newTable))
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTableEncode), 0666))
	return table
}
func (table table) GetColumns() []string {
	return getColumns(table.rawTable)
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
	newTableEncode := utilities.Must(globalEncoderKey.Encode(newTable))
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTableEncode), 0666))
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
	newTableEncode := utilities.Must(globalEncoderKey.Encode(newTable))
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTableEncode), 0666))
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
			rowSlice := strings.Split(row.Value, "|")
			rowSlice[i+1] = " " + newValue + " "
			row.Value = strings.Join(rowSlice, "|")
			updateTable, err := updateRow(table.rawTable, id, row.Value)
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
	newTableEncode := utilities.Must(globalEncoderKey.Encode(newTable))
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTableEncode), 0666))
	return table, nil
}
func (table table) GetRows() Rows {
	values := getValues(table.rawTable)
	return values
}
func (table table) GetRowById(id string) (Row, error) {
	rows := getValues(table.rawTable)
	for i := 1; i < len(rows); i++ {
		s := strings.Split(rows[i].Value, "|")
		if strings.TrimSpace(s[2]) == id {
			return rows[i], nil
		}
	}
	return Row{}, &NotFoundError{itemName: "Row"}
}
func (table table) DeleteRow(id string) (Table, error) {
	row, err := table.GetRowById(id)
	if err != nil {
		return nil, err
	}
	rowSlice := strings.Split(table.rawTable, "\n")
	index := slices.Index(rowSlice, row.Value)
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
	newTableEncode := utilities.Must(globalEncoderKey.Encode(newTable))
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTableEncode), 0666))
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
	newTableEncode := utilities.Must(globalEncoderKey.Encode(tablesWithFrontier))
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTableEncode), 0666))
	return table, nil
}
func (r Rows) String() string {
	s := make([]string, len(r))
	return strings.Join(s, "\n")
}

func (r Rows) OrderBy(column string) Rows {

	newSlice := make([]Row, len(r)-1)
	for i := 1; i < len(r); i++ {
		newSlice[i-1] = r[i]
	}
	sRow := strings.Split(r[0].Value, " ")

	var columnIndex int
	for i, v := range sRow {
		if v == column {
			columnIndex = i
		}
	}

	sort.Slice(newSlice, func(i, j int) bool {
		s := strings.Split(newSlice[i].Value, " ")
		s2 := strings.Split(newSlice[j].Value, " ")
		return s[columnIndex] > s2[columnIndex]
	})
	newSlice = append(Rows{r[0]}, newSlice...)
	return newSlice
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
func valuesBuilder(table string, values []Row) string {
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
	for _, r := range values {

		if n > count {
			break
		}
		co[n] = strings.ReplaceAll(r.Value, " ", "U+0020")
		n = n + 2

	}
	co = append(co, "\n!*!")
	return strings.Join(co, " ")
}
