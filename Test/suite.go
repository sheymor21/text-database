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
	db       pkg.Db
	dbConfig pkg.DbConfig
}

type databaseSuite struct {
	suite.Suite
	db       pkg.Db
	dbConfig pkg.DbConfig
}

var testConfig = []pkg.DbConfig{
	{SecurityKey: "", DatabaseName: "testDb.txt"},
	{SecurityKey: "testKey123", DatabaseName: "testDbWithKey.txt"},
}

func (s *databaseSuite) SetupTest() {
	config := pkg.DbConfig{SecurityKey: "", DatabaseName: "testDb.txt"}
	s.db = config.CreateDatabase()
}

func (s *databaseSuite) TearDownTest() {
	utilities.ErrorHandler(os.Remove("testDb.txt"))
}

func (s *tableSuite) SetupTest() {
	config := pkg.DbConfig{SecurityKey: "", DatabaseName: "testDb.txt"}
	s.db = config.CreateDatabase()
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
