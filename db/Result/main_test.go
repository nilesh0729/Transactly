package Anuskh

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nilesh0729/Transactly/util"
)

var testQueries *Queries
var TestDb *sql.DB

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Can't Load Config: ", err)
	}

	TestDb, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Can't connect to Database:", err)
	}

	testQueries = New(TestDb)

	os.Exit(m.Run())

}
