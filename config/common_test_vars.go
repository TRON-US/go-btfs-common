package config

import (
	"errors"
	"fmt"
	"os"
)

var (
	DbStatusURL         string
	DbGuardURL          string
	RdURL               string
	FoundPgStatusString bool
	FoundPgGuardString  bool
	FoundDbRdString     bool
)

func InitTestDB() error {
	//get db and redis connection strings
	DbStatusURL, FoundPgStatusString = os.LookupEnv("TEST_DB_URL_STATUS")
	DbGuardURL, FoundPgGuardString = os.LookupEnv("TEST_DB_URL_STATUS")
	RdURL, FoundDbRdString = os.LookupEnv("TEST_RD_URL")
	if FoundPgStatusString == false || FoundPgGuardString == false || FoundDbRdString == false {
		return errors.New(fmt.Sprintf("TEST_DB_URL_STATUS or TEST_DB_URL_STATUS or TEST_RD_URL env vars need to be set before running test"))
	}
	return nil
}
