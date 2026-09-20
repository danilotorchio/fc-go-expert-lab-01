package main

import (
	"fmt"
	"log/slog"
	"os"
)

func run() error {
	fmt.Println("Hello from Deploy com Cloud Run")
	return nil
}

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
