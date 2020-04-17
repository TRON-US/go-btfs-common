package utils

import (
	"context"

	sharedpb "github.com/tron-us/go-btfs-common/protos/shared"

	"github.com/tron-us/go-common/v2/constant"
	"github.com/tron-us/go-common/v2/db"
	"github.com/tron-us/go-common/v2/db/postgres"
	"github.com/tron-us/go-common/v2/db/redis"
	"github.com/tron-us/go-common/v2/log"

	"go.uber.org/zap"
)

const DBURLDNE = "DB URL does not exist !!"
const RDURLDNE = "RD URL does not exist !!"

func CheckRuntime(ctx context.Context, runtime *sharedpb.SignedRuntimeInfoRequest, connection db.ConnectionUrls) (*sharedpb.RuntimeInfoReport, error) {
	// db runtime
	report := new(sharedpb.RuntimeInfoReport)
	report.Status = sharedpb.RuntimeInfoReport_RUNNING
	report.DbStatusExtra = map[string]string{}

	for key, url := range connection.PgURL {
		if url != "" {
			// Check postgres dbWrite
			PGDBWrite := postgres.CreateTGPGDB(url)
			if err := PGDBWrite.Ping(); err != nil {
				report.DbStatusExtra[key] = constant.DBWriteConnectionError
				report.Status = sharedpb.RuntimeInfoReport_SICK
				log.Error(constant.DBWriteConnectionError, zap.Error(err))
			}
			// Check postgres dbRead
			PGDBRead := postgres.CreateTGPGDB(url)
			if err := PGDBRead.Ping(); err != nil {
				report.DbStatusExtra[key] = constant.DBReadConnectionError
				report.Status = sharedpb.RuntimeInfoReport_SICK
				log.Error(constant.DBReadConnectionError, zap.Error(err))
			}
			// Assume the database connection is healthy
			report.DbStatusExtra[key] = constant.DBConnectionHealthy
		} else {
			report.DbStatusExtra[key] = DBURLDNE
		}
		if report.DbStatusExtra[key] == constant.DBConnectionHealthy {
			log.Info(key +":"+ constant.DBConnectionHealthy)
		}
	}

	// Check redis environment variable
	if connection.RdURL != "" {
		opts, errParse := redis.ParseRedisURL(connection.RdURL)
		if errParse != nil {
			log.Error(constant.RDURLParseError, zap.Error(errParse))
		}
		errConn := redis.CheckRedisConnection(redis.NewRedisConn(opts))
		if errConn != nil {
			report.RdStatusExtra = constant.RDConnectionError
			report.Status = sharedpb.RuntimeInfoReport_SICK
			log.Error(constant.RDConnectionError, zap.Error(errConn))
		}
		// Assume the redis connection is healthy
		report.RdStatusExtra = constant.RDConnectionHealthy
	} else {
		report.RdStatusExtra = RDURLDNE
	}
	if report.RdStatusExtra == constant.RDConnectionHealthy {
		log.Info(constant.RDConnectionHealthy)
	}

	// Remaining fields will be populated by the calling service
	// Reserve: only pass fatal error to higher level
	return report, nil
}
