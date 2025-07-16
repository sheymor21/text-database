package text_database

import (
	"fmt"
	"github.com/sheymor21/text-database/utilities"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (c DbConfig) CreateMigration(migrationName string) {
	path := utilities.Must(os.Getwd())
	migrationPath := path + "/migrations"
	constructorPath := migrationPath + "/constructor.go"
	if IsFileNameExist(migrationPath, migrationName) {
		panic("Migration already exist")
	}
	fileRoute := fmt.Sprintf("%s/%s_Migration_%d.go", migrationPath, migrationName, time.Now().Unix())
	if !utilities.IsFileExist(fileRoute) {
		code := migrationBuilder(c, migrationName)
		utilities.ErrorHandler(os.MkdirAll(migrationPath, 0755))
		utilities.ErrorHandler(os.WriteFile(fileRoute, code, 0755))
		if !utilities.IsFileExist(constructorPath) {
			constructorCode := constructorBuilder(migrationName)
			utilities.ErrorHandler(os.WriteFile(constructorPath, constructorCode, 0755))
		} else {
			err := os.Remove(constructorPath)
			if err != nil {
				return
			}
			constructorCode := constructorBuilder(migrationName)
			utilities.ErrorHandler(os.WriteFile(constructorPath, constructorCode, 0755))

		}
	}
}

func migrationBuilder(c DbConfig, migrationName string) []byte {
	var builder strings.Builder

	imports := `
package migrations

import "text-database/tdb"
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
	config := tdb.DbConfig{EncryptionKey: "", DatabaseName: "%s"}
	db, err := config.CreateDatabase()
	if err != nil {
		return
	}

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

func constructorBuilder(migrationName string) []byte {
	var builder strings.Builder
	var namesBuilder strings.Builder
	namesBuilder.WriteString("// Migrations Order:\n")
	for i, m := range getMigrationsNames() {
		names := fmt.Sprintf("// [%d] %s", i+1, m)
		namesBuilder.WriteString(names)
		namesBuilder.WriteString("\n")

	}
	imports := `
package migrations

`
	types := `
type table struct {
	name    string
	columns []string
	values  []value
}
type value struct {
	value []string
}`
	lastMigration := fmt.Sprintf(`
func GenerateMigration() {
	generate%s()
}`, upperCase(migrationName))

	builder.WriteString(imports)
	builder.WriteString(namesBuilder.String())
	builder.WriteString(types)
	builder.WriteString(lastMigration)
	return []byte(builder.String())
}

func migrationTableBuilder(c DbConfig) string {
	var builder strings.Builder
	tables := getTables(false)

	for _, t := range tables {
		tableName := strings.ReplaceAll(t.nameRaw, "-", "")
		columns := migrationColumnBuilder(t.columns)
		var values string
		if c.DataConfig != nil {
			values = migrationValuesBuilder(t.nameRaw, c.DataConfig)
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
func IsFileNameExist(path string, fileName string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		name := strings.Split(entry.Name(), "_Migration_")
		if name[0] == fileName {
			return true
		}
	}
	return false
}

func getMigrationsNames() []string {
	path := utilities.Must(os.Getwd())
	migrationPath := path + "/migrations"
	entries, err := os.ReadDir(migrationPath)
	if err != nil {
		panic(err)
	}
	var fileNames []string
	for _, entry := range entries {
		if entry.Name() != "constructor.go" {
			fileNames = append(fileNames, entry.Name())
		}
	}
	sort.Slice(fileNames, func(i, j int) bool {
		numberStrI := strings.Split(fileNames[i], "_Migration_")[1]
		numberStrJ := strings.Split(fileNames[j], "_Migration_")[1]
		numberStrI = strings.ReplaceAll(numberStrI, ".go", "")
		numberStrJ = strings.ReplaceAll(numberStrJ, ".go", "")
		numberI, _ := strconv.Atoi(numberStrI)
		numberJ, _ := strconv.Atoi(numberStrJ)
		return numberI < numberJ
	})
	return fileNames
}
