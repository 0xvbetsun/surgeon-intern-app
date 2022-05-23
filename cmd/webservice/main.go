package main

import (
	"context"

	webservice "github.com/vbetsun/surgeon-intern-app/internal/logbookService/app"
	db2 "github.com/vbetsun/surgeon-intern-app/internal/pkg/db"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/services/authorization"
	"github.com/vbetsun/surgeon-intern-app/internal/util"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.uber.org/zap"
)

func main() {
	logger := initializeLogging()
	defer logger.Sync() // flushes buffer, if any
	ctx := context.Background()
	if util.GetViper().GetString("BOIL_DEBUG") == "true" {
		boil.DebugMode = true
	}
	app, err := webservice.InitializeApp(ctx,
		db2.ConnectionParameters(util.GetViper().GetString("DB_CONN_STRING")),
		util.GetViper().GetString("CASBIN_PG_TABLE"),
		authorization.CasbinConfigFilePath(util.GetViper().GetString("CASBIN_CONFIG_FILE_PATH")))
	if err != nil {
		panic(err)
	}
	app.Run()
}

func initializeLogging() *zap.Logger {
	logger := util.Logger()
	zap.ReplaceGlobals(logger)
	return logger
}
