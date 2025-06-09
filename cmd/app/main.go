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
			Values:    []pkg.Values{{"carlos", "32"}, {"jose", "23"}},
		},

		{
			TableName: "Segunda Tabla",
			Columns:   []string{"Nombre", "Direccion"},
			Values:    []pkg.Values{{"Perez", "Manuguayabo"}},
		},
	}
	config := pkg.DbConfig{EncryptionKey: "", DatabaseName: "database.txt", DataConfig: data}
	_, err := config.CreateDatabase()
	migrationName := flag.String("ma", "", "Name of the migration")
	flag.Parse()
	if migrationName != nil {
		config.CreateMigration(*migrationName)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

}
