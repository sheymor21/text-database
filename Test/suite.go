package Test

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"reflect"
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
func (s *tableSuite) ErrFail(err error) {
	expected := fmt.Sprintf("Expected %s", reflect.TypeOf(&pkg.NotFoundError{}))
	recibe := fmt.Sprintf("Recibe: %s", reflect.TypeOf(err))
	message := fmt.Sprintf("Message: %s", err.Error())
	s.Fail(expected, recibe, message)
}

func (s *databaseSuite) ErrFail(err error) {
	expected := fmt.Sprintf("Expected %s", reflect.TypeOf(&pkg.NotFoundError{}))
	recibe := fmt.Sprintf("Recibe: %s", reflect.TypeOf(err))
	message := fmt.Sprintf("Message: %s", err.Error())
	s.Fail(expected, recibe, message)
}
