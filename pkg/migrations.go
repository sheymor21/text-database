package pkg

import (
	"fmt"
	"os"
	"strings"
	"text-database/pkg/utilities"
	"time"
)

func (c DbConfig) CreateMigration(migrationName string) {
	fileRoute := fmt.Sprintf("../../migrations/%s_%d.go", migrationName, time.Now().Unix())
	if !utilities.IsFileExist(fileRoute) {
		code := migrationBuilder(c, migrationName)
		utilities.ErrorHandler(os.MkdirAll("../../migrations", 0755))
		utilities.ErrorHandler(os.WriteFile(fileRoute, code, 0755))
		if !utilities.IsFileExist("../../migrations/constructor.go") {
			constructorCode := constructorBuilder(c.DatabaseName)
			utilities.ErrorHandler(os.WriteFile("../../migrations/constructor.go", constructorCode, 0755))
		}
	}
}

func migrationBuilder(c DbConfig, migrationName string) []byte {
	var builder strings.Builder

	imports := `
package migrations

import "text-database/pkg"
`
	constants := fmt.Sprintf(`
	// migrationName        = %q
	// databaseName         = %q
	// migrationVersion     = 1
	// migrationDate        = %q
	// migrationDescription = %q
`, migrationName, c.DatabaseName, time.Now().UTC().String(), "Initial Migration")

	functionGenerateTable := fmt.Sprintf(`
func generate%s() {
%s
	config := pkg.DbConfig{EncryptionKey: "", DatabaseName: "%s"}
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
}`, upperCase(migrationName), migrationTableBuilder(c), c.DatabaseName)

	builder.WriteString(imports)
	builder.WriteString(constants)
	builder.WriteString(functionGenerateTable)
	return []byte(builder.String())
}

func constructorBuilder(databaseName string) []byte {
	var builder strings.Builder
	imports := `
package migrations

import "os"
`
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
	functionCleanDatabase := fmt.Sprintf(`
func cleanDatabase() {
databaseName := "%s"
	file, err := os.ReadFile(databaseName)
	if err != nil {
		return
	}
	file = []byte("")
	err = os.WriteFile(databaseName, file, 0755)
	if err != nil {
		return
	}
}`, databaseName)
	builder.WriteString(imports)
	builder.WriteString(types)
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

func upperCase(s string) string {
	sS := strings.Split(s, "")
	sS[0] = strings.ToUpper(sS[0])
	return strings.Join(sS, "")

}
