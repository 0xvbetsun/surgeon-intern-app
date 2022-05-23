package logbookService

import (
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/graph/generated"
	"go.uber.org/zap"
)

type (
	App struct {
		qlConfig generated.Config
	}
)

func NewApp(c generated.Config) *App {
	return &App{qlConfig: c}
}

func (a *App) Run() error {
	if err := NewGraphQL("", a.qlConfig).Server.Run(); err != nil {
		zap.S().Fatalf("Failed to launch http server %s", err.Error())
	}
	return nil
}
