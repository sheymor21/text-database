package Test

import (
	"errors"
	"fmt"
	"github.com/sheymor21/text-database/tdb"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

func (s *tableSuite) TestGetRowById() {
	tb, _ := s.db.GetTableByName("Users")
	row, _ := tb.GetRowById("2")
	result := strings.TrimSpace(row.String())
	if result != "|1| 2 |2| juan |3| 54" {
		s.Fail("Expected |1| 2 |2| juan |3| 54", fmt.Sprintf("Recibe: %s", row))
	}
}
func (s *tableSuite) TestGetRows() {
	tb, _ := s.db.GetTableByName("Users")
	count := len(tb.GetRows())
	if count != 4 {
		s.Fail("Expected 4 Rows", fmt.Sprintf("Recibe: %d", count))
	}
}
func (s *tableSuite) TestAddValue() {
	tb, _ := s.db.GetTableByName("Users")
	_ = tb.AddValue("name", "test")
	tb, _ = s.db.GetTableByName("Users")
	rows := tb.GetRows()
	id, index := getIdAndIndex(rows)
	if rows[index].String() != fmt.Sprintf("|1| %s |2| test |3| null", id) {
		s.Fail(fmt.Sprintf("Expected |1| %s |2| test |3| null", id), fmt.Sprintf("Recibe: %s", rows[2]))
	}

}
func (s *tableSuite) TestAddValues() {
	tb, _ := s.db.GetTableByName("Users")
	tb.AddValues("test", "20")
	tb, _ = s.db.GetTableByName("Users")
	rows := tb.GetRows()
	id, index := getIdAndIndex(rows)
	if rows[index].String() != fmt.Sprintf("|1| %s |2| test |3| 20", id) {
		s.Fail(fmt.Sprintf("Expected |1| %s |2| test |3| 20", id), fmt.Sprintf("Recibe: %s", rows[2]))
	}
}
func (s *tableSuite) TestUpdateTableName() {
	tb, _ := s.db.GetTableByName("Users")
	tb.UpdateTableName("Test")
	tb, _ = s.db.GetTableByName("Test")
	if tb.GetName() != "-----Test-----" {
		s.Fail("Expected Test Table", fmt.Sprintf("Recibe: %s", tb))
	}
}
func (s *tableSuite) TestUpdateColumnName() {
	tb, _ := s.db.GetTableByName("Users")
	_ = tb.UpdateColumnName("name", "username")
	tb, _ = s.db.GetTableByName("Users")
	columns := strings.TrimSpace(strings.Join(tb.GetColumns(), " "))
	if columns != "[1] id [2] username [3] age" {
		s.Fail("Expected [1] id [2] username |3| age ", fmt.Sprintf("Recibe: %s", columns))
	}
}
func (s *tableSuite) TestUpdateValue() {
	tb, _ := s.db.GetTableByName("Users")
	_ = tb.UpdateValue("age", "2", "30")
	tb, _ = s.db.GetTableByName("Users")
	rows := tb.GetRows()
	if rows[1].String() != "|1| 2 |2| juan |3| 30" {
		s.Fail("Expected |1| 2 |2| juan |3| 30", fmt.Sprintf("Recibe: %s", rows[2]))
	}
}
func (s *tableSuite) TestDeleteRow() {
	tb, _ := s.db.GetTableByName("Users")
	_ = tb.DeleteRow("1", false)
	tb, _ = s.db.GetTableByName("Users")
	rows := tb.GetRows()
	if len(rows) != 3 {
		s.Fail("Expected len of 3", fmt.Sprintf("Recibe: %d", len(rows)))
	}

	if rows[0].String() != "|1| 2 |2| juan |3| 54" {
		s.Fail("Expected |1| 1 |2| juan |3| 54", fmt.Sprintf("Recibe: %s", rows[1]))
	}
}
func (s *tableSuite) TestDeleteColumn() {
	tb, _ := s.db.GetTableByName("Users")
	_ = tb.DeleteColumn("age")
	tb, _ = s.db.GetTableByName("Users")
	columns := strings.TrimSpace(strings.Join(tb.GetColumns(), " "))
	if columns != "[1] id [2] name" {
		s.Fail("Expected [1] id [2] name", fmt.Sprintf("Recibe: %s", columns))
	}
}
func (s *tableSuite) TestAddValue_ReturnColumnError() {
	tb, _ := s.db.GetTableByName("Users")
	err := tb.AddValue("test", "value")
	var example *tdb.NotFoundError
	if !errors.As(err, &example) {
		s.ErrFail(err)
	}
}
func (s *tableSuite) TestUpdateColumnName_ReturnColumnError() {
	tb, _ := s.db.GetTableByName("Users")
	err := tb.UpdateColumnName("test", "value")
	var example *tdb.NotFoundError
	if !errors.As(err, &example) {
		s.ErrFail(err)
	}
}
func (s *tableSuite) TestUpdateValue_ReturnColumnError() {
	tb, _ := s.db.GetTableByName("Users")
	err := tb.UpdateValue("test", "value", "value")
	var example *tdb.NotFoundError
	if !errors.As(err, &example) {
		s.ErrFail(err)
	}
}
func (s *tableSuite) TestUpdateValue_ReturnIdError() {
	tb, _ := s.db.GetTableByName("Users")
	err := tb.UpdateValue("name", "test", "value")
	var example *tdb.NotFoundError
	if !errors.As(err, &example) {
		s.ErrFail(err)
	}
}
func (s *tableSuite) TestDeleteRow_ReturnIdError() {
	tb, _ := s.db.GetTableByName("Users")
	err := tb.DeleteRow("test", false)
	var example *tdb.NotFoundError
	if !errors.As(err, &example) {
		s.ErrFail(err)
	}
}
func (s *tableSuite) TestDeleteColumn_ReturnIdError() {
	tb, _ := s.db.GetTableByName("Users")
	err := tb.DeleteColumn("test")
	var example *tdb.NotFoundError
	if !errors.As(err, &example) {
		s.ErrFail(err)
	}
}
func (s *tableSuite) TestGetRowById_ReturnIdError() {
	tb, _ := s.db.GetTableByName("Users")
	_, err := tb.GetRowById("test")
	var example *tdb.NotFoundError
	if !errors.As(err, &example) {
		s.ErrFail(err)
	}
}
func (s *tableSuite) TestOrderByAscend_Numbers() {
	tb, _ := s.db.GetTableByName("Users")
	rows := tb.GetRows()
	_ = rows.OrderByAscend("age")
	if rows[1].String() != "|1| 2 |2| juan |3| 54" {
		s.Fail("Expected |1| 2 |2| juan |3| 54", fmt.Sprintf("Recibe: %s", rows[1]))
	}
}
func (s *tableSuite) TestOrderByDescend_Numbers() {
	tb, _ := s.db.GetTableByName("Users")
	rows := tb.GetRows()
	_ = rows.OrderByDescend("age")
	if rows[0].String() != "|1| 1 |2| pedro |3| 32" {
		s.Fail("Expected |1| 1 |2| pedro |3| 32", fmt.Sprintf("Recibe: %s", rows[1]))
	}
}
func (s *tableSuite) TestOrderByAscend_Letters() {
	tb, _ := s.db.GetTableByName("Users")
	rows := tb.GetRows()
	_ = rows.OrderByAscend("name")
	if rows[0].String() != "|1| 1 |2| pedro |3| 32" {
		s.Fail("Expected |1| 1 |2| pedro |3| 32", fmt.Sprintf("Recibe: %s", rows[1]))
	}
}
func (s *tableSuite) TestOrderByDescend_Letters() {
	tb, _ := s.db.GetTableByName("Users")
	rows := tb.GetRows()
	_ = rows.OrderByDescend("name")
	if rows[1].String() != "|1| 2 |2| juan |3| 54" {
		s.Fail("Expected |1| 2 |2| juan |3| 54", fmt.Sprintf("Recibe: %s", rows[1]))
	}
}
func (s *tableSuite) TestOrderBy_ReturnColumnError() {
	tb, _ := s.db.GetTableByName("Users")
	rows := tb.GetRows()
	err := rows.OrderByAscend("Email")
	var example *tdb.NotFoundError
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
	fk := &tdb.ForeignKey{
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
	fk := &tdb.ForeignKey{
		TableName:         "test",
		ColumnName:        "id",
		ForeignTableName:  "Houses",
		ForeignColumnName: "id_owner",
	}
	err := s.db.AddForeignKey(*fk)
	var example *tdb.NotFoundError
	if !errors.As(err, &example) {
	}
}
func (s *tableSuiteWithStaticData) TestSearchByForeignKey() {

	fk := &tdb.ForeignKey{
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
func (s *tableSuiteWithStaticData) TestDeleteByForeignKey() {
	fk := &tdb.ForeignKey{
		TableName:         "Users",
		ColumnName:        "id",
		ForeignTableName:  "Houses",
		ForeignColumnName: "id_owner",
	}
	_ = s.db.AddForeignKey(*fk)

	usersTb, _ := s.db.GetTableByName("Users")
	err := usersTb.DeleteRow("1", true)
	if err != nil {
		s.Fail(err.Error())
	}
	housesTb, _ := s.db.GetTableByName("Houses")
	rows := housesTb.GetRows()
	if len(rows) != 2 {
		s.Fail("Expected len of 2", fmt.Sprintf("Recibe: %d", len(rows)))
	}
	if rows[0].String() != "|1| 3 |2| juan avenue |3| 2" {

		s.Fail("Expected |1| 3 |2| juan avenue |3| 2", fmt.Sprintf("Recibe: %d", len(rows)))
	}
	if rows[1].String() != "|1| 4 |2| carlos avenue |3| 3" {
		s.Fail("Expected |1| 4 |2| carlos avenue |3| 3", fmt.Sprintf("Recibe: %d", len(rows)))
	}
}

func getIdAndIndex(r tdb.Rows) (string, int) {
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
