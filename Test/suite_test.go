package Test

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"strings"
	"testing"
	"text-database/pkg"
	"text-database/pkg/utilities"
)

type Suite struct {
	suite.Suite
	db pkg.Db
}

func (s *Suite) SetupTest() {
	s.db = pkg.CreateDatabase("testDb.txt")
}

func (s *Suite) TearDownTest() {
	utilities.ErrorHandler(os.Remove("testDb.txt"))
}
func (s *Suite) TestGetTableByName() {
	tb := s.db.GetTableByName("Users")
	if tb.GetName() != "-----Users-----" {
		s.Fail("Expected Users Table", fmt.Sprintf("Recibe: %s", tb))

	}
}
func (s *Suite) TestGetTables() {
	tbs := s.db.GetTables()
	if len(tbs) != 1 {
		s.Fail("Expected 1 Table")
	}
	for _, t := range tbs {
		if t.GetName() != "-----Users-----" {
			s.Fail("Expected Users Table", fmt.Sprintf("Recibe: %s", t))
		}
	}
}
func (s *Suite) TestGetRowById() {
	tb := s.db.GetTableByName("Users")
	row := tb.GetRowById("2")
	row = strings.TrimSpace(row)
	if row != "|1| 2 |2| juan |3| 54" {
		s.Fail("Expected |1| 2 |2| juan |3| 54", fmt.Sprintf("Recibe: %s", row))
	}
}
func (s *Suite) TestGetRows() {
	tb := s.db.GetTableByName("Users")
	count := len(tb.GetRows())
	if count != 2 {
		s.Fail("Expected 2 Rows", fmt.Sprintf("Recibe: %d", count))
	}
}
func (s *Suite) TestNewTable() {
	tb := s.db.NewTable("Test", []string{"name", "age"})
	if tb.GetName() != "-----Test-----" {
		s.Fail("Expected Test Table", fmt.Sprintf("Recibe: %s", tb))
	}
}
func (s *Suite) TestDeleteTable() {
	s.db.NewTable("Test", []string{"name", "age"})
	s.db.DeleteTable("Test")
	tbs := s.db.GetTables()
	for _, t := range tbs {
		if t.GetName() == "-----Test-----" {
			s.Fail("Expected Test Table deleted", fmt.Sprintf("Recibe: %s", t))
		}
	}
}
func (s *Suite) TestAddValue() {
	tb := s.db.GetTableByName("Users")
	tb = tb.AddValue("name", "Jose")
	rows := tb.GetRows()
	id := getId(rows[2])
	if rows[2] != fmt.Sprintf("|1| %s |2| Jose |3| null ", id) {
		s.Fail(fmt.Sprintf("Expected |1| %s |2| Jose |3| null", id), fmt.Sprintf("Recibe: %s", rows[2]))
	}

}
func (s *Suite) TestAddValues() {
	tb := s.db.GetTableByName("Users")
	tb = tb.AddValues([]string{"Jose", "20"})
	rows := tb.GetRows()
	id := getId(rows[2])
	if rows[2] != fmt.Sprintf("|1| %s |2| Jose |3| 20 ", id) {
		s.Fail(fmt.Sprintf("Expected |1| %s |2| Jose |3| 20", id), fmt.Sprintf("Recibe: %s", rows[2]))
	}
}
func (s *Suite) TestUpdateTableName() {
	tb := s.db.GetTableByName("Users")
	tb = tb.UpdateTableName("Test")
	if tb.GetName() != "-----Test-----" {
		s.Fail("Expected Test Table", fmt.Sprintf("Recibe: %s", tb))
	}
}
func (s *Suite) TestUpdateColumnName() {
	tb := s.db.GetTableByName("Users")
	tb = tb.UpdateColumnName("name", "username")
	columns := strings.TrimSpace(strings.Join(tb.GetColumns(), " "))
	if columns != "[1] id [2] username [3] age" {
		s.Fail("Expected [1] id [2] username |3| age ", fmt.Sprintf("Recibe: %s", columns))
	}
}
func (s *Suite) TestUpdateValue() {
	tb := s.db.GetTableByName("Users")
	tb = tb.UpdateValue("age", "2", "30")
	rows := tb.GetRows()
	if rows[1] != "|1| 2 |2| juan |3| 30 " {
		s.Fail("Expected |1| 2 |2| juan |3| 30", fmt.Sprintf("Recibe: %s", rows[1]))
	}
}
func (s *Suite) TestDeleteRow() {
	tb := s.db.GetTableByName("Users")
	tb = tb.DeleteRow("1")
	rows := tb.GetRows()
	if rows[0] != "|1| 2 |2| juan |3| 54" || len(rows) != 1 {
		s.Fail("Expected |1| 1 |2| juan |3| 54", fmt.Sprintf("Recibe: %s", rows[0]))
	}
}
func (s *Suite) TestDeleteColumn() {
	tb := s.db.GetTableByName("Users")
	tb = tb.DeleteColumn("age")
	columns := strings.TrimSpace(strings.Join(tb.GetColumns(), " "))
	if columns != "[1] id [2] name" {
		s.Fail("Expected [1] id [2] name", fmt.Sprintf("Recibe: %s", columns))
	}
}
func TestDatabase(t *testing.T) {
	suite.Run(t, new(Suite))
}
func getId(row string) string {
	split := strings.Split(row, "|")
	return strings.TrimSpace(split[2])
}
