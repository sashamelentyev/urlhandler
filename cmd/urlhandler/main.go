package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sashamelentyev/urlhandler/internal/urlhandler"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	set := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	filename := set.String("filename", "", "file name")
	reqTimeout := set.Duration("request-timeout", 10*time.Second, "")

	if err := set.Parse(os.Args[1:]); err != nil {
		return fmt.Errorf("parse args: %w", err)
	}

	h := urlhandler.New(*filename, *reqTimeout)

	ctx := context.Background()

	if err := h.Run(ctx); err != nil {
		return fmt.Errorf("run url handler: %w", err)
	}

	return nil
}
