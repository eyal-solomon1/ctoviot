package db

import (
	"database/sql"
	util "github.com/eyal-solomon1/ctoviot/util"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error
	var config util.Config

	config, err = util.LoadConfig(util.WithConfigFileName("conf.test"), util.WithConfigFilePath("../.."))

	if err != nil {
		log.Fatal("couldn't find test config file", err)
	}
	testDb, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(testDb)

	os.Exit(m.Run())
}
