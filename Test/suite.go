package Test

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"reflect"
	"text-database/pkg"
	"text-database/pkg/utilities"
)

type databaseSuite struct {
	suite.Suite
	db       pkg.Db
	dbConfig pkg.DbConfig
}

type tableSuite struct {
	suite.Suite
	db       pkg.Db
	dbConfig pkg.DbConfig
}
type databaseWithEncryptionSuite struct {
	databaseSuite
}
type databaseWithStaticDataSuite struct {
	suite.Suite
	db       pkg.Db
	dbConfig pkg.DbConfig
}

func (s *databaseSuite) SetupTest() {

	config := pkg.DbConfig{EncryptionKey: "", DatabaseName: "testDb.txt"}
	s.db, _ = config.CreateDatabase()
}

func (s *databaseSuite) TearDownTest() {
	utilities.ErrorHandler(os.Remove("testDb.txt"))
}

func (s *tableSuite) SetupTest() {
	config := pkg.DbConfig{EncryptionKey: "", DatabaseName: "testDb.txt"}
	s.db, _ = config.CreateDatabase()
}

func (s *tableSuite) TearDownTest() {
	utilities.ErrorHandler(os.Remove("testDb.txt"))
}

func (s *databaseWithEncryptionSuite) SetupTest() {
	config := pkg.DbConfig{EncryptionKey: "", DatabaseName: "testDbWithEncryption.txt"}
	s.db, _ = config.CreateDatabase()
}

func (s *databaseWithEncryptionSuite) TearDownTest() {
	utilities.ErrorHandler(os.Remove("testDbWithEncryption.txt"))
}

func (s *databaseWithStaticDataSuite) SetupTest() {
	dataConfig := []pkg.DataConfig{
		{
			TableName: "DataTest",
			Columns:   []string{"name", "age"},
			Values:    []pkg.Values{{"1", "carlos", "32"}, {"2", "jose", "23"}},
		},
	}
	config := pkg.DbConfig{EncryptionKey: "", DatabaseName: "testDbWithStaticData.txt", DataConfig: dataConfig}
	s.db = utilities.Must(config.CreateDatabase())
}

func (s *databaseWithStaticDataSuite) TearDownTest() {
	utilities.ErrorHandler(os.Remove("testDbWithStaticData.txt"))
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

func (s *databaseWithStaticDataSuite) ErrFail(err error) {
	expected := fmt.Sprintf("Expected %s", reflect.TypeOf(&pkg.NotFoundError{}))
	recibe := fmt.Sprintf("Recibe: %s", reflect.TypeOf(err))
	message := fmt.Sprintf("Message: %s", err.Error())
	s.Fail(expected, recibe, message)
}
