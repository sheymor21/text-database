package main

import "text-database/pkg"

func getExampleDataConfig() []pkg.DataConfig {
	data := []pkg.DataConfig{
		{
			TableName: "Users",
			Columns:   []string{"name", "age"},
			Values:    []pkg.Values{{"1", "carlos", "32"}, {"2", "Emilio", "23"}, {"3", "Manuel", "76"}},
		},
		{
			TableName: "Clients",
			Columns:   []string{"Name", "LastName", "Email"},
			Values: []pkg.Values{
				{"1", "Jose", "Coronel", "ejemplo@outlook.com"},
				{"2", "Fabricio", "Carrasco", "Example@gmail.com"},
				{"3", "Emanuel", "Ramirez", "Sipo@hotmail.com"},
			},
		},
		{
			TableName: "Products",
			Columns:   []string{"Name", "Price", "Coin"},
			Values: []pkg.Values{
				{"1", "Soup", "23", "Dollar"},
				{"2", "Ketchup", "19", "Dollar"},
				{"3", "Salami", "30", "Euro"},
				{"4", "Bed", "120", "Dollar"},
			},
		},
		{
			TableName: "Invoice",
			Columns:   []string{"Product_Id", "Client_Id"},
			Values:    []pkg.Values{},
		},
	}
	return data
}

func getExampleForeignKeys() []pkg.ForeignKey {
	foreignKey := &[]pkg.ForeignKey{
		{
			TableName:         "Clients",
			ColumnName:        "id",
			ForeignTableName:  "Invoice",
			ForeignColumnName: "Client_Id",
		},
		{
			TableName:         "Products",
			ColumnName:        "id",
			ForeignTableName:  "Invoice",
			ForeignColumnName: "Product_Id",
		},
	}
	return *foreignKey
}
func s() {

}
