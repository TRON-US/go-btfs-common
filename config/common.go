package config

import (
	"github.com/tron-us/go-btfs-common/protos/guard"
	"strconv"

	"github.com/tron-us/go-common/v2/env"
	"github.com/tron-us/go-common/v2/log"

	"go.uber.org/zap"
)

var (
	ConstMinQuestionsCountPerChallenge = 100
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
	if status == nil || status.FileStoreMeta.CheckFrequency == 0 {
		//if it is 0 or status is nil, means using default frequency, 12 per year
		return ConstMinQuestionsCountPerChallenge * 12
	}

	return int(status.FileStoreMeta.CheckFrequency) * ConstMinQuestionsCountPerChallenge
}
