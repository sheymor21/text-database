package text_database

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sheymor21/text-database/utilities"
	"os"
	"slices"
	"sort"
	"strings"
)

type Table interface {
	AddValue(column string, value string) error
	AddValues(values ...string)
	addValuesIdGenerationOff(values []string)
	UpdateTableName(newName string)
	UpdateColumnName(oldColumnName string, newColumnName string) error
	UpdateValue(columnName string, id string, newValue string) error
	DeleteRow(id string, cascade bool) error
	DeleteColumn(columnName string) error
	GetRowById(id string) (Row, error)
	GetRows() Rows
	GetColumns() []string
	PrintTable()
	GetName() string
	SearchOne(column string, value string) (Row, error)
	SearchAll(column string, value string) Rows
	SearchByForeignKey(id string) ([]ComplexRow, error)
	getSimpleName() string
	table() table
	save()
}
type Rows []Row
type Row struct {
	columns []string
	value   string
}

type ComplexRow struct {
	Table Table
	Rows  Rows
}
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
func (t *table) AddValue(column string, value string) error {
	tables := getTables(false)
	s, err := valueBuilder(*t, column, value)
	if err != nil {
		return err
	}

	t.rawTable = strings.Replace(t.rawTable, "!*!", s, 1)
	for i, t := range tables {
		if strings.Contains(t.nameRaw, t.nameRaw) {
			tables[i] = t
			break
		}
	}
	saveTables(tables)
	return nil
}
func (t *table) PrintTable() {
	fmt.Println(t.rawTable)
}
func (t *table) GetName() string {
	return t.nameRaw
}
func (t *table) AddValues(values ...string) {
	*t = addValues(*t, values, true)
}
func (t *table) addValuesIdGenerationOff(values []string) {
	*t = addValues(*t, values, false)
}
func (t *table) GetColumns() []string {
	return getColumns(t.rawTable)
}
func (t *table) UpdateTableName(newName string) {
	formatName := strings.Replace(t.nameRaw, "-----", "", 2)
	formatName = formatName + "_End"
	formatName = fmt.Sprintf("-----%s-----", formatName)
	rawNewName := fmt.Sprintf("-----%s-----", newName)
	rawNewNameEnd := fmt.Sprintf("-----%s-----", newName+"_End")

	t.nameRaw = rawNewName
	t.rawTable = strings.Replace(t.rawTable, t.nameRaw, rawNewName, 1)
	t.rawTable = strings.Replace(t.rawTable, formatName, rawNewNameEnd, 1)
	tables := getTables(false)
	for i, t := range tables {
		if t.nameRaw == t.nameRaw {
			tables[i].nameRaw = t.nameRaw
			tables[i] = t
			break
		}
	}
	saveTables(tables)
}
func (t *table) UpdateColumnName(oldColumnName string, newColumnName string) error {
	if !slices.Contains(t.columns, oldColumnName) {
		return &NotFoundError{itemName: "Column"}
	}
	tables := getTables(false)
	for i, c := range t.columns {
		if c == oldColumnName {
			t.columns[i] = newColumnName
			break
		}
	}

	t.rawTable = strings.Replace(t.rawTable, oldColumnName, newColumnName, 1)
	for i, tb := range tables {
		if tb.nameRaw == tb.nameRaw {
			tables[i] = tb
			break
		}
	}

	saveTables(tables)
	return nil

}
func (t *table) UpdateValue(columnName string, id string, newValue string) error {

	if !slices.Contains(t.columns, columnName) {
		return &NotFoundError{itemName: "Column"}
	}
	for i := 3; i < len(t.columns); i += 2 {
		if strings.TrimSpace(t.columns[i]) == columnName {
			row, rowErr := t.GetRowById(id)
			if rowErr != nil {
				return rowErr
			}
			rowSlice := strings.Split(row.value, "|")
			rowSlice[i+1] = " " + newValue + " "
			row.value = strings.Join(rowSlice, "|")
			row.value = strings.Trim(row.value, " ")
			updateTable, err := updateRow(t.rawTable, id, row.value)
			if err != nil {
				return err
			}
			t.rawTable = updateTable
			break
		}
	}

	tables := getTables(false)
	for i, tb := range tables {
		if tb.nameRaw == tb.nameRaw {
			tables[i].rawTable = tb.rawTable
			break
		}
	}
	saveTables(tables)
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
		if strings.TrimSpace(s[1]) == id {
			return rows[i], nil
		}
	}
	return Row{}, &NotFoundError{itemName: "Row"}
}
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
	if !slices.Contains(t.columns, columnName) {
		return &NotFoundError{itemName: "Column"}
	}
	for n, c := range t.columns {
		if c == columnName {
			t.columns = slices.Delete(t.columns, n, n+1)
			position := fmt.Sprintf("[%d]", n-1)
			t.columns = slices.Replace(t.columns, n-1, n, position)
			if n < len(t.columns) {
				t.columns = slices.Delete(t.columns, n, n+1)
			} else {
				t.columns = slices.Delete(t.columns, n-1, n)
			}
			rawColumn := strings.Join(t.columns, " ")
			newTable := strings.Split(t.rawTable, "\n")
			newTable[2] = rawColumn
			t.rawTable = strings.Join(newTable, "\n")
			t.rawTable = deleteColumnData(t.rawTable, n)
			break
		}
	}
	tables := getTables(false)
	for i, tb := range tables {
		if tb.nameRaw == tb.nameRaw {
			tables[i] = tb
			break
		}
	}
	saveTables(tables)
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
func (r *Rows) OrderByAscend(column string) error {
	order, err := orderBy(*r, column, true)
	if err != nil {
		return err
	}
	*r = order
	return err
}
func (r *Rows) OrderByDescend(column string) error {
	order, err := orderBy(*r, column, false)
	if err != nil {
		return err
	}
	*r = order
	return nil
}
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

func deleteRow(tb table, id string) (table, error) {
	row, err := tb.GetRowById(id)
	if err != nil {
		return table{}, err
	}
	rowSlice := strings.Split(tb.rawTable, "\n")
	index := slices.Index(rowSlice, row.value)
	newRow := slices.Replace(rowSlice, index, index+1, "")
	newRow = utilities.RemoveEmptyIndex(newRow)
	rowString := strings.Join(newRow, "\n")
	tb.rawTable = "\n" + rowString + "\n"

	tables := getTables(false)
	for i, t := range tables {
		if t.nameRaw == tb.nameRaw {
			tables[i] = tb
			break
		}
	}

	saveTables(tables)
	return tb, nil

}
func removeStrConv(r Rows) Rows {
	for i := 0; i < len(r); i++ {
		value := strings.ReplaceAll(r[i].value, "U+0020", " ")
		r[i].value = value
	}
	return r
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
	var newTable string
	if len(tables) != 0 {
		newTable = addTableFrontiers(tables)
	}
	if encryptionKeyExist {
		newTable = utilities.Must(globalEncoderKey.Encode(newTable))
	}
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTable), 0666))
}
