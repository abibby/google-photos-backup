package main

import (
	"context"
	"log"

	"github.com/abibby/google-photos-backup/app"
)

func main() {
	ctx := context.Background()
	err := app.Kernel.Bootstrap(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = app.Kernel.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
