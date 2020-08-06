package config

import (
	"strconv"

	"github.com/tron-us/go-btfs-common/protos/guard"
	"github.com/tron-us/go-common/v2/env"
	"github.com/tron-us/go-common/v2/log"

	"go.uber.org/zap"
)

var (
	ConstMinQuestionsCountPerChallenge = 100
	ConstMinQuestionsCountPerShard     = 100

	ConstRequestPayoutBatchPageSize = 100
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

func GetRenewContingencyPercentage() (percent int) {
	return 10
}
