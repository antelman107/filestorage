package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/antelman107/filestorage/internal/providers"
	"github.com/antelman107/filestorage/internal/services"
)

func main() {
	a := services.NewStorageApp(providers.DefaultStorageAppProviders)
	if err := a.Init(); err != nil {
		log.Fatal(err)
	}

	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelFunc()

	if err := a.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
