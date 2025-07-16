package Test

import (
	"fmt"
	"github.com/sheymor21/text-database/tdb"
	"github.com/sheymor21/text-database/tdb/utilities"
	"github.com/stretchr/testify/suite"
	"os"
	"reflect"
)

type databaseSuite struct {
	suite.Suite
	db       tdb.Db
	dbConfig tdb.DbConfig
}

type tableSuite struct {
	suite.Suite
	db       tdb.Db
	dbConfig tdb.DbConfig
}

type tableSuiteWithStaticData struct {
	tableSuite
}
type databaseWithEncryptionSuite struct {
	databaseSuite
}
type databaseWithStaticDataSuite struct {
	suite.Suite
	db       tdb.Db
	dbConfig tdb.DbConfig
}

func (s *databaseSuite) SetupTest() {

	config := tdb.DbConfig{EncryptionKey: "", DatabaseName: "testDb.txt"}
	s.db, _ = config.CreateDatabase()
}

func (s *databaseSuite) TearDownTest() {
	utilities.ErrorHandler(os.Remove("testDb.txt"))
}

func (s *tableSuite) SetupTest() {
	config := tdb.DbConfig{EncryptionKey: "", DatabaseName: "testDb.txt"}
	s.db, _ = config.CreateDatabase()
}

func (s *tableSuite) TearDownTest() {
	utilities.ErrorHandler(os.Remove("testDb.txt"))
}
func (s *tableSuiteWithStaticData) SetupTest() {
	dataConfig := []tdb.DataConfig{
		{
			TableName: "Users",
			Columns:   []string{"name", "age"},
			Values: []tdb.Values{
				{"1", "pedro", "32"},
				{"2", "juan", "54"},
				{"3", "carlos", "62"},
				{"4", "manuel", "54"},
			},
		},
		{
			TableName: "DataTest",
			Columns:   []string{"name", "age"},
			Values:    []tdb.Values{{"1", "carlos", "32"}, {"2", "jose", "23"}},
		},
		{
			TableName: "Houses",
			Columns:   []string{"direction", "id_owner"},
			Values: []tdb.Values{
				{"1", "pedro avenue", "1"},
				{"2", "pedro avenue", "1"},
				{"3", "juan avenue", "2"},
				{"4", "carlos avenue", "3"},
			},
		},
	}

	config := tdb.DbConfig{EncryptionKey: "", DatabaseName: "testDbTableWithStaticData.txt", DataConfig: dataConfig}
	s.db, _ = config.CreateDatabase()
}

func (s *tableSuiteWithStaticData) TearDownTest() {
	utilities.ErrorHandler(os.Remove("testDbTableWithStaticData.txt"))
}
func (s *databaseWithEncryptionSuite) SetupTest() {
	config := tdb.DbConfig{EncryptionKey: "", DatabaseName: "testDbWithEncryption.txt"}
	s.db, _ = config.CreateDatabase()
}

func (s *databaseWithEncryptionSuite) TearDownTest() {
	utilities.ErrorHandler(os.Remove("testDbWithEncryption.txt"))
}

func (s *databaseWithStaticDataSuite) SetupTest() {
	dataConfig := []tdb.DataConfig{
		{
			TableName: "DataTest",
			Columns:   []string{"name", "age"},
			Values:    []tdb.Values{{"1", "carlos", "32"}, {"2", "jose", "23"}},
		},
	}
	config := tdb.DbConfig{EncryptionKey: "", DatabaseName: "testDbWithStaticData.txt", DataConfig: dataConfig}
	s.db = utilities.Must(config.CreateDatabase())
}

func (s *databaseWithStaticDataSuite) TearDownTest() {
	utilities.ErrorHandler(os.Remove("testDbWithStaticData.txt"))
}

func (s *tableSuite) ErrFail(err error) {
	expected := fmt.Sprintf("Expected %s", reflect.TypeOf(&tdb.NotFoundError{}))
	recibe := fmt.Sprintf("Recibe: %s", reflect.TypeOf(err))
	message := fmt.Sprintf("Message: %s", err.Error())
	s.Fail(expected, recibe, message)
}

func (s *databaseSuite) ErrFail(err error) {
	expected := fmt.Sprintf("Expected %s", reflect.TypeOf(&tdb.NotFoundError{}))
	recibe := fmt.Sprintf("Recibe: %s", reflect.TypeOf(err))
	message := fmt.Sprintf("Message: %s", err.Error())
	s.Fail(expected, recibe, message)
}

func (s *databaseWithStaticDataSuite) ErrFail(err error) {
	expected := fmt.Sprintf("Expected %s", reflect.TypeOf(&tdb.NotFoundError{}))
	recibe := fmt.Sprintf("Recibe: %s", reflect.TypeOf(err))
	message := fmt.Sprintf("Message: %s", err.Error())
	s.Fail(expected, recibe, message)
}
