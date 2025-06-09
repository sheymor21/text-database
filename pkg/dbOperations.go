package pkg

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"text-database/pkg/utilities"
)

type Db interface {
	GetName() string
	GetTables() []Table
	GetTableByName(name string) (Table, error)
	PrintTables()
	NewTable(name string, columns []string) Table
	DeleteTable(tableName string)
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
		if c.DataConfig == nil {
			setDefaultData(c)
		} else {
			utilities.ErrorHandler(os.WriteFile(c.DatabaseName, []byte{}, 0644))
			newDb := db{name: c.DatabaseName, tables: getTables()}
			addData(newDb, c.DataConfig)
			return newDb, nil
		}
	}
	return db{name: c.DatabaseName, tables: getTables()}, nil

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
	t := &table{name, columns, nil, ""}
	tb := d.addTable(*t)
	return tb
}
func (d db) GetTableByName(name string) (Table, error) {
	return getTableByName(name)
}
func (d db) addTable(table table) Table {
	//data := utilities.Must(os.ReadFile(dbName))
	data := globalEncoderKey.ReadAndDecode(dbName)
	dataByte := []byte(data)
	raw := tableBuilder(table)
	dataByte = append(dataByte, []byte(raw)...)
	dataEncode := utilities.Must(globalEncoderKey.Encode(string(dataByte)))

	utilities.ErrorHandler(os.WriteFile(dbName, []byte(dataEncode), 0666))
	return utilities.Must(d.GetTableByName(table.name))

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
	newTableEncode := utilities.Must(globalEncoderKey.Encode(newTable))
	utilities.ErrorHandler(os.WriteFile(dbName, []byte(newTableEncode), 0666))
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
	data := globalEncoderKey.ReadAndDecode(dbName)
	data = strings.ReplaceAll(data, "\r", "")
	data = strings.ReplaceAll(data, "U+0020", " ")
	s := strings.Split(data, "////")
	sif := utilities.RemoveEmptyIndex(s)
	tables := make([]table, len(sif))
	for i, t := range sif {
		name := getTableName(t)
		values := getValues(t)
		tables[i] = table{name, getColumns(t), values, t}
	}
	return tables
}
func tableBuilder(table table) string {
	columnsRaw := columnsBuilder(table.columns)
	tableRaw := fmt.Sprintf("\n-----%s-----\n"+
		"[1] id %s"+
		"\n!*!"+
		"\n-----%s_End-----\n////", table.name, columnsRaw, table.name)
	if table.values != nil {
		values := valuesBuilder(tableRaw, table.values)
		tableRaw = strings.Replace(tableRaw, "!*!", values, 1)
	}
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
func getValues(table string) []Row {
	row := strings.Split(table, "\n")
	newRow := make([]Row, len(row)-5)
	n := 0
	for i := 2; i < len(row)-3; i++ {
		newRow[n].Value = row[i]
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
	for _, v := range d {
		for _, iv := range v.Values {

			if len(v.Columns) != len(iv) {
				message := fmt.Sprintf("columns and values must have the same length, %d != %d, table: %s", len(v.Columns), len(iv), v.TableName)
				return errors.New(message)
			}
		}
		if v.TableName == "" {
			return errors.New("table name is required")
		}
		if len(v.Columns) == 0 {
			return errors.New("columns are required")
		}
	}
	return nil
}
func addData(db db, d []DataConfig) {
	for _, v := range d {
		tb := db.NewTable(v.TableName, v.Columns)
		for _, iv := range v.Values {
			tb = tb.AddValues(iv)
		}
	}
}
func setDefaultData(c DbConfig) {
	initData, err := os.ReadFile("internal/layout.txt")
	// if you use the test command will trigger this route
	if err != nil {
		initData = utilities.Must(os.ReadFile("../internal/layout.txt"))
	}
	if strings.TrimSpace(c.EncryptionKey) != "" {
		encodeData := utilities.Must(globalEncoderKey.Encode(string(initData)))
		utilities.ErrorHandler(os.WriteFile(c.DatabaseName, []byte(encodeData), 0644))
	} else {
		utilities.ErrorHandler(os.WriteFile(c.DatabaseName, initData, 0644))
		encryptionKeyExist = false
	}
}
