package tdb

import (
	"errors"
	"fmt"
	"github.com/sheymor21/text-database/tdb/utilities"
	"os"
	"slices"
	"strings"
	"sync"
)

type Db interface {
	GetName() string
	GetTables() []Table
	GetTableByName(name string) (Table, error)
	PrintTables()
	NewTable(name string, columns []string) Table
	DeleteTable(tableName string) error
	AddForeignKey(key ForeignKey) error
	AddForeignKeys(keys []ForeignKey) error
	FromSql(sql string) (SqlRows, error)
}
type db struct {
	name   string
	tables []table
}
type DataConfig struct {
	TableName string
	Columns   []string
	Values    []Values
}
type Values []string
type DbConfig struct {
	EncryptionKey string
	DatabaseName  string
	DataConfig    []DataConfig
}
type ForeignKey struct {
	TableName         string
	ColumnName        string
	ForeignTableName  string
	ForeignColumnName string
}

var encryptionKeyExist bool
var dbName string

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
		globalEncoderKey = *NewSecureTextEncoder(c.EncryptionKey)
		encryptionKeyExist = true
	}
	if !utilities.IsFileExist(c.DatabaseName) {
		utilities.ErrorHandler(os.WriteFile(c.DatabaseName, []byte{}, 0644))
		if c.DataConfig == nil {
			setDefaultData(c)
		}

	} else {
		data := string(utilities.Must(os.ReadFile(c.DatabaseName)))
		if !IsEncode(data) && encryptionKeyExist {
			EncodeAndSave(data)
		}
	}

	if c.DataConfig != nil {
		newDb := setDatabaseData(c)
		return &newDb, nil
	}

	return &db{name: c.DatabaseName, tables: getTables(true)}, nil
}
func (c DbConfig) RemoveEncryption() error {
	if dbName == "" {
		return &NotFoundError{itemName: "Database"}
	}
	if strings.TrimSpace(c.EncryptionKey) != "" {
		data := string(utilities.Must(os.ReadFile(c.DatabaseName)))
		if IsEncode(data) {
			DecodeAndSave(data)
			return nil
		}
	}
	return &NotFoundError{itemName: "EncryptionKey"}
}
func (d *db) GetName() string {
	return dbName
}
func (d *db) GetTables() []Table {
	tables := getTables(true)
	iTables := make([]Table, len(tables))
	for i, t := range tables {
		iTables[i] = &t
	}
	return iTables
}
func (d *db) PrintTables() {
	tables := getTables(true)
	for _, t := range tables {
		fmt.Println(t.rawTable)
	}
}
func (d *db) NewTable(name string, columns []string) Table {
	t := &table{name, columns, nil, ""}
	tb := d.addTable(*t)
	return tb
}
func (d *db) GetTableByName(name string) (Table, error) {
	tb, err := getTableByName(name, true)
	return &tb, err
}
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
		data := string(utilities.Must(os.ReadFile(dbName)))
		linkAdded := string(linkTableLayout()) + data
		utilities.ErrorHandler(os.WriteFile(dbName, []byte(linkAdded), 0666))
	}

	linkTb, _ := getTableByName("Links", false)
	err := validateForeignKey(linkTb, key)
	if err != nil {
		return err
	}
	linkTb.AddValues(key.TableName, key.ColumnName, key.ForeignTableName, key.ForeignColumnName)
	return nil
}
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
	data := globalEncoderKey.ReadAndDecode(dbName)
	dataByte := []byte(data)
	raw := tableBuilder(table)
	dataByte = append(dataByte, []byte(raw)...)
	if encryptionKeyExist {
		dataEncode := utilities.Must(globalEncoderKey.Encode(string(dataByte)))
		utilities.ErrorHandler(os.WriteFile(dbName, []byte(dataEncode), 0666))
		return utilities.Must(d.GetTableByName(table.nameRaw))
	}
	utilities.ErrorHandler(os.WriteFile(dbName, dataByte, 0666))
	return utilities.Must(d.GetTableByName(table.nameRaw))

}
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
func (d *db) FromSql(sql string) (SqlRows, error) {
	return validateSql(*d, sql)
}
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
func getTables(strConv bool) []table {
	data := globalEncoderKey.ReadAndDecode(dbName)
	data = strings.ReplaceAll(data, "\r", "")
	if strConv {
		data = strings.ReplaceAll(data, "U+0020", " ")
	}
	s := strings.Split(data, "////")
	sif := utilities.RemoveEmptyIndex(s)
	tables := make([]table, len(sif))
	for i, t := range sif {
		name := getTableName(t)
		values := getRows(t)
		tables[i] = table{name, getColumns(t), values, t}
	}
	return tables
}
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
func addStaticData(db db, v DataConfig) {
	tb, _ := db.GetTableByName(v.TableName)
	for _, iv := range v.Values {
		if !areValuesInDatabase(v.TableName, iv[0]) {
			tb.addValuesIdGenerationOff(iv)
		}
	}
}
func generateStaticData(db db, v DataConfig) {
	tb := db.NewTable(v.TableName, v.Columns)
	if v.Values != nil || len(v.Values) != 0 {
		for _, iv := range v.Values {
			tb.addValuesIdGenerationOff(iv)
		}
	}
}
func setDefaultData(c DbConfig) {
	if encryptionKeyExist {
		EncodeAndSave(string(getLayout()))
	} else {
		utilities.ErrorHandler(os.WriteFile(c.DatabaseName, getLayout(), 0644))
	}
}
func isTableInDatabase(tableName string) bool {
	_, err := getTableByName(tableName, false)
	if err != nil {
		return false
	}
	return true
}
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
func validateForeignKey(linkTb table, key ForeignKey) error {
	linkRows := getRows(linkTb.rawTable)
	for i := 1; i < len(linkRows); i++ {
		s := strings.Split(linkRows[i].value, "|")
		s = utilities.RemoveEmptyIndex(s)
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
