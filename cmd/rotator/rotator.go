package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/ekhvalov/otus-banners-rotation/internal/environment/config"
	internalgrpc "github.com/ekhvalov/otus-banners-rotation/internal/environment/server/grpc"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
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

	server := internalgrpc.NewServer(internalgrpc.NewConfig(v))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*3)
		defer shutdownCancel()

		if grpcErr := server.Shutdown(shutdownCtx); grpcErr != nil {
			err = multierror.Append(err, grpcErr)

			//logg.Error("failed to stop grpc server: " + grpcErr.Error())
		}
	}()

	if serveErr := server.ListenAndServe(); serveErr != nil {
		err = multierror.Append(err, serveErr)
	}

	return err
}
