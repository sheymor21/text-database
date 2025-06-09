package pkg

import (
	"fmt"
	"os"
	"strings"
	"text-database/pkg/utilities"
	"time"
)

func (c DbConfig) CreateMigration() {
	if !utilities.IsFileExist("migrations/initialMigration.go") {
		utilities.ErrorHandler(os.MkdirAll("migrations", 0755))
		code := migrationBuilder(c)
		utilities.ErrorHandler(os.WriteFile("migrations/initialMigration.go", code, 0755))
	}
}

func migrationBuilder(c DbConfig) []byte {
	var builder strings.Builder

	imports := `
package migrations

import (
	"os"
	"text-database/pkg"
)
`
	constants := fmt.Sprintf(`const (
	migrationName        = %q
	databaseName         = %q
	migrationVersion     = 1
	migrationDate        = %q
	migrationDescription = %q
)`, "initialMigration", c.DatabaseName, time.Now().UTC().String(), "Initial Migration")

	types := `
type table struct {
	name    string
	columns []string
	values  []value
}
type value struct {
	value []string
}
`
	functionGenerateTable := fmt.Sprintf(`
func GenerateTable() {
%s
	config := pkg.DbConfig{EncryptionKey: "", DatabaseName: databaseName}
	db, err := config.CreateDatabase()
	if err != nil {
		return
	}
	cleanDatabase()

	for _, t := range *tablesS {
		tb := db.NewTable(t.name, t.columns)
		for _, v := range t.values {
			tb = tb.AddValues(v.value)
		}
	}
}`, migrationTableBuilder(c))

	functionCleanDatabase := fmt.Sprintf(`
func cleanDatabase() {
	file, err := os.ReadFile(databaseName)
	if err != nil {
		return
	}
	file = []byte("")
	err = os.WriteFile(databaseName, file, 0755)
	if err != nil {
		return
	}
}`)
	builder.WriteString(imports)
	builder.WriteString(constants)
	builder.WriteString(types)
	builder.WriteString(functionGenerateTable)
	builder.WriteString(functionCleanDatabase)
	return []byte(builder.String())
}

func migrationTableBuilder(c DbConfig) string {
	var builder strings.Builder
	tables := getTables()

	for _, t := range tables {
		tableName := strings.ReplaceAll(t.name, "-", "")
		columns := migrationColumnBuilder(t.columns)
		var values string
		if c.DataConfig != nil {
			values = migrationValuesBuilder(t.name, c.DataConfig)
		}
		//values := "[]value{}"
		tableStr := fmt.Sprintf(`{
			name:    %q,
			columns: %s,
			values:  %s,
		},`,
			tableName, columns, values)
		builder.WriteString(tableStr)
	}
	tableStrF := fmt.Sprintf(`tablesS := &[]table{
		%s
	}`, builder.String())
	return tableStrF
}
func migrationColumnBuilder(column []string) string {
	columnNames := make([]string, (len(column)/2)-1)
	n := 0
	for i := 3; i < len(column); i = i + 2 {
		columnNames[n] = column[i]
		n++
	}

	slices := stringSliceBuilder(columnNames)
	columnStr := fmt.Sprintf("[]string{%s}", slices)
	return columnStr
}
func migrationValuesBuilder(tableName string, r []DataConfig) string {
	var builder strings.Builder
	for _, v := range r {
		if v.TableName == strings.ReplaceAll(tableName, "-", "") {
			for _, iv := range v.Values {
				valuesS := stringSliceBuilder(iv)
				a := fmt.Sprintf("{value: []string{%s}}", valuesS)
				builder.WriteString(a)
				builder.WriteString(",")
			}
			break
		}
	}
	dataName := fmt.Sprintf("[]value{%s}", builder.String())
	return dataName
}
func stringSliceBuilder(s []string) string {
	var builder strings.Builder
	for i := 0; i < len(s); i++ {

		builder.WriteString(fmt.Sprintf(`"%s"`, s[i]))
		if i != len(s)-1 {
			builder.WriteString(",")
		}
	}
	return builder.String()
}
