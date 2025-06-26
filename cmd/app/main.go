package main

import (
	"flag"
	"fmt"
	"text-database/pkg"
)

func main() {
	data := []pkg.DataConfig{
		{
			TableName: "Esto",
			Columns:   []string{"name", "age"},
			Values:    []pkg.Values{{"1", "carlos", "32"}, {"2", "jose", "23"}},
		},

		{
			TableName: "Segunda Tabla",
			Columns:   []string{"Nombre", "Direccion"},
			Values:    []pkg.Values{{"1", "Perez", "Manuguayabo"}},
		},
	}
	config := pkg.DbConfig{EncryptionKey: "", DatabaseName: "database.txt", DataConfig: data}
	_, err := config.CreateDatabase()
	if err != nil {
		fmt.Println(err)
		return
	}
	migrationName := flag.String("ma", "", "Name of the migration")
	flag.Parse()
	if migrationName != nil && *migrationName != "" {
		config.CreateMigration(*migrationName)
	}

}
