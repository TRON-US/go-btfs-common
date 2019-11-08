package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v7"
	sharedpb "github.com/tron-us/go-btfs-common/protos/shared"
	"github.com/tron-us/go-common/db/postgres"
	dbenv "github.com/tron-us/go-common/env/db"
	"github.com/tron-us/go-common/log"
	"go.uber.org/zap"
	"os"
	"strconv"
)

const (
	DBConnectionHealthy = "DB Connection Healthy"
	DBWriteConnectionError = "Cannot connect to the write database!"
	DBReadConnectionError  = "Cannot connect to the read database!"
	DBEnvironmentError = "Cannot parse the database environment variable"
	RedisEnvironmentError = "Cannot parse the redis database environment variable"
	RedisConnectionError = "Cannot connect to the redis database!"
)

var (
	REDISClient *redis.Client
	RedisHost     = "127.0.0.1"
	RedisPort     = 6379
	RedisPassword = ""
	RedisDb       = 0
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func ConnectRedis(url string, port int, password string, db int) *redis.Client {
	var client = redis.NewClient(&redis.Options{
		Addr:     url + ":" + strconv.Itoa(port),
		Password: password, // no password set
		DB:       db,       // use default DB
	})
	return client
}

func CheckRedisConnection(connection *redis.Client) error {
	val, err := connection.Ping().Result()
	if err != nil {
		return err
	}
	if val != "PONG" {
		return fmt.Errorf("Redis respond error %s %d ", RedisHost, RedisPort)
	}
	return err
}

func CheckRuntime(ctx context.Context, runtime *sharedpb.RuntimeInfoRequest) (*sharedpb.RuntimeInfoReport, error) {
	// db runtime
	report := new(sharedpb.RuntimeInfoReport)

	// Check postgres dbWrite
	report.DbStatusExtra = []byte(DBConnectionHealthy)
	report.Signature = runtime.Signature
	report.Status = sharedpb.RuntimeInfoReport_RUNNING

	// Check database environment variable
	dbEnv, err := strconv.ParseBool(getEnv("DB_URL", "false"))
	if err != nil {
		log.Error(DBEnvironmentError, zap.Error(err))
	}
	if dbEnv == true {
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
	}

	// Check redis environment variable
	redisEnv, err := strconv.ParseBool(getEnv("REDIS_URL", "false"))
	if err != nil {
		log.Error(RedisEnvironmentError, zap.Error(err))
	}
	if redisEnv == true {
		// Check redis connection
		error := CheckRedisConnection(ConnectRedis(RedisHost, RedisPort, RedisPassword, RedisDb))
		if error != nil {
			report.DbStatusExtra = []byte(RedisConnectionError)
			log.Error(RedisConnectionError, zap.Error(err))
		}
	}

	// Remaining fields will be populated by the calling service
	// Reserve: only pass fatal error to higher level
	return report, nil
}

