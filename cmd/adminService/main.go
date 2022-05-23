package main

import (
	"context"

	adminService "github.com/vbetsun/surgeon-intern-app/internal/adminService/app"
	"github.com/vbetsun/surgeon-intern-app/internal/adminService/pkg/services/encrypt"
	"github.com/vbetsun/surgeon-intern-app/internal/adminService/pkg/services/users"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/db"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/services/authorization"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/thirdparty"
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
	app, err := adminService.InitializeApp(ctx,
		db.ConnectionParameters(util.GetViper().GetString("DB_CONN_STRING")),
		util.GetViper().GetString("CASBIN_PG_TABLE"),
		authorization.CasbinConfigFilePath(util.GetViper().GetString("CASBIN_CONFIG_FILE_PATH")),
		&thirdparty.MailgunCredentials{
			ApiKey:            util.GetViper().GetString("MAILGUN_APIKEY"),
			Domain:            util.GetViper().GetString("MAILGUN_DOMAIN"),
			InviteEmailSender: util.GetViper().GetString("MAILGUN_SENDER_EMAIL"),
		},
		&users.Auth0UserDbParameters{
			DbConnectionId:   util.GetViper().GetString("AUTH0_DB_CONNECTION_ID"),
			DbConnectionName: util.GetViper().GetString("AUTH0_DB_CONNECTION_NAME"),
		},
		&encrypt.HmacParams{
			SecretKey: util.GetViper().GetString("HMAC_SECRET"),
		})
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
