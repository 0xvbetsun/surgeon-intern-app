package main

import (
	"context"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	db "github.com/vbetsun/surgeon-intern-app/internal/pkg/db"
	"github.com/vbetsun/surgeon-intern-app/internal/seeder"
	"github.com/vbetsun/surgeon-intern-app/internal/util"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.uber.org/zap"
)

func main() {
	// logging
	logger := initializeLogging()
	defer logger.Sync() // flushes buffer, if any
	ctx := context.Background()
	v := viper.GetViper()
	v.AutomaticEnv()
	if v.GetString("BOIL_DEBUG") == "true" {
		boil.DebugMode = true
	}

	s, err := seeder.InitializeSeeder(ctx, db.ConnectionParameters(v.GetString("DB_CONN_STRING")))
	if err != nil {
		panic(err)
	}
	s.SeedDevelopmentData(v.GetString("BASE_SEED_JSON_PATH"))
	zap.S().Info("Done seeding.")
}

func initializeLogging() *zap.Logger {
	logger := util.Logger()
	zap.ReplaceGlobals(logger)
	return logger
}
