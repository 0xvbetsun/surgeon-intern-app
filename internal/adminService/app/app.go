package adminService

import (
	"github.com/vbetsun/surgeon-intern-app/internal/adminService/graph/generated"
	"go.uber.org/zap"
)

type (
	App struct {
		qlConfig generated.Config
		restApi  *RestApi
	}
)

func NewApp(c generated.Config, restApi *RestApi) *App {
	return &App{qlConfig: c, restApi: restApi}
}

func (a *App) Run() error {
	if err := NewGraphQL("", a.qlConfig, a.restApi).Server.Run(); err != nil {
		zap.S().Fatalf("Failed to launch http server %s", err.Error())
	}
	return nil
}
