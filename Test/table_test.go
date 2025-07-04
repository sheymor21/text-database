package Test

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
	"text-database/pkg"
	"text-database/pkg/utilities"
)

func (s *tableSuite) TestGetRowById() {
	tb := utilities.Must(s.db.GetTableByName("Users"))
	row, _ := tb.GetRowById("2")
	result := strings.TrimSpace(row.String())
	if result != "|1| 2 |2| juan |3| 54" {
		s.Fail("Expected |1| 2 |2| juan |3| 54", fmt.Sprintf("Recibe: %s", row))
	}
}
func (s *tableSuite) TestGetRows() {
	tb := utilities.Must(s.db.GetTableByName("Users"))
	count := len(tb.GetRows())
	if count != 4 {
		s.Fail("Expected 4 Rows", fmt.Sprintf("Recibe: %d", count))
	}
}
func (s *tableSuite) TestAddValue() {
	tb := utilities.Must(s.db.GetTableByName("Users"))
	tb, _ = tb.AddValue("name", "test")
	rows := tb.GetRows()
	id, index := getIdAndIndex(rows)
	if rows[index].String() != fmt.Sprintf("|1| %s |2| test |3| null", id) {
		s.Fail(fmt.Sprintf("Expected |1| %s |2| test |3| null", id), fmt.Sprintf("Recibe: %s", rows[2]))
	}

}
func (s *tableSuite) TestAddValues() {
	tb := utilities.Must(s.db.GetTableByName("Users"))
	tb = tb.AddValues("test", "20")
	rows := tb.GetRows()
	id, index := getIdAndIndex(rows)
	if rows[index].String() != fmt.Sprintf("|1| %s |2| test |3| 20", id) {
		s.Fail(fmt.Sprintf("Expected |1| %s |2| test |3| 20", id), fmt.Sprintf("Recibe: %s", rows[2]))
	}
}
func (s *tableSuite) TestUpdateTableName() {
	tb := utilities.Must(s.db.GetTableByName("Users"))
	tb = tb.UpdateTableName("Test")
	if tb.GetName() != "-----Test-----" {
		s.Fail("Expected Test Table", fmt.Sprintf("Recibe: %s", tb))
	}
}
func (s *tableSuite) TestUpdateColumnName() {
	tb := utilities.Must(s.db.GetTableByName("Users"))
	tb, _ = tb.UpdateColumnName("name", "username")
	columns := strings.TrimSpace(strings.Join(tb.GetColumns(), " "))
	if columns != "[1] id [2] username [3] age" {
		s.Fail("Expected [1] id [2] username |3| age ", fmt.Sprintf("Recibe: %s", columns))
	}
}
func (s *tableSuite) TestUpdateValue() {
	tb := utilities.Must(s.db.GetTableByName("Users"))
	tb, _ = tb.UpdateValue("age", "2", "30")
	rows := tb.GetRows()
	if rows[1].String() != "|1| 2 |2| juan |3| 30 " {
		s.Fail("Expected |1| 2 |2| juan |3| 30", fmt.Sprintf("Recibe: %s", rows[2]))
	}
}
func (s *tableSuite) TestDeleteRow() {
	tb := utilities.Must(s.db.GetTableByName("Users"))
	tb, _ = tb.DeleteRow("1")
	rows := tb.GetRows()
	if len(rows) != 3 {
		s.Fail("Expected len of 3", fmt.Sprintf("Recibe: %d", len(rows)))
	}

	if rows[0].String() != "|1| 2 |2| juan |3| 54" {
		s.Fail("Expected |1| 1 |2| juan |3| 54", fmt.Sprintf("Recibe: %s", rows[1]))
	}
}
func (s *tableSuite) TestDeleteColumn() {
	tb := utilities.Must(s.db.GetTableByName("Users"))
	tb, _ = tb.DeleteColumn("age")
	columns := strings.TrimSpace(strings.Join(tb.GetColumns(), " "))
	if columns != "[1] id [2] name" {
		s.Fail("Expected [1] id [2] name", fmt.Sprintf("Recibe: %s", columns))
	}
}
func (s *tableSuite) TestAddValue_ReturnColumnError() {
	tb, _ := s.db.GetTableByName("Users")
	_, err := tb.AddValue("test", "value")
	var example *pkg.NotFoundError
	if !errors.As(err, &example) {
		s.ErrFail(err)
	}
}
func (s *tableSuite) TestUpdateColumnName_ReturnColumnError() {
	tb, _ := s.db.GetTableByName("Users")
	_, err := tb.UpdateColumnName("test", "value")
	var example *pkg.NotFoundError
	if !errors.As(err, &example) {
		s.ErrFail(err)
	}
}
func (s *tableSuite) TestUpdateValue_ReturnColumnError() {
	tb, _ := s.db.GetTableByName("Users")
	_, err := tb.UpdateValue("test", "value", "value")
	var example *pkg.NotFoundError
	if !errors.As(err, &example) {
		s.ErrFail(err)
	}
}
func (s *tableSuite) TestUpdateValue_ReturnIdError() {
	tb, _ := s.db.GetTableByName("Users")
	_, err := tb.UpdateValue("name", "test", "value")
	var example *pkg.NotFoundError
	if !errors.As(err, &example) {
		s.ErrFail(err)
	}
}
func (s *tableSuite) TestDeleteRow_ReturnIdError() {
	tb, _ := s.db.GetTableByName("Users")
	_, err := tb.DeleteRow("test")
	var example *pkg.NotFoundError
	if !errors.As(err, &example) {
		s.ErrFail(err)
	}
}
func (s *tableSuite) TestDeleteColumn_ReturnIdError() {
	tb, _ := s.db.GetTableByName("Users")
	_, err := tb.DeleteColumn("test")
	var example *pkg.NotFoundError
	if !errors.As(err, &example) {
		s.ErrFail(err)
	}
}
func (s *tableSuite) TestGetRowById_ReturnIdError() {
	tb, _ := s.db.GetTableByName("Users")
	_, err := tb.GetRowById("test")
	var example *pkg.NotFoundError
	if !errors.As(err, &example) {
		s.ErrFail(err)
	}
}
func (s *tableSuite) TestOrderByAscend_Numbers() {
	tb, _ := s.db.GetTableByName("Users")
	rows := tb.GetRows()
	newRow, _ := rows.OrderByAscend("age")
	if newRow[1].String() != "|1| 2 |2| juan |3| 54" {
		s.Fail("Expected |1| 2 |2| juan |3| 54", fmt.Sprintf("Recibe: %s", newRow[1]))
	}
}
func (s *tableSuite) TestOrderByDescend_Numbers() {
	tb, _ := s.db.GetTableByName("Users")
	rows := tb.GetRows()
	newRow, _ := rows.OrderByDescend("age")
	if newRow[0].String() != "|1| 1 |2| pedro |3| 32" {
		s.Fail("Expected |1| 1 |2| pedro |3| 32", fmt.Sprintf("Recibe: %s", newRow[1]))
	}
}
func (s *tableSuite) TestOrderByAscend_Letters() {
	tb, _ := s.db.GetTableByName("Users")
	rows := tb.GetRows()
	newRow, _ := rows.OrderByAscend("name")
	if newRow[0].String() != "|1| 1 |2| pedro |3| 32" {
		s.Fail("Expected |1| 1 |2| pedro |3| 32", fmt.Sprintf("Recibe: %s", newRow[1]))
	}
}
func (s *tableSuite) TestOrderByDescend_Letters() {
	tb, _ := s.db.GetTableByName("Users")
	rows := tb.GetRows()
	newRow, _ := rows.OrderByDescend("name")
	if newRow[1].String() != "|1| 2 |2| juan |3| 54" {
		s.Fail("Expected |1| 2 |2| juan |3| 54", fmt.Sprintf("Recibe: %s", newRow[1]))
	}
}
func (s *tableSuite) TestOrderBy_ReturnColumnError() {
	tb, _ := s.db.GetTableByName("Users")
	rows := tb.GetRows()
	_, err := rows.OrderByAscend("Email")
	var example *pkg.NotFoundError
	if !errors.As(err, &example) {
		s.ErrFail(err)
	}
}
func (s *tableSuite) TestSearchOne() {
	tb, _ := s.db.GetTableByName("Users")
	result, _ := tb.SearchOne("age", "54")
	if result.String() != "|1| 2 |2| juan |3| 54" {
		s.Fail("Expected |1| 2 |2| juan |3| 54", fmt.Sprintf("Recibe: %s", result))
	}
}
func (s *tableSuite) TestSearchAll() {
	tb, _ := s.db.GetTableByName("Users")
	result := tb.SearchAll("age", "54")
	if len(result) < 2 {

		s.Fail("Expected slice len greater than 2", fmt.Sprintf("Recibe: %d", len(result)))
	}
	if result[0].String() != "|1| 2 |2| juan |3| 54" {
		s.Fail("Expected |1| 2 |2| juan |3| 54", fmt.Sprintf("Recibe: %s", result))
	}
	if result[1].String() != "|1| 4 |2| manuel |3| 54" {
		s.Fail("Expected |1| 4 |2| manuel |3| 54", fmt.Sprintf("Recibe: %s", result))
	}
}

func (s *tableSuiteWithStaticData) TestAddForeignKeys() {
	fk := &pkg.ForeignKey{
		TableName:         "Users",
		ColumnName:        "id",
		ForeignTableName:  "Houses",
		ForeignColumnName: "id_owner",
	}
	err := s.db.AddForeignKey(*fk)
	if err != nil {
		s.Fail("Expected nil", fmt.Sprintf("Recibe: %s", err))
	}
}

func (s *tableSuiteWithStaticData) TestAddForeignKeys_ReturnTableNotFoundError() {
	fk := &pkg.ForeignKey{
		TableName:         "test",
		ColumnName:        "id",
		ForeignTableName:  "Houses",
		ForeignColumnName: "id_owner",
	}
	err := s.db.AddForeignKey(*fk)
	var example *pkg.NotFoundError
	if !errors.As(err, &example) {
	}
}

func (s *tableSuiteWithStaticData) TestSearchByForeignKey() {

	fk := &pkg.ForeignKey{
		TableName:         "Users",
		ColumnName:        "id",
		ForeignTableName:  "Houses",
		ForeignColumnName: "id_owner",
	}
	err := s.db.AddForeignKey(*fk)
	if err != nil {
		s.Fail("Expected nil", fmt.Sprintf("Recibe: %s", err))
	}

	tb, _ := s.db.GetTableByName("Users")
	complexRow, err := tb.SearchByForeignKey("1")
	if err != nil {
		s.Fail("Expected nil", fmt.Sprintf("Recibe: %s", err))
	}
	if len(complexRow[0].Rows) < 2 {
		s.Fail("Expected slice greater than 2", fmt.Sprintf("Recibe: %d", len(complexRow)))
	}
	if complexRow[0].Rows[0].String() != "|1| 1 |2| pedro avenue |3| 1" {
		s.Fail("Expected |1| 1 |2| pedro avenue |3| 1", fmt.Sprintf("Recibe: %s", complexRow[0].Rows[0]))
	}
	if complexRow[0].Rows[1].String() != "|1| 2 |2| pedro avenue |3| 1" {
		s.Fail("Expected |1| 2 |2| pedro avenue |3| 1", fmt.Sprintf("Recibe: %s", complexRow[0].Rows[1]))
	}

}

func getIdAndIndex(r pkg.Rows) (string, int) {
	var index int
	var id string
	for i, row := range r {
		value := row.SearchValue("name")
		if value == "test" {
			id = getId(row.String())
			index = i
			return id, index
		}
	}
	return "", 0
}

func TestTable(t *testing.T) {
	t.Run("Test Set: Tables", func(t *testing.T) {
		suite.Run(t, &tableSuite{})
	})

	t.Run("Test Set: TablesWithStaticData", func(t *testing.T) {
		suite.Run(t, &tableSuiteWithStaticData{})
	})
}
