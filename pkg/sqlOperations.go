package pkg

import (
	"errors"
	"slices"
	"strings"
	"sync"
	"text-database/pkg/utilities"
)

type SqlRows struct {
	AffectRows int
	Rows       Rows
}

func validateSql(d db, sql string) (SqlRows, error) {
	sql = strings.ReplaceAll(sql, ",", " ")
	sql = strings.ReplaceAll(sql, "(", " ")
	sql = strings.ReplaceAll(sql, ")", " ")
	sqlS := strings.Split(sql, " ")
	sqlS = utilities.RemoveEmptyIndex(sqlS)
	sqlS[0] = strings.ToUpper(sqlS[0])
	switch sqlS[0] {
	case "SELECT":
		upper := strings.ToUpper(sql)
		if !strings.Contains(upper, "FROM") {
			return SqlRows{}, errors.New("invalid sql")
		}
		result := sqlSelect(sqlS)
		return result, nil
	case "CREATE":
		break
	case "UPDATE":
		return sqlUpdate(sqlS)
	case "DELETE":
		return sqlDelete(sqlS)
	case "INSERT":
		if strings.ToUpper(sqlS[1]) != "INTO" {
			return SqlRows{}, errors.New("invalid sql")
		}
		return sqlInsert(sqlS)
	case "DROP":
		err := sqlDrop(d, sqlS)
		return SqlRows{}, err
	default:
	}
	return SqlRows{}, nil
}
func sqlDrop(d db, sqlS []string) error {
	tableName := sqlS[2]
	err := d.DeleteTable(tableName)
	if err != nil {
		return err
	}
	return nil
}
func sqlSelect(sqlS []string) SqlRows {
	index := slices.Index(sqlS, "FROM")
	tableName := sqlS[index+1]
	tb, _ := getTableByName(tableName, true)
	rows := tb.GetRows()
	columns := getSqlColumns(tb, sqlS)
	whereParams := sqlWhere(sqlS)

	var result []string
	for i := 0; i < len(rows); i++ {
		for _, v := range columns {
			value := rows[i].SearchValue(v)
			result = append(result, value)
		}
	}

	sqlValues := valuesBuilderSql(columns, result)
	var finalResult Rows
	if whereParams != nil {
		for _, v := range sqlValues {
			if v.SearchValue(whereParams[0]) == whereParams[2] {
				finalResult = append(finalResult, v)
			}
		}
	} else {
		finalResult = sqlValues
	}
	sqlRows := &SqlRows{
		AffectRows: 0,
		Rows:       finalResult,
	}
	return *sqlRows
}
func sqlWhere(sqlS []string) []string {
	index := slices.Index(sqlS, "WHERE")
	if index == -1 {
		return nil
	}
	params := sqlS[index+1:]
	params = fixSqlParams(params)
	return params
}
func sqlUpdate(sqlS []string) (SqlRows, error) {
	updateIndex := slices.Index(sqlS, "UPDATE")
	tableName := sqlS[updateIndex+1]
	tb, _ := getTableByName(tableName, true)
	setIndex := slices.Index(sqlS, "SET")
	whereParams := sqlWhere(sqlS)
	var whereIndex int
	if whereParams != nil {
		whereIndex = slices.Index(sqlS, "WHERE")
	}
	newS := sqlS[setIndex+1 : whereIndex]
	newS = fixSqlParams(newS)
	rows := tb.SearchAll(whereParams[0], whereParams[2])
	for i := 0; i < len(newS); i += 3 {
		for _, row := range rows {
			err := tb.UpdateValue(newS[i], row.SearchValue("id"), newS[i+2])
			if err != nil {
				return SqlRows{}, err
			}
		}
	}
	tb.save()
	return SqlRows{
		AffectRows: len(rows),
		Rows:       nil,
	}, nil

}
func sqlDelete(sqlS []string) (SqlRows, error) {
	fromIndex := slices.Index(sqlS, "FROM")
	tableName := sqlS[fromIndex+1]
	tb, _ := getTableByName(tableName, true)
	whereParams := sqlWhere(sqlS)
	rows := tb.SearchAll(whereParams[0], whereParams[2])
	for _, row := range rows {
		err := tb.DeleteRow(row.SearchValue("id"), false)
		if err != nil {
			return SqlRows{}, err
		}
	}
	tb.save()
	return SqlRows{
		AffectRows: len(rows),
		Rows:       nil,
	}, nil
}
func sqlInsert(sqlS []string) (SqlRows, error) {
	insertIndex := slices.Index(sqlS, "INSERT")
	tableName := sqlS[insertIndex+2]
	valuesIndex := slices.Index(sqlS, "VALUES")
	tb, _ := getTableByName(tableName, true)
	columns := sqlS[insertIndex+3 : valuesIndex]
	for _, v := range columns {
		if !slices.Contains(tb.columns, v) {
			return SqlRows{}, &NotFoundError{itemName: "Column: " + v + " in table: " + tableName}
		}
	}
	if len(columns) != len(tb.columns)/2 {
		return SqlRows{}, errors.New("column number does not match")
	}

	values := sqlS[valuesIndex+1:]
	a := divideEachNewRow(len(columns), values)
	for _, v := range a {
		tb.addValuesIdGenerationOff(v)
	}
	tb.save()
	d := len(a)
	return SqlRows{
		AffectRows: d,
		Rows:       nil,
	}, nil
}
func divideEachNewRow(columns int, values []string) [][]string {
	if columns == len(values) {
		return [][]string{values}
	}
	a := &[][]string{}
	var b []string
	count := 0
	for _, v := range values {
		if columns != count {
			b = append(b, v)
			count++
		}
		if columns == count {
			*a = append(*a, b)
			b = []string{}
			count = 0
		}
	}
	return *a
}
func fixSqlParams(params []string) []string {
	var wg sync.WaitGroup
	ch := make(chan []string)
	var checkedParams []string
	for _, v := range params {
		wg.Add(1)
		go func(val string) {
			defer wg.Done()
			newS := replaceComparativeSymbol(val)
			var newParams []string
			if newS != nil {
				newParams = append(newParams, newS...)
			} else {
				newParams = append(newParams, val)
			}
			ch <- newParams
		}(v)
		checkedParams = append(checkedParams, <-ch...)
	}
	wg.Wait()
	close(ch)
	return checkedParams
}
func replaceComparativeSymbol(value string) []string {
	symbols := []string{"=", ">", "<", ">=", "<="}
	for _, s := range symbols {
		if strings.Contains(value, s) && value != s {
			newS := strings.Split(value, s)
			index := slices.Index(newS, "")
			if index != -1 {
				newS[index] = s
			} else {
				newS = slices.Insert(newS, 1, s)
			}
			return newS
		}
	}
	return nil
}
func getSqlColumns(tb table, sqlS []string) []string {
	fromIndex := slices.Index(sqlS, "FROM")

	var columns []string
	if sqlS[1] != "*" {
		for i := 1; i < fromIndex; i++ {
			columns = append(columns, sqlS[i])
		}
	} else {
		tbColumns := tb.GetColumns()
		for i := 1; i < len(tbColumns); i += 2 {
			value := tbColumns[i]
			columns = append(columns, value)
		}
	}
	return columns
}
func valuesBuilderSql(sqlColumns []string, sqlValues []string) Rows {
	columns := columnsBuilder(sqlColumns)
	columnsS := strings.Split(columns, " ")
	formattedColumns := utilities.RemoveEmptyIndex(columnsS)
	count := len(columnsS)
	for i := 0; i < count; i += 2 {
		formattedColumns[i] = strings.ReplaceAll(formattedColumns[i], "[", "|")
		formattedColumns[i] = strings.ReplaceAll(formattedColumns[i], "]", "|")
	}
	n := 0
	var result Rows
	count = len(sqlColumns)
	for i := 0; i < len(sqlValues); i++ {
		if count == n {
			n = 0
		}
		index := slices.Index(columnsS, sqlColumns[n])
		formattedColumns[index] = sqlValues[i]
		formattedResult := strings.Join(formattedColumns, " ")
		formattedResult = strings.Trim(formattedResult, " ")
		n++
		if count == n {
			row := &Row{
				columns: columnsS,
				value:   formattedResult,
			}
			result = append(result, *row)
		}
	}
	return result

}
