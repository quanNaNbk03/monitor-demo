package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/viper"

	"git.ocn.com.vn/ocn/common/httpbase"
	"git.ocn.com.vn/ocn/common/logger"
	"git.ocn.com.vn/ocn/common/storage"
	"github.com/quanNaNbk03/monitor-demo/internal/server/router"
	"github.com/quanNaNbk03/monitor-demo/pkg/validator"
)

type Server struct {
	httpServer    *httpbase.Server
	mysqlInstance *storage.MySQL
}

func NewServer() (*Server, error) {
	validator.Init()
	mysqlInstance, err := storage.NewMySQL(viper.GetString("mysql.dsn"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}
	serverOptions := make([]httpbase.ServerOption, 0)
	if viper.GetBool("http.server.tls.enabled") {
		tlsCfg, err := httpbase.LoadServerConfig(
			viper.GetString("http.server.tls.caCert"),
			viper.GetString("http.server.tls.serverCert"),
			viper.GetString("http.server.tls.serverKey"),
		)
		if err != nil {
			return nil, err
		}
		serverOptions = append(serverOptions, httpbase.WithTLSConfig(tlsCfg))
	}

	server := httpbase.NewServer(
		fmt.Sprintf("%s:%s", viper.GetString("http.server.host"), viper.GetString("http.server.port")),
		router.Setup(mysqlInstance),
		httpbase.ConfigTimeout{
			ReadTimeout:       viper.GetDuration("http.server.readTimeout"),
			ReadHeaderTimeout: viper.GetDuration("http.server.readHeaderTimeout"),
			WriteTimeout:      viper.GetDuration("http.server.writeTimeout"),
			IdleTimeout:       viper.GetDuration("http.server.idleTimeout"),
		},
		serverOptions...,
	)

	return &Server{
		httpServer:    server,
		mysqlInstance: mysqlInstance,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		if err := s.httpServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Failed to start HTTP server", "err", err)
		}
	}()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.httpServer.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop HTTP server: %w", err)
	}
	return nil
}
