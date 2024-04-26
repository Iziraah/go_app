package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/EfimVelichkin/3rd_module_GO/03-03-umanager/internal/env"
)

const ShutdownTimeout = 3 * time.Second

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err := runMain(ctx); err != nil {
		log.Fatal(err)
	}
}

func runMain(ctx context.Context) error {
	e, c, err := env.Setup(ctx)
	if err != nil {
		return fmt.Errorf("setup.Setup: %w", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	grpcServer := e.LinksGRPCServer

	go func() {
		<-ctx.Done()
		grpcServer.Stop()
	}()

	go func() {
		defer wg.Done()
		if err := e.LinkUpdater.Run(ctx); err != nil {
			slog.Error("link updater Run: %w", err)
		}
	}()

	go func() {
		defer wg.Done()

		slog.Info(fmt.Sprintf("links grpc was started %s", e.Config.LinksService.GRPCServer.Addr))

		lis, err := net.Listen("tcp", e.Config.LinksService.GRPCServer.Addr)
		if err != nil {
			slog.Error("net Listen", slog.Any("err", err))
			return
		}

		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("net Listen", slog.Any("err", err))
			return
		}
	}()

	wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	c.Close(ctx)

	return nil
}
