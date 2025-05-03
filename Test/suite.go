package Test

import (
	"github.com/stretchr/testify/suite"
	"os"
	"text-database/pkg"
	"text-database/pkg/utilities"
)

type tableSuite struct {
	suite.Suite
	db pkg.Db
}

type databaseSuite struct {
	suite.Suite
	db pkg.Db
}

func (s *databaseSuite) SetupTest() {
	s.db = pkg.CreateDatabase("testDb.txt")
}

func (s *databaseSuite) TearDownTest() {
	utilities.ErrorHandler(os.Remove("testDb.txt"))
}

func (s *tableSuite) SetupTest() {
	s.db = pkg.CreateDatabase("testDb.txt")
}

func (s *tableSuite) TearDownTest() {
	utilities.ErrorHandler(os.Remove("testDb.txt"))
}
