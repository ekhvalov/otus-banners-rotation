package internalgrpc

import (
	"context"
	"fmt"
	"net"

	"github.com/ekhvalov/otus-banners-rotation/internal/app"
	grpcapi "github.com/ekhvalov/otus-banners-rotation/pkg/api/grpc"
	"google.golang.org/grpc"
)

//nolint:lll // Ignore long line
//go:generate protoc ../../../../api/grpc/v1/rotator.proto -I ../../../../api/grpc  --go_out=../../../../pkg/api/grpc --go-grpc_out=../../../../pkg/api/grpc

func NewServer(c Config, rotator app.Rotator, logger app.Logger) Server {
	return Server{config: c, rotator: rotator, logger: logger}
}

type Server struct {
	config  Config
	rotator app.Rotator
	logger  app.Logger
	server  *grpc.Server
}

func (s *Server) ListenAndServe() error {
	address := s.config.GetAddress()
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("listen address '%s' error: %w", address, err)
	}
	s.server = grpc.NewServer()
	grpcapi.RegisterRotatorServer(s.server, &handler{rotator: s.rotator})
	s.logger.Info(fmt.Sprintf("listen on: %s", address))
	return s.server.Serve(listener)
}

func (s *Server) Shutdown(ctx context.Context) error {
	stopCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		s.server.GracefulStop()
		cancel()
	}()
	select {
	case <-ctx.Done():
		s.server.Stop()
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("stop server error: %w", err)
		}
		return nil
	case <-stopCtx.Done():
		return nil
	}
}
