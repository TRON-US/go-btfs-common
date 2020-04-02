package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/tron-us/go-btfs-common/protos/guard"
	"github.com/tron-us/go-common/v2/env"
	"github.com/tron-us/go-common/v2/log"

	"go.uber.org/zap"
)

var (
	ConstMinQuestionsCountPerChallenge = 100
	ConstMinQuestionsCountPerShard     = 100
	DbStatusURL string
	DbGuardURL string
	DbRdURL string
	FoundPgStatusString bool
	FoundPgGuardString bool
	FoundDbRdString bool
)

func init() {
	if _, s := env.GetEnv("CONST_MIN_QUESTIONS_PER_CHALLENGE"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil {
			log.Panic("Get const_min_questions_per_challenge as string", zap.Error(err))
		}
		ConstMinQuestionsCountPerChallenge = v
	}
}

func GetMinimumQuestionsCountPerShard(status *guard.FileStoreStatus) (val int) {
	//Below comment codes were reserved here
	//if status == nil || status.FileStoreMeta.CheckFrequency == 0 {
	//	//if it is 0 or status is nil, means using default frequency
	//	return ConstMinQuestionsCountPerChallenge * 52
	//}
	//
	//return int(status.FileStoreMeta.CheckFrequency) * ConstMinQuestionsCountPerChallenge
	return ConstMinQuestionsCountPerShard
}

func InitDB() {
	//get db and redis connection strings
	DbStatusURL, FoundPgStatusString = os.LookupEnv("TEST_DB_URL_STATUS")
	DbGuardURL, FoundPgGuardString = os.LookupEnv("TEST_DB_URL_GUARD")
	DbRdURL, FoundDbRdString = os.LookupEnv("TEST_RD_URL")
	if FoundPgStatusString == false || FoundPgGuardString == false || FoundDbRdString == false {
		log.Error(fmt.Sprintf("dbStatusURL and dbStatusURL env vars need to be set before running test"))
	}
}
