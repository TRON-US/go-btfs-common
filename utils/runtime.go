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

type dbObj struct {
	dbConn dbConnection
}

type dbConnection interface {
	checkConnection(ctx context.Context, report *sharedpb.RuntimeInfoReport)
}

type redisObj struct {
	url string
}

type postgresObj struct {
	urls map[string]string
}

func CheckDBConnection(ctx context.Context, req *sharedpb.SignedRuntimeInfoRequest, connection db.ConnectionUrls) (*sharedpb.RuntimeInfoReport, error) {
	// db runtime
	report := new(sharedpb.RuntimeInfoReport)
	report.Status = sharedpb.RuntimeInfoReport_RUNNING
	report.DbStatusExtra = make(map[string]string)

	var checker []*dbObj

	checker = append(checker, &dbObj{dbConn: &postgresObj{urls: connection.PgURL}})
	checker = append(checker, &dbObj{dbConn: &redisObj{url: connection.RdURL}})

	for _, check := range checker {
		check.dbConn.checkConnection(ctx, report)
	}

	return report, nil
}

func (p *postgresObj) checkConnection(ctx context.Context, report *sharedpb.RuntimeInfoReport) {
	for key, url := range p.urls {
		// Assume the database connection is healthy
		report.DbStatusExtra[key] = DBURLDNE

		if url != "" {
			// Log the connection string
			pgdb := postgres.CreateTGPGDB(url)
			log.Info("Postgres",
				zap.String("name", key),
				zap.String("user", pgdb.Options().User),
				zap.String("host", pgdb.Options().Addr),
				zap.String("db", pgdb.Options().Database),
			)

			// Ping the database.
			// TODO: Separate function to test only ro and rw access
			if err := pgdb.Ping(); err != nil {
				report.DbStatusExtra[key] = constant.DBWriteConnectionError
				report.Status = sharedpb.RuntimeInfoReport_SICK
				log.Error(constant.DBWriteConnectionError, zap.Error(err))
			}

			// Set the database connection is healthy
			report.DbStatusExtra[key] = constant.DBConnectionHealthy
		}

		// Log status
		log.Info(key + ":" + report.DbStatusExtra[key])
	}
}

func (r *redisObj) checkConnection(ctx context.Context, report *sharedpb.RuntimeInfoReport) {

	// Assume the redis connection is not present
	report.RdStatusExtra = RDURLDNE

	// Check redis environment variable
	if r.url != "" {
		// Parse redis url
		opts, errParse := redis.ParseRedisURL(r.url)
		if errParse != nil {
			log.Error(constant.RDURLParseError, zap.Error(errParse))
		}

		// Log connection string
		log.Info("Redis URL",
			zap.String("host", opts.Addr),
			zap.String("db", string(opts.DB)),
		)

		// Check redis connection
		// TODO: Separate function to test only ro and rw access
		errConn := redis.CheckRedisConnection(redis.NewRedisConn(opts))
		if errConn != nil {
			report.RdStatusExtra = constant.RDConnectionError
			report.Status = sharedpb.RuntimeInfoReport_SICK
			log.Error(constant.RDConnectionError, zap.Error(errConn))
		}

		// Set redis connection to healthy
		report.RdStatusExtra = constant.RDConnectionHealthy
	}

	// Log status
	log.Info(report.RdStatusExtra)
}
