package main

import "github.com/barisaydogdu/magic-byte-organizer/internal/cli"

func main() {
	cli.Execute()

	//go run ./cmd/main.go start --dir /home/baris/Downloads/
	//go run ./cmd/main.go scan -d /home/baris/Downloads --dry-run
}
