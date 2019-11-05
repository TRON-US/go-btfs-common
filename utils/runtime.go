package utils

import (
	"context"

	"github.com/tron-us/go-common/db/postgres"
	dbenv "github.com/tron-us/go-common/env/db"

	sharedpb "github.com/tron-us/go-btfs-common/protos/shared"
	"github.com/tron-us/go-common/log"

	"go.uber.org/zap"
)

const (
	DBConnectionHealthy = "DB Connection Healthy"
	DBWriteConnectionError = "Cannot connect to the write database!"
	DBReadConnectionError  = "Cannot connect to the read database!"
)

type Server struct {
}

func (s *Server) CheckRuntime(ctx context.Context, runtime *sharedpb.RuntimeInfoRequest) (*sharedpb.RuntimeInfoReport, error) {

	report := new(sharedpb.RuntimeInfoReport)
	report.DbStatusExtra = []byte(DBConnectionHealthy)
	report.Status = sharedpb.RuntimeInfoReport_RUNNING
	// db runtime
	// Check postgres dbWrite
	PGDBWrite := postgres.CreateTGPGDB(dbenv.DBWriteURL)
	if err := PGDBWrite.Ping(); err != nil {
		report.DbStatusExtra = []byte(DBWriteConnectionError)
		report.Status = sharedpb.RuntimeInfoReport_SICK
		log.Error(DBWriteConnectionError, zap.Error(err))
	}
	// Check postgres dbRead
	PGDBRead := postgres.CreateTGPGDB(dbenv.DBReadURL)
	if err := PGDBRead.Ping(); err != nil {
		report.DbStatusExtra = []byte(DBReadConnectionError)
		report.Status = sharedpb.RuntimeInfoReport_SICK
		log.Error(DBReadConnectionError, zap.Error(err))
	}
	// Reserve: only pass fatal error to higher level
	return report, nil
}

