package utils

import (
	"context"
	sharedpb "github.com/tron-us/go-btfs-common/protos/shared"
	"github.com/tron-us/go-common/db/postgres"
	dbenv "github.com/tron-us/go-common/env/db"
	"github.com/tron-us/go-common/log"
	"go.uber.org/zap"
)

const (
	DBConnectionHealthy = "DB Connection Healthy"
	DBWriteConnectionError = "Cannot connect to the write database!"
	DBReadConnectionError  = "Cannot connect to the read database!"
)

func CheckRuntime(ctx context.Context, runtime *sharedpb.RuntimeInfoRequest) (*sharedpb.RuntimeInfoReport, error) {

	// db runtime
	report := new(sharedpb.RuntimeInfoReport)

	// Check postgres dbWrite
	report.DbStatusExtra = []byte(DBConnectionHealthy)
	report.Signature = runtime.Signature

	PGDBWrite := postgres.CreateTGPGDB(dbenv.DBWriteURL)
	if err := PGDBWrite.Ping(); err != nil {
		report.DbStatusExtra = []byte(DBWriteConnectionError)
		report.Signature = runtime.Signature
		log.Error(DBWriteConnectionError, zap.Error(err))
	}
	// Check postgres dbRead
	PGDBRead := postgres.CreateTGPGDB(dbenv.DBReadURL)
	if err := PGDBRead.Ping(); err != nil {
		report.DbStatusExtra = []byte(DBReadConnectionError)
		report.Status = sharedpb.RuntimeInfoReport_SICK
		log.Error(DBReadConnectionError, zap.Error(err))
	}
	//remaining fields will be populated by the calling service

	// Reserve: only pass fatal error to higher level
	return report, nil
}

