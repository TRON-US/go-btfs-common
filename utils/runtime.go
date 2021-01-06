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

// DBURLDNE database does not exist or unreachable
const DBURLDNE = "DB URL does not exist !!"

// RDURLDNE redis does not exist or unreachable
const RDURLDNE = "RD URL does not exist !!"

type dbObj struct {
	dbConn DbConnection
}

type DbConnection interface {
	checkConnection(ctx context.Context, report *sharedpb.RuntimeInfoReport)
}

type RedisObj struct {
	url string
	DB  *redis.TGRDDB
}

type PostgresObj struct {
	urls map[string]string
	DBs  map[string]*postgres.TGPGDB
}

// CheckDBConnection checks database connections.
func CheckDBConnection(ctx context.Context, runtime *sharedpb.SignedRuntimeInfoRequest,
	connection db.ConnectionUrls) (*sharedpb.RuntimeInfoReport, []DbConnection, error) {
	// db runtime
	report := new(sharedpb.RuntimeInfoReport)
	report.Status = sharedpb.RuntimeInfoReport_RUNNING
	report.DbStatusExtra = make(map[string]string)

	var checker []*dbObj
	checker = append(checker,
		&dbObj{dbConn: &PostgresObj{
			urls: connection.PgURL,
			DBs:  map[string]*postgres.TGPGDB{},
		}},
		&dbObj{dbConn: &RedisObj{
			url: connection.RdURL,
			DB:  nil,
		}},
	)

	// Check connection for each dbObj
	var dbConns []DbConnection
	for _, check := range checker {
		check.dbConn.checkConnection(ctx, report)
		dbConns = append(dbConns, check.dbConn)
	}

	// Remaining fields will be populated by the calling service
	// Reserve: only pass fatal error to higher level
	return report, dbConns, nil
}

func (p *PostgresObj) checkConnection(ctx context.Context, report *sharedpb.RuntimeInfoReport) {
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
			if err := pgdb.Ping(); err != nil {
				report.DbStatusExtra[key] = constant.DBWriteConnectionError
				report.Status = sharedpb.RuntimeInfoReport_SICK
				log.Error(constant.DBWriteConnectionError, zap.Error(err))
				// Do not set object as it is error - report back up
			} else {
				// Set the database connection is healthy
				report.DbStatusExtra[key] = constant.DBConnectionHealthy
				report.Status = sharedpb.RuntimeInfoReport_RUNNING
				// Set conn object
				p.DBs[key] = pgdb
			}
		}

		// Log status
		log.Info(key + ":" + report.DbStatusExtra[key])
	}
}

func (r *RedisObj) checkConnection(ctx context.Context, report *sharedpb.RuntimeInfoReport) {

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
		rddb := redis.NewRedisConn(opts)
		log.Info("Redis URL",
			zap.String("host", opts.Addr),
			zap.Int("db", opts.DB),
		)

		// Check redis connection
		errConn := redis.CheckRedisConnection(rddb)
		if errConn != nil {
			report.RdStatusExtra = constant.RDConnectionError
			report.Status = sharedpb.RuntimeInfoReport_SICK
			log.Error(constant.RDConnectionError, zap.Error(errConn))
			// Do not set object as it is error - report back up
		} else {
			// Set redis connection to healthy
			report.RdStatusExtra = constant.RDConnectionHealthy
			report.Status = sharedpb.RuntimeInfoReport_RUNNING
			// Set conn object
			r.DB = rddb
		}
	}

	// Log status
	log.Info(report.RdStatusExtra)
}
