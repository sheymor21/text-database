package Test

import (
	"github.com/stretchr/testify/suite"
	"os"
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
		s.Fail("Expected Users Table")

	}
}
func (s *Suite) TestGetTables() {
	tbs := s.db.GetTables()
	if len(tbs) != 1 {
		s.Fail("Expected 1 Table")
	}
	for _, t := range tbs {
		if t.GetName() != "-----Users-----" {
			s.Fail("Expected Users Table")
		}
	}
}
func (s *Suite) TestNewTable() {
	tb := s.db.NewTable("Test", []string{"name", "age"})
	if tb.GetName() != "-----Test-----" {
		s.Fail("Expected Test Table")
	}
}
func (s *Suite) TestDeleteTable() {
	s.db.NewTable("Test", []string{"name", "age"})
	s.db.DeleteTable("Test")
	tbs := s.db.GetTables()
	for _, t := range tbs {
		if t.GetName() == "-----Test-----" {
			s.Fail("Expected Test Table deleted")
		}
	}
}

func TestAll(t *testing.T) {
	suite.Run(t, new(Suite))
}
