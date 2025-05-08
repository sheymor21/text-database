package main

import (
	"text-database/pkg"
)

func main() {
	config := pkg.DbConfig{SecurityKey: "", DatabaseName: "database.txt"}
	config.CreateDatabase()

}
