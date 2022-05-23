package util

import (
	"github.com/gofrs/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Logger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("Unable to initialize Zap")
	}
	return logger
}

func GetViper() *viper.Viper {
	v := viper.New()
	v.AutomaticEnv()
	return v
}

func IsValidUUID(in string) bool {
	_, err := uuid.FromString(in)
	return err == nil
}
