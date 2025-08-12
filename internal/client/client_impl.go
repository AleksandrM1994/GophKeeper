package client

import (
	"net/http"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type ClientImpl struct {
	lg         *zap.SugaredLogger
	cfg        *viper.Viper
	httpClient *http.Client
}

func NewClient(lg *zap.SugaredLogger, cfg *viper.Viper, httpClient *http.Client) Client {
	return &ClientImpl{lg: lg, cfg: cfg, httpClient: httpClient}
}
