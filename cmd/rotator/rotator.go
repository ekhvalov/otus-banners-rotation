package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ekhvalov/otus-banners-rotation/internal/app"
	"github.com/ekhvalov/otus-banners-rotation/internal/environment/config"
	"github.com/ekhvalov/otus-banners-rotation/internal/environment/logger"
	internalgrpc "github.com/ekhvalov/otus-banners-rotation/internal/environment/server/grpc"
	"github.com/ekhvalov/otus-banners-rotation/internal/environment/storage/redis"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	rotatorCmd = &cobra.Command{
		Use: "rotator",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

const configEnvPrefix = "rotator"

func init() {
	rotatorCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Path to config file")
}

func run() error {
	var err error
	v, err := config.NewViper(cfgFile, configEnvPrefix, config.DefaultEnvKeyReplacer)
	if err != nil {
		return fmt.Errorf("create viper error: %w", err)
	}

	storage := createStorage(v)
	queue := createEventQueue()
	logg := createLogger(v)
	rotator := app.NewRotator(storage, queue, logg)
	server := internalgrpc.NewServer(internalgrpc.NewConfig(v), rotator, logg)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*3)
		defer shutdownCancel()

		if grpcErr := server.Shutdown(shutdownCtx); grpcErr != nil {
			err = multierror.Append(err, grpcErr)
		}
	}()

	if serveErr := server.ListenAndServe(); serveErr != nil {
		err = multierror.Append(err, serveErr)
	}

	return err
}

func createStorage(v *viper.Viper) app.Storage {
	return redis.NewRedis(redis.NewConfig(v), redis.NewUUIDGenerator())
}

func createEventQueue() app.EventQueue {
	return eventQueue{}
}

func createLogger(v *viper.Viper) app.Logger {
	cfg := logger.NewConfig(v)
	return logger.NewLogger(cfg, os.Stdout)
}

type eventQueue struct{}

func (q eventQueue) Put(_ context.Context, _ app.Event) error {
	return nil
}
