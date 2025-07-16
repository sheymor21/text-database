package main

import (
	"flag"
	"fmt"
	"text-database/pkg"
)

func main() {

	data := getExampleDataConfig()
	config := pkg.DbConfig{EncryptionKey: "", DatabaseName: "database.txt", DataConfig: data}
	db, err := config.CreateDatabase()
	if err != nil {
		fmt.Println(err)
		return
	}

	foreignKeys := getExampleForeignKeys()
	errF := db.AddForeignKeys(foreignKeys)
	if errF != nil {
		fmt.Println(errF)
	}
	migrationName := flag.String("ma", "", "Name of the migration")
	flag.Parse()
	if migrationName != nil && *migrationName != "" {
		config.CreateMigration(*migrationName)
	}

}
