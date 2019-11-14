package utils

import (
	"context"
	sharedpb "github.com/tron-us/go-btfs-common/protos/shared"
	"github.com/tron-us/go-common/constant"
	"github.com/tron-us/go-common/db"
	"github.com/tron-us/go-common/db/postgres"
	"github.com/tron-us/go-common/db/redis"
	"github.com/tron-us/go-common/log"

	"go.uber.org/zap"
)

func CheckRuntime(ctx context.Context, runtime *sharedpb.SignedRuntimeInfoRequest, connection db.ConnectionUrls) (*sharedpb.RuntimeInfoReport, error) {
	// db runtime
	report := new(sharedpb.RuntimeInfoReport)
	report.Status = sharedpb.RuntimeInfoReport_RUNNING

	if connection.PgURL != "" {
		// Check postgres dbWrite
		PGDBWrite := postgres.CreateTGPGDB(connection.PgURL)
		if err := PGDBWrite.Ping(); err != nil {
			report.DbStatusExtra = []byte(constant.DBWriteConnectionError)
			report.Status = sharedpb.RuntimeInfoReport_SICK
			log.Error(constant.DBWriteConnectionError, zap.Error(err))
		}
		// Check postgres dbRead
		PGDBRead := postgres.CreateTGPGDB(connection.PgURL)
		if err := PGDBRead.Ping(); err != nil {
			report.DbStatusExtra = []byte(constant.DBReadConnectionError)
			report.Status = sharedpb.RuntimeInfoReport_SICK
			log.Error(constant.DBReadConnectionError, zap.Error(err))
		}
		// Assume the database connection is healthy
		report.DbStatusExtra = []byte(constant.DBConnectionHealthy)
	} else {
		report.DbStatusExtra = nil
	}

	// Check redis environment variable
	if connection.RdURL != "" {
		err := redis.CheckRedisConnection(redis.CreateTGRDDB(connection.RdURL))
		if err != nil {
			report.RdStatusExtra = []byte(constant.RDConnectionError)
			report.Status = sharedpb.RuntimeInfoReport_SICK
			log.Error(constant.RDConnectionError, zap.Error(err))
		}
		// Assume the redis connection is healthy
		report.RdStatusExtra = []byte(constant.RDConnectionHealthy)
	} else {
		report.RdStatusExtra = nil
	}
	// Remaining fields will be populated by the calling service
	// Reserve: only pass fatal error to higher level
	return report, nil
}
