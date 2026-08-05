package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/spf13/viper"

	"git.ocn.com.vn/ocn/common/appcfg"
	"git.ocn.com.vn/ocn/common/logger"
	"git.ocn.com.vn/ocn/common/util/viperutil"
	"github.com/quanNaNbk03/monitor-demo/buildinfo"
	"github.com/quanNaNbk03/monitor-demo/internal/server"
)

const (
	defaultGracefulTimeout = 30 * time.Second
)

type App interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

func main() {
	// Define command-line flags
	configPath := flag.String("config", "", "Path to configuration file (default: auto-search in . and ./config)")
	versionFlag := flag.Bool("version", false, "Show version information.")
	flag.Parse()

	version := fmt.Sprintf("Version: \t %s \nBranch: \t %s\nBuild: \t %s\nOrganization: \t %s", buildinfo.Version, buildinfo.GitBranch, buildinfo.BuildDate, buildinfo.Organization)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, version)
	w.Flush()

	if *versionFlag {
		os.Exit(0)
	}

	if err := initConfig(*configPath); err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	logger.Init()

	appcfg.AppName = viper.GetString("app.name")
	if viper.GetString("app.mode") == appcfg.ModeProduction.String() {
		appcfg.SetMode(appcfg.ModeProduction)
	} else {
		appcfg.SetMode(appcfg.ModeDevelopment)
	}

	app, err := server.NewServer()
	if err != nil {
		logger.Fatal("failed to create server", "err", err)
	}

	startCtx, cancel := context.WithTimeout(context.Background(), defaultGracefulTimeout)
	defer cancel()
	if err = app.Start(startCtx); err != nil {
		logger.Fatal("failed to start app", "err", err)
	}

	interruptHandle(app)
}

func initConfig(path string) error {
	initConfigWithEnv()
	if path != "" {
		return initConfigWithFile(path)
	}
	return nil
}

func initConfigWithEnv() {
	// load env
	viper.AutomaticEnv()
	// if we use env instead of YAML, it will replace . -> _
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
}

func initConfigWithFile(path string) error {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func interruptHandle(app App) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("listening shutdown signal...")
	s := <-c
	logger.Info(fmt.Sprintf("received signal: %s. shutting down server...", s))

	stopCtx, cancel := context.WithTimeout(context.Background(), viperutil.GetWithDefault("app.gracefulTimeout", defaultGracefulTimeout))
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		logger.Error("failed to gracefully shut down server", "err", err)
	}
}
