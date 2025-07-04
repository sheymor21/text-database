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
	AddValues(values ...string) Table
	addValuesIdGenerationOff(values []string) Table
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
	SearchOne(column string, value string) (Row, error)
	SearchAll(column string, value string) Rows
}

type Rows []Row
type Row struct {
	columns []string
	value   string
}

type ComplexRow struct {
	Table table
	Rows  Rows
}
type table struct {
	nameRaw  string
	columns  []string
	values   []Row
	rawTable string
}

func (table table) AddValue(column string, value string) (Table, error) {
	tables := getTables(false)
	s, err := valueBuilder(table, column, value)
	if err != nil {
		return nil, err
	}

	table.rawTable = strings.Replace(table.rawTable, "!*!", s, 1)
	for i, t := range tables {
		if strings.Contains(t.nameRaw, table.nameRaw) {
			tables[i] = table
			break
		}
	}
	saveTables(tables)
	return table, nil
}
func (table table) PrintTable() {
	fmt.Println(table.rawTable)
}
func (table table) GetName() string {
	return table.nameRaw
}
func (table table) AddValues(values ...string) Table {
	return addValues(table, values, true)
}
func (table table) addValuesIdGenerationOff(values []string) Table {
	return addValues(table, values, false)
}
func (table table) GetColumns() []string {
	return getColumns(table.rawTable)
}
func (table table) UpdateTableName(newName string) Table {
	formatName := strings.Replace(table.nameRaw, "-----", "", 2)
	formatName = formatName + "_End"
	formatName = fmt.Sprintf("-----%s-----", formatName)
	rawNewName := fmt.Sprintf("-----%s-----", newName)
	rawNewNameEnd := fmt.Sprintf("-----%s-----", newName+"_End")

	table.nameRaw = rawNewName
	table.rawTable = strings.Replace(table.rawTable, table.nameRaw, rawNewName, 1)
	table.rawTable = strings.Replace(table.rawTable, formatName, rawNewNameEnd, 1)
	tables := getTables(false)
	for i, t := range tables {
		if t.nameRaw == table.nameRaw {
			tables[i].nameRaw = table.nameRaw
			tables[i] = table
			break
		}
	}
	saveTables(tables)
	return table
}
func (table table) UpdateColumnName(oldColumnName string, newColumnName string) (Table, error) {
	if !slices.Contains(table.columns, oldColumnName) {
		return nil, &NotFoundError{itemName: "Column"}
	}
	tables := getTables(false)
	for i, c := range table.columns {
		if c == oldColumnName {
			table.columns[i] = newColumnName
			break
		}
	}

	table.rawTable = strings.Replace(table.rawTable, oldColumnName, newColumnName, 1)
	for i, t := range tables {
		if t.nameRaw == table.nameRaw {
			tables[i] = table
			break
		}
	}

	saveTables(tables)
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
			rowSlice := strings.Split(row.value, "|")
			rowSlice[i+1] = " " + newValue + " "
			row.value = strings.Join(rowSlice, "|")
			updateTable, err := updateRow(table.rawTable, id, row.value)
			if err != nil {
				return nil, err
			}
			table.rawTable = updateTable
			break
		}
	}

	tables := getTables(false)
	for i, t := range tables {
		if t.nameRaw == table.nameRaw {
			tables[i].rawTable = table.rawTable
			break
		}
	}
	saveTables(tables)
	return table, nil
}
func (table table) GetRows() Rows {
	values := getRows(table.rawTable)
	return values
}
func (table table) GetRowById(id string) (Row, error) {
	rows := getRows(table.rawTable)
	for i, row := range rows {
		s := strings.Split(row.value, " ")
		if strings.TrimSpace(s[1]) == id {
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
	index := slices.Index(rowSlice, row.value)
	newRow := slices.Replace(rowSlice, index, index+1, "")
	newRow = utilities.RemoveEmptyIndex(newRow)
	rowString := strings.Join(newRow, "\n")
	table.rawTable = "\n" + rowString + "\n"

	tables := getTables(false)
	for i, t := range tables {
		if t.nameRaw == table.nameRaw {
			tables[i] = table
			break
		}
	}

	saveTables(tables)
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
	tables := getTables(false)
	for i, t := range tables {
		if t.nameRaw == table.nameRaw {
			tables[i] = table
			break
		}
	}
	saveTables(tables)
	return table, nil
}
func (table table) SearchOne(column string, value string) (Row, error) {
	rows := table.GetRows()
	index := slices.Index(rows[0].columns, column)
	for _, r := range rows {
		row := strings.Split(r.value, " ")
		if row[index] == value {
			return r, nil
		}
	}
	return Row{}, &NotFoundError{itemName: value}
}
func (table table) SearchAll(column string, value string) Rows {
	return searchAll(table, column, value)
}
func (r Rows) String() string {
	s := make([]string, len(r))
	return strings.Join(s, "\n")
}
func (r Rows) OrderByAscend(column string) (Rows, error) {
	return orderBy(r, column, true)
}
func (r Rows) OrderByDescend(column string) (Rows, error) {
	return orderBy(r, column, false)
}
func (r Row) SearchValue(column string) string {
	index := slices.Index(r.columns, column)
	if index == -1 {
		return ""
	}
	s := strings.Split(r.value, " ")
	return s[index]
}
func (r Row) String() string {
	return r.value
}
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
func valueBuilder(table table, columnName string, value string) (string, error) {
	co := getColumns(table.rawTable)
	co = utilities.RemoveEmptyIndex(co)
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
func valuesBuilder(table string, values []Row, idGenerate bool) string {
	co := getColumns(table)
	co = utilities.RemoveEmptyIndex(co)
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
	tables := getTables(false)
	for i, t := range tables {
		if strings.Contains(t.nameRaw, table.nameRaw) {
			tables[i] = table
			break
		}
	}
	saveTables(tables)
	return table
}
func saveTables(tables []table) {
	newTable := addTableFrontiers(tables)
	if encryptionKeyExist {
		newTable = utilities.Must(globalEncoderKey.Encode(newTable))
	}
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTable), 0666))
}
