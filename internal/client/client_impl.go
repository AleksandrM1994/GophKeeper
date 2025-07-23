package client

import (
	"net/http"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type ClientImpl struct {
	lg         *zap.SugaredLogger
	cfg        *viper.Viper
	httpClient *http.Client
}

func NewClient(lg *zap.SugaredLogger, cfg *viper.Viper) *ClientImpl {
	srv := &ClientImpl{lg: lg, cfg: cfg}
	srv.httpClient = &http.Client{
		Timeout: time.Second * 30,
	}
	return srv
}
