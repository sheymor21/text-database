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
	row.Value = strings.TrimSpace(row.Value)
	if row.Value != "|1| 2 |2| juan |3| 54" {
		s.Fail("Expected |1| 2 |2| juan |3| 54", fmt.Sprintf("Recibe: %s", row))
	}
}
func (s *tableSuite) TestGetRows() {
	tb := utilities.Must(s.db.GetTableByName("Users"))
	count := len(tb.GetRows())
	if count != 3 {
		s.Fail("Expected 2 Rows", fmt.Sprintf("Recibe: %d", count))
	}
}
func (s *tableSuite) TestAddValue() {
	tb := utilities.Must(s.db.GetTableByName("Users"))
	tb, _ = tb.AddValue("name", "Jose")
	rows := tb.GetRows()
	id := getId(rows[3].Value)
	if rows[3].Value != fmt.Sprintf("|1| %s |2| Jose |3| null ", id) {
		s.Fail(fmt.Sprintf("Expected |1| %s |2| Jose |3| null", id), fmt.Sprintf("Recibe: %s", rows[2]))
	}

}
func (s *tableSuite) TestAddValues() {
	tb := utilities.Must(s.db.GetTableByName("Users"))
	tb = tb.AddValues([]string{"Jose", "20"})
	rows := tb.GetRows()
	id := getId(rows[3].Value)
	if rows[3].Value != fmt.Sprintf("|1| %s |2| Jose |3| 20 ", id) {
		s.Fail(fmt.Sprintf("Expected |1| %s |2| Jose |3| 20", id), fmt.Sprintf("Recibe: %s", rows[2]))
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
	if rows[2].Value != "|1| 2 |2| juan |3| 30 " {
		s.Fail("Expected |1| 2 |2| juan |3| 30", fmt.Sprintf("Recibe: %s", rows[2]))
	}
}
func (s *tableSuite) TestDeleteRow() {
	tb := utilities.Must(s.db.GetTableByName("Users"))
	tb, _ = tb.DeleteRow("1")
	rows := tb.GetRows()
	if rows[1].Value != "|1| 2 |2| juan |3| 54" || len(rows) != 2 {
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
func (s *tableSuite) TestAddValue_ReturnError() {
	tb, _ := s.db.GetTableByName("Users")
	_, err := tb.AddValue("test", "value")
	var example *pkg.NotFoundError
	if !errors.As(err, &example) {
		s.ErrFail(err)
	}
}
func (s *tableSuite) TestUpdateColumnName_ReturnError() {
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
func (s *tableSuite) TestDeleteRow_ReturnError() {
	tb, _ := s.db.GetTableByName("Users")
	_, err := tb.DeleteRow("test")
	var example *pkg.NotFoundError
	if !errors.As(err, &example) {
		s.ErrFail(err)
	}
}
func (s *tableSuite) TestDeleteColumn_ReturnError() {
	tb, _ := s.db.GetTableByName("Users")
	_, err := tb.DeleteColumn("test")
	var example *pkg.NotFoundError
	if !errors.As(err, &example) {
		s.ErrFail(err)
	}
}
func (s *tableSuite) TestGetRowById_ReturnError() {
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
	if newRow[1].Value != "|1| 2 |2| juan |3| 54" {
		s.Fail("Expected |1| 2 |2| juan |3| 54", fmt.Sprintf("Recibe: %s", newRow[1]))
	}
}
func (s *tableSuite) TestOrderByDescend_Numbers() {
	tb, _ := s.db.GetTableByName("Users")
	rows := tb.GetRows()
	newRow, _ := rows.OrderByDescend("age")
	if newRow[1].Value != "|1| 1 |2| pedro |3| 32" {
		s.Fail("Expected |1| 1 |2| pedro |3| 32", fmt.Sprintf("Recibe: %s", newRow[1]))
	}
}
func (s *tableSuite) TestOrderByAscend_Letters() {
	tb, _ := s.db.GetTableByName("Users")
	rows := tb.GetRows()
	newRow, _ := rows.OrderByAscend("name")
	if newRow[1].Value != "|1| 1 |2| pedro |3| 32" {
		s.Fail("Expected |1| 1 |2| pedro |3| 32", fmt.Sprintf("Recibe: %s", newRow[1]))
	}
}
func (s *tableSuite) TestOrderByDescend_Letters() {
	tb, _ := s.db.GetTableByName("Users")
	rows := tb.GetRows()
	newRow, _ := rows.OrderByDescend("name")
	if newRow[1].Value != "|1| 2 |2| juan |3| 54" {
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

func TestTable(t *testing.T) {
	for _, config := range testConfig {
		t.Run(fmt.Sprintf("DbConfig: %s", config.DatabaseName), func(t *testing.T) {
			suite.Run(t, &tableSuite{dbConfig: config})
		})
	}
}
