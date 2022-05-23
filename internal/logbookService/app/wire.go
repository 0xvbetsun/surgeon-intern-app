//go:build wireinject
// +build wireinject

package logbookService

import (
	"context"
	casbinpgadapter "github.com/cychiuae/casbin-pg-adapter"
	"github.com/google/wire"
	db2 "github.com/vbetsun/surgeon-intern-app/internal/pkg/db"
	authorization2 "github.com/vbetsun/surgeon-intern-app/internal/pkg/services/authorization"
	"github.com/vbetsun/surgeon-intern-app/internal/util"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/ogbook"
)

func InitializeApp(ctx context.Context, pgConnectionParameters db2.ConnectionParameters, casbinTableName string, casbinConfigFilePath authorization2.CasbinConfigFilePath) (*App, error) {
	wire.Build(NewApp, NewQlConfig, NewResolver, NewAuthDirectives, ServiceSet, ogbook.RepoSet, casbinpgadapter.NewAdapter, db2.NewSqlConnection, wire.FieldsOf(new(*db2.Postgresql), "conn"), util.GetViper)
	return &App{}, nil
}
