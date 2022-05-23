//go:build wireinject
// +build wireinject

package adminService

import (
	"context"
	casbinpgadapter "github.com/cychiuae/casbin-pg-adapter"
	"github.com/google/wire"
	"github.com/vbetsun/surgeon-intern-app/internal/adminService/pkg/services/encrypt"
	"github.com/vbetsun/surgeon-intern-app/internal/adminService/pkg/services/users"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/db"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/services/authorization"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/thirdparty"
	"github.com/vbetsun/surgeon-intern-app/internal/util"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/ogbook"
)

func InitializeApp(ctx context.Context, pgConnectionParameters db.ConnectionParameters, casbinTableName string, casbinConfigFilePath authorization.CasbinConfigFilePath, credentials *thirdparty.MailgunCredentials, auth0Parameters *users.Auth0UserDbParameters, hmacParams *encrypt.HmacParams) (*App, error) {
	wire.Build(NewApp, NewQlConfig, NewResolver, ServiceSet, NewRestApi, ogbook.RepoSet, casbinpgadapter.NewAdapter, db.NewSqlConnection, wire.FieldsOf(new(*db.Postgresql), "conn"), util.GetViper)
	return &App{}, nil
}
