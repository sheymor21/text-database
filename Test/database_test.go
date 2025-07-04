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
	s.db.DeleteTable("Test")
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
