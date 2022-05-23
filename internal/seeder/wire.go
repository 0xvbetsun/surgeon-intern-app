//go:build wireinject
// +build wireinject

package seeder

import (
	"context"
	"github.com/google/wire"
	db2 "github.com/vbetsun/surgeon-intern-app/internal/pkg/db"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/ogbook"
)

func InitializeSeeder(ctx context.Context, pgConnectionParameters db2.ConnectionParameters) (*Seeder, error) {
	wire.Build(New, ogbook.RepoSet, db2.NewSqlConnection, wire.FieldsOf(new(*db2.Postgresql), "conn"))
	return &Seeder{}, nil
}
