package main

import (
	"github.com/alexivanenko/skype-backup/core"

	"log"
	"os"
	"time"

	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()

	stdLogger := log.New(os.Stdout, "", 0)
	stdLogger.Printf("[LOG] %v | %s\n", time.Now().Format("01/02/2006 - 15:04:05"), "Run Backup")

	srv, err := core.Connect(ctx)
	if err != nil {
		log.Fatalf("Unable to retrieve drive Client %v", err)
	}

	core.UploadAll(srv, ctx, stdLogger)

	stdLogger.Printf("[LOG] %v | %s\n", time.Now().Format("01/02/2006 - 15:04:05"), "Complete Backup")
}
