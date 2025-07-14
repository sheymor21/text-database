package Test

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
	"text-database/pkg"
	"text-database/pkg/utilities"
)

func (s *databaseSuite) TestGetTables() {
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
func (s *databaseSuite) TestNewTable() {
	tb := s.db.NewTable("Test", []string{"name", "age"})
	if tb.GetName() != "-----Test-----" {
		s.Fail("Expected Test Table", fmt.Sprintf("Recibe: %s", tb))
	}
}
func (s *databaseSuite) TestDeleteTable() {
	s.db.NewTable("Test", []string{"name", "age"})
	err := s.db.DeleteTable("Test")
	if err != nil {
		s.ErrFail(err)
	}
	tbs := s.db.GetTables()
	for _, t := range tbs {
		if t.GetName() == "-----Test-----" {
			s.Fail("Expected Test Table deleted", fmt.Sprintf("Recibe: %s", t))
		}
	}
}
func (s *databaseWithStaticDataSuite) TestStaticData() {
	tb, err := s.db.GetTableByName("DataTest")
	if err != nil {
		s.ErrFail(err)
	}
	rows := tb.GetRows()
	if rows[0].String() != "|1| 1 |2| carlos |3| 32" && rows[1].String() != "|1| 2 |2| jose |3| 23" {
		s.Fail("Expected |1| 1 |2| carlos |3| 32 and |1| 2 |2| jose |3| 23", fmt.Sprintf("Recibe: %s", rows))
	}
}
func (s *databaseSuite) TestGetTableByName() {
	tb := utilities.Must(s.db.GetTableByName("Users"))
	if tb.GetName() != "-----Users-----" {
		s.Fail("Expected Users Table", fmt.Sprintf("Recibe: %s", tb))

	}
}
func (s *databaseSuite) TestGetTableByName_ReturnNameError() {
	_, err := s.db.GetTableByName("test")
	var example *pkg.NotFoundError
	if !errors.As(err, &example) {
		s.ErrFail(err)
	}
}

func (s *databaseSuite) TestFromSql_Select_All() {
	tb := utilities.Must(s.db.GetTableByName("Users"))
	data, _ := s.db.FromSql("SELECT * FROM Users")
	count := len(tb.GetRows())
	if len(data.Rows) != count {
		s.Fail("Expected len of 4", fmt.Sprintf("Recibe: %d", len(data.Rows)))
	}
}
func (s *databaseSuite) TestFromSql_Select() {
	tb := utilities.Must(s.db.GetTableByName("Users"))
	data, _ := s.db.FromSql("SELECT name , age FROM Users")
	count := len(tb.GetRows())
	if len(data.Rows) != count {
		s.Fail("Expected len of 4", fmt.Sprintf("Recibe: %d", len(data.Rows)))
	}
}
func (s *databaseSuite) TestFromSql_Select_Where() {
	data, err := s.db.FromSql("SELECT name , age FROM Users WHERE age = 54")

	if err != nil {
		s.ErrFail(err)
	}
	if len(data.Rows) != 2 {
		s.Fail("Expected len of 2", fmt.Sprintf("Recibe: %d", len(data.Rows)))
	}
}

func (s *databaseSuite) TestFromSql_Update() {
	data, err := s.db.FromSql("Update Users SET age = 25,name = pepe WHERE age = 32")
	if err != nil {
		s.ErrFail(err)
	}
	if data.AffectRows != 1 {
		s.Fail("Expected 1 affected row", fmt.Sprintf("Recibe: %d", data.AffectRows))
	}
	tb, _ := s.db.GetTableByName("Users")
	user, _ := tb.GetRowById("1")
	if user.String() != "|1| 1 |2| pepe |3| 25" {
		s.Fail("Expected |1| 1 |2| pepe |3| 25", fmt.Sprintf("Recibe: %s", user))
	}
}
func (s *databaseSuite) TestFromSql_Delete() {
	data, err := s.db.FromSql("Delete FROM Users WHERE age =54")

	if err != nil {
		s.ErrFail(err)
	}
	if data.AffectRows != 2 {
		s.Fail("Expected 2 affected row", fmt.Sprintf("Recibe: %d", data.AffectRows))
	}
	tb, _ := s.db.GetTableByName("Users")
	rows := tb.GetRows()
	if len(rows) != 2 {
		s.Fail("Expected len of 2", fmt.Sprintf("Recibe: %d", len(rows)))
	}
}
func (s *databaseSuite) TestFromSql_Insert() {
	data, err := s.db.FromSql("Insert INTO Users (id,name,age) VALUES (5,maria,20,6,carlitos,32)")
	if err != nil {
		s.ErrFail(err)
	}
	if data.AffectRows != 2 {
		s.Fail("Expected 1 affected row", fmt.Sprintf("Recibe: %d", data.AffectRows))
	}
	tb, _ := s.db.GetTableByName("Users")
	rows := tb.GetRows()
	if len(rows) != 6 {
		s.Fail("Expected len of 6", fmt.Sprintf("Recibe: %d", len(rows)))
	}
}
func TestDatabase(t *testing.T) {
	t.Run("TestSet: Database", func(t *testing.T) {
		suite.Run(t, &databaseSuite{})
	})

	t.Run("TestSet: DatabaseWithEncryption", func(t *testing.T) {
		suite.Run(t, &databaseWithEncryptionSuite{})
	})

	t.Run("TestSet: DatabaseWithStaticData", func(t *testing.T) {
		suite.Run(t, &databaseWithStaticDataSuite{})
	})
}
